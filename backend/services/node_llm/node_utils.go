package node_llm

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/stardustagi/TopLib/libs/jwt"
	"github.com/stardustagi/TopLib/libs/redis"
	"github.com/stardustagi/TopLib/libs/uuid"
	"github.com/stardustagi/TopModelsNode/constants"
	"github.com/stardustagi/TopModelsNode/models"
	"go.uber.org/zap"
)

// NodeUserMailEncodeToken 节点用户邮箱加密
func (n *NodeHttpService) nodeUserMailEncodeToken(mail, password string, fixedSalt string) (string, error) {
	return jwt.EncryptWithFixedSalt(password, constants.ObtenationIterations, mail, fixedSalt)
}

// NodeUserMailDecodeToken 节点用户邮箱解密
func (n *NodeHttpService) nodeUserMailDecodeToken(inPassword, dbPassword string, fixedSalt string) (string, error) {
	return jwt.DecryptWithFixedSalt(inPassword, constants.ObtenationIterations, dbPassword, fixedSalt)
}

func (n *NodeHttpService) updateNodeUser(nodeUser *models.NodeUsers) error {
	session := n.dao.NewSession()
	defer session.Close()
	_, err := session.UpdateById(nodeUser.Id, nodeUser)
	return err
}

// generateNodeUserToken 生成用户Token
func (n *NodeHttpService) generateNodeUserToken(ak, once string, nodeUserId int64) (string, string, error) {
	// 生成Token
	n.logger.Info("生成节点用户Token", zap.Int64("nodeUserId", nodeUserId), zap.String("accessKey", ak), zap.String("once", once))

	nodeUserKey := constants.NodeAccessKey(nodeUserId, ak)
	ok, err := n.rds.Exists(n.ctx, nodeUserKey)
	if err != nil && !errors.Is(err, redis.Nil) {
		n.logger.Error("检查Redis节点用户Token失败", zap.Error(err), zap.String("nodeUserKey", nodeUserKey))
		return "", "", err
	}

	// 生成token
	token := uuid.GenBytes(32)
	err = n.rds.Set(n.ctx, nodeUserKey, token, constants.NodeUserTokenExpireTimeString)
	// 查找Ak对应的Sk
	keyModel := &models.NodeKeys{
		AccessKey:  ak,
		NodeUserId: nodeUserId,
	}

	ok, err = n.dao.FindOne(keyModel)
	if err != nil || !ok {
		n.logger.Error("查找节点用户密钥失败", zap.Error(err), zap.String("accessKey", ak))
		return "", "", fmt.Errorf("节点用户密钥不存在，请联系管理员")
	}
	if keyModel.NodeId == "" {
		n.logger.Debug("新节点，需要生成节点ID", zap.String("accessKey", ak))
		keyModel.NodeId = uuid.GetUuidString()
	}
	// 设置节点用户密钥到Redis
	nodeAccessKey := constants.NodeAccessModelsKey(keyModel.NodeId)
	if err = n.rds.HSet(n.ctx, nodeAccessKey, "token", []byte(token)); err != nil {
		n.logger.Error("保存节点用户Token到Redis失败", zap.Error(err), zap.String("nodeAccessKey", nodeAccessKey))
		return "", "", err
	}
	err = n.rds.HSet(n.ctx, nodeAccessKey, "accessKey", []byte(ak))
	if err != nil {
		n.logger.Error("保存节点用户AccessKey到Redis失败", zap.Error(err), zap.String("nodeAccessKey", nodeAccessKey))
		return "", "", err
	}
	err = n.rds.HSet(n.ctx, nodeAccessKey, "securityKey", []byte(keyModel.SecurityKey))
	if err != nil {
		n.logger.Error("保存节点用户SecurityKey到Redis失败", zap.Error(err), zap.String("nodeAccessKey", nodeAccessKey))
		return "", "", err
	}
	err = n.rds.HSet(n.ctx, nodeAccessKey, "once", []byte(once))
	if err != nil {
		n.logger.Error("保存节点用户OnceToken到Redis失败", zap.Error(err), zap.String("nodeAccessKey", nodeAccessKey))
		return "", "", err
	}
	// 设置redis过期时间
	err = n.rds.Expire(n.ctx, nodeAccessKey, constants.NodeUserTokenExpireTimeString)
	if err != nil {
		n.logger.Error("设置过期失败", zap.Error(err), zap.String("nodeAccessKey", nodeAccessKey))
		return "", "", err
	}
	// 组装密钥
	jwtKey := fmt.Sprintf("%s-%s-%s-%s", ak, keyModel.SecurityKey, once, keyModel.NodeId)
	jwtString := jwt.JWTEncrypt(fmt.Sprintf("%d", nodeUserId), string(token), jwtKey)

	// 更新数据库
	session := n.dao.NewSession()
	_, err = session.UpdateById(keyModel.Id, keyModel)
	if err != nil {
		n.logger.Error("更新节点用户密钥失败", zap.Error(err), zap.String("accessKey", ak))
		// 清理redis
		_, err = n.rds.Del(n.ctx, nodeUserKey)
		_, err = n.rds.Del(n.ctx, nodeAccessKey)
		return "", "", err
	}
	return keyModel.NodeId, jwtString, nil
}

func (n *NodeHttpService) generateNodeUserJwt(nodeUserId int64) (string, string, error) {
	// 生成JWT token
	token := uuid.GetUuidString()
	jwtKey := fmt.Sprintf("%s-%s-%d", constants.AppName, constants.AppVersion, nodeUserId)
	jwtString := jwt.JWTEncrypt(fmt.Sprintf("%d", nodeUserId), token, jwtKey)

	// 存储到Redis
	key := constants.NodeUserTokenKey(nodeUserId)
	err := n.rds.Set(n.ctx, key, []byte(token), constants.NodeUserTokenExpire)
	if err != nil {
		n.logger.Error("Failed to store Token in Redis", zap.Error(err), zap.String("key", key))
		return "", "", fmt.Errorf("failed to store JWT: %w", err)
	}
	key = constants.NodeUserJwtKey(nodeUserId)
	err = n.rds.Set(n.ctx, key, []byte(jwtString), constants.NodeUserJwtExpire)
	if err != nil {
		n.logger.Error("Failed to store JWT in Redis", zap.Error(err), zap.String("key", key))
		return "", "", fmt.Errorf("failed to store JWT: %w", err)
	}
	return token, jwtString, nil
}

// startRedisKeyExpireListener 启动Redis键过期事件监听器
func (n *NodeHttpService) startRedisKeyExpireListener() {
	defer n.wg.Done()

	n.logger.Info("启动Redis键过期事件监听器")

	// 订阅Redis键过期事件，使用专门的消息通知Redis连接
	// __keyevent@0__:expired 是Redis键过期事件的通知频道
	pubSub, err := n.notifyRds.PSubscribe(n.ctx, "__keyevent@*__:expired")
	if err != nil {
		n.logger.Error("订阅Redis键过期事件失败", zap.Error(err))
		return
	}
	defer func(pubSub *redis.PubSub) {
		_ = pubSub.Close()
	}(pubSub)

	n.logger.Info("成功订阅Redis键过期事件")

	// 获取消息通道
	ch := pubSub.Channel()

	for {
		select {
		case <-n.stopCh:
			n.logger.Info("收到停止信号，退出Redis键过期事件监听器")
			return
		case msg := <-ch:
			if msg == nil {
				continue
			}

			n.logger.Debug("收到Redis键过期事件",
				zap.String("channel", msg.Channel),
				zap.String("payload", msg.Payload))

			// 检查是否是模型相关的键
			if n.isModelKey(msg.Payload) {
				n.logger.Info("检测到模型键过期", zap.String("key", msg.Payload))

				// 发送过期通知
				select {
				case n.expireNotifyCh <- msg.Payload:
					n.logger.Debug("成功发送过期键通知", zap.String("key", msg.Payload))
				default:
					n.logger.Warn("过期键通知通道已满，丢弃通知", zap.String("key", msg.Payload))
				}
			}
		case <-n.ctx.Done():
			n.logger.Info("上下文已取消，退出Redis键过期事件监听器")
			return
		}
	}
}

// handleExpiredKeys 处理过期的键
func (n *NodeHttpService) handleExpiredKeys() {
	defer n.wg.Done()

	n.logger.Info("启动过期键处理器")

	for {
		select {
		case <-n.stopCh:
			n.logger.Info("收到停止信号，退出过期键处理器")
			return
		case expiredKey := <-n.expireNotifyCh:
			n.logger.Info("处理过期键", zap.String("key", expiredKey))

			// 处理模型键过期
			err := n.handleModelKeyExpiration(expiredKey)
			if err != nil {
				n.logger.Error("处理模型键过期失败",
					zap.Error(err),
					zap.String("key", expiredKey))
			}
		case <-n.ctx.Done():
			n.logger.Info("上下文已取消，退出过期键处理器")
			return
		}
	}
}

// handleModelKeyExpiration 处理模型键过期事件
func (n *NodeHttpService) handleModelKeyExpiration(expiredKey string) error {
	n.logger.Info("处理模型键过期事件", zap.String("key", expiredKey))

	// 从键中提取NodeId
	nodeId := n.extractNodeIdFromKey(expiredKey)
	if nodeId == "" {
		n.logger.Warn("无法从键中提取NodeId", zap.String("key", expiredKey))
		return fmt.Errorf("无法从键中提取NodeId: %s", expiredKey)
	}

	n.logger.Info("提取到NodeId",
		zap.String("key", expiredKey),
		zap.String("nodeId", nodeId))

	// 从数据库查找对应的模型信息
	_, modelInfos, err := n.readModelInfo2DB(nodeId)
	if err != nil {
		n.logger.Error("从数据库读取模型信息失败",
			zap.Error(err),
			zap.String("nodeId", nodeId))
		return err
	}

	if len(modelInfos) == 0 {
		n.logger.Info("数据库中未找到对应的模型信息",
			zap.String("nodeId", nodeId))
		return nil
	}

	n.logger.Info("从数据库找到模型信息，准备更新状态为offline",
		zap.String("nodeId", nodeId),
		zap.Int("modelCount", len(modelInfos)))

	// 更新模型状态为offline
	err = n.updateModelStatusToOffline(nodeId)
	if err != nil {
		n.logger.Error("更新模型状态为offline失败",
			zap.Error(err),
			zap.String("nodeId", nodeId))
		return err
	}

	n.logger.Info("成功处理模型键过期事件",
		zap.String("key", expiredKey),
		zap.String("nodeId", nodeId))

	// 注销nodeId对应的节点
	if ok, err := n.unRegisterNodes(nodeId); err != nil || !ok {
		n.logger.Error("注销节点失败",
			zap.Error(err),
			zap.String("nodeId", nodeId))
		return err
	}
	return nil
}

// updateModelStatusToOffline 将指定nodeId的所有模型状态更新为offline
func (n *NodeHttpService) updateModelStatusToOffline(nodeId string) error {
	n.logger.Info("更新模型状态为offline", zap.String("nodeId", nodeId))

	// 先查询需要更新的模型信息
	var modelsInfoList []*models.ModelsInfo
	err := n.dao.FindMany(&modelsInfoList, "id ASC", &models.ModelsInfo{NodeId: nodeId})
	if err != nil {
		n.logger.Error("查询模型信息失败", zap.Error(err), zap.String("nodeId", nodeId))
		return err
	}

	if len(modelsInfoList) == 0 {
		n.logger.Info("未找到需要更新的模型信息", zap.String("nodeId", nodeId))
		return nil
	}

	// 使用数据库会话逐个更新模型状态
	session := n.dao.NewSession()
	defer session.Close()

	updatedCount := 0
	for _, modelInfo := range modelsInfoList {
		// 更新模型状态为offline
		modelInfo.Status = "0"
		modelInfo.LastUpdate = time.Now().Unix()

		// 使用UpdateById更新单个记录
		_, err := session.UpdateById(modelInfo.Id, modelInfo)
		if err != nil {
			n.logger.Error("更新模型状态失败",
				zap.Error(err),
				zap.String("nodeId", nodeId),
				zap.String("modelId", modelInfo.ModelId),
				zap.Int64("id", modelInfo.Id))
			continue
		}
		updatedCount++
	}

	n.logger.Info("成功更新模型状态为offline",
		zap.String("nodeId", nodeId),
		zap.Int("totalModels", len(modelsInfoList)),
		zap.Int("updatedCount", updatedCount))

	return nil
}

// isModelKey 检查键是否是模型相关的键
func (n *NodeHttpService) isModelKey(key string) bool {
	// 检查键是否以模型前缀开头
	return strings.HasPrefix(key, constants.ModelsKeyPrefix)
}

// extractNodeIdFromKey 从Redis键中提取NodeId
// 键的格式应该是: ModelsKeyPrefix:NodeId
func (n *NodeHttpService) extractNodeIdFromKey(key string) string {
	// 例如: "models:node123"
	prefix := constants.ModelsKeyPrefix + ":"
	if after, ok := strings.CutPrefix(key, prefix); ok {
		return after
	}
	return ""
}

func (n *NodeHttpService) checkModelIdExists(modelId int64) bool {
	modelInfo := &models.ModelsInfo{}
	ok, err := n.dao.FindById(modelId, modelInfo)
	if err != nil || !ok {
		n.logger.Info("checkModelExists call failed", zap.Error(err))
		return false
	}
	return true
}

func (n *NodeHttpService) checkNodeIdExists(nodeId int64) bool {
	nodeInfo := &models.Nodes{}
	ok, err := n.dao.FindById(nodeId, nodeInfo)
	if err != nil || !ok {
		n.logger.Info("checkNodeExists call failed", zap.Error(err))
		return false
	}
	return true
}

func (n *NodeHttpService) ownerNodeCheck(nodeUserId int64, nodeId string) (bool, error) {
	n.logger.Info("LLMUserOwnerNodeCheck called", zap.Int64("nodeUserId", nodeUserId), zap.String("nodeId", nodeId))
	nodeKeysModel := &models.NodeKeys{
		NodeUserId: nodeUserId,
		NodeId:     nodeId,
	}
	session := n.dao.NewSession()
	defer session.Close()
	return session.FindOne(nodeKeysModel)
}

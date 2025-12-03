package node_llm

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/stardustagi/TopLib/libs/redis"
	"github.com/stardustagi/TopModelsNode/constants"
	"github.com/stardustagi/TopModelsNode/models"
	"go.uber.org/zap"
)

func (n *NodeHttpService) updateNodeKeepLive(nodeId int64, keepInfo []NodeKeepLiveInfo) error {
	n.logger.Info("updateModelKeepLive called", zap.Int64("nodeId", nodeId))

	modelKeepLiveKey := constants.ModelsNodeKeepLiveKey(nodeId)
	_ = n.rds.Expire(n.ctx, modelKeepLiveKey, constants.ModelsNodeExpireTimeString)
	for _, info := range keepInfo {
		keepLiveDataBytes, err := json.Marshal(info)
		if err != nil {
			n.logger.Error("模型KeepLive信息序列化失败", zap.Error(err), zap.Int64("modelName", info.ModelId))
			return err
		}
		err = n.rds.HSet(n.ctx, modelKeepLiveKey, fmt.Sprintf("%d", info.ModelId), keepLiveDataBytes)
		if err != nil {
			n.logger.Error("写入Redis模型KeepLive信息失败", zap.Error(err), zap.Int64("modelName", info.ModelId))
			return err
		}
	}
	_ = n.rds.Expire(n.ctx, modelKeepLiveKey, constants.ModelsNodeExpireTimeString)

	n.logger.Info("模型KeepLive信息更新成功", zap.Int64("nodeId", nodeId))
	return nil
}

// writeModelInfo2Redis 保存模型信息到Redis
func (n *NodeHttpService) writeModelInfo2Redis(nodeId int64, modelsInfo []ModelInfo) (bool, error) {
	n.logger.Info("SaveModelInfo2Redis called", zap.Int64("nodeId", nodeId), zap.Int("modelCount", len(modelsInfo)))

	nodeInfoKey := constants.NodeInfoKey(nodeId)
	// 写入redis
	for _, model := range modelsInfo {
		modelKey := fmt.Sprintf("%d", model.ID)
		modelData, err := json.Marshal(model)
		if err != nil {
			n.logger.Error("模型序列化失败", zap.Error(err), zap.Int64("modelId", model.ID))
			return false, err
		}
		err = n.rds.HSet(n.ctx, nodeInfoKey, modelKey, modelData)
		if err != nil {
			n.logger.Error("写入Redis失败", zap.Error(err), zap.Int64("modelId", model.ID))
			return false, err
		}
		_ = n.rds.Expire(n.ctx, nodeInfoKey, constants.ModelsNodeExpireTimeString)
		// 增加KeepLive key
		modelKeepLiveKey := constants.ModelsNodeKeepLiveKey(nodeId)
		keepLiveData := NodeKeepLiveInfo{
			ModelId:    model.ID,
			Metrics:    model.Metrics,
			ExpireTime: model.ExpireTime,
			KeepLive:   time.Now().Unix(),
		}
		keepLiveDataBytes, err := json.Marshal(keepLiveData)
		if err != nil {
			n.logger.Error("模型KeepLive信息序列化失败", zap.Error(err), zap.Int64("modelId", model.ID))
			return false, err
		}
		err = n.rds.HSet(n.ctx, modelKeepLiveKey, modelKey, keepLiveDataBytes)
		if err != nil {
			n.logger.Error("写入Redis模型KeepLive信息失败", zap.Error(err), zap.Int64("modelId", model.ID))
			return false, err
		}
		_ = n.rds.Expire(n.ctx, modelKeepLiveKey, constants.ModelsNodeExpireTimeString)
	}

	n.logger.Info("模型信息保存到Redis成功", zap.Int64("nodeId", nodeId), zap.Int("modelCount", len(modelsInfo)))
	return true, nil
}

// readModelInfo2DB 从数据库读取模型信息
func (n *NodeHttpService) readModelInfo2DB(nodeId int64) ([]ModelInfo, error) {
	n.logger.Info("readModelInfo2DB called", zap.Int64("nodeName", nodeId))
	session := n.dao.NewSession()
	defer session.Close()
	// 1. 查询所有映射
	var mapInfos []models.NodeModelsInfoMaps
	err := session.FindMany(&mapInfos, "id desc", &models.NodeModelsInfoMaps{NodeId: nodeId})
	if err != nil {
		return nil, err
	}

	// 2. 按ModelId分组聚合ProviderId
	modelProvidersMap := make(map[int64][]int64)
	for _, m := range mapInfos {
		modelProvidersMap[m.ModelId] = append(modelProvidersMap[m.ModelId], m.ModelProviderId)
	}

	var modelInfos []ModelInfo
	for modelId, providerIds := range modelProvidersMap {
		// 查模型
		dbModel := &models.ModelsInfo{Id: modelId}
		ok, err := session.FindOne(dbModel)
		if !ok || err != nil {
			continue
		}

		// 查所有provider
		var providers []Provider
		for _, pid := range providerIds {
			dbProvider := &models.ModelsProvider{Id: pid}
			ok, err := session.FindOne(dbProvider)
			if !ok || err != nil {
				continue
			}
			providers = append(providers, Provider{
				ID:          dbProvider.ProviderId,
				Type:        dbProvider.Type,
				Name:        dbProvider.Name,
				Endpoint:    dbProvider.Endpoint,
				APIType:     dbProvider.ApiType,
				APIKeys:     strings.Split(dbProvider.ApiKeys, ","),
				ModelName:   dbProvider.ModelName,
				InputPrice:  dbProvider.InputPrice,
				OutputPrice: dbProvider.OutputPrice,
				CachePrice:  dbProvider.CachePrice,
			})
		}

		// 组装ModelInfo
		modelInfos = append(modelInfos, ModelInfo{
			ID:          dbModel.Id,
			Name:        dbModel.Name,
			APIVersion:  dbModel.ApiVersion,
			DeployName:  dbModel.DeployName,
			InputPrice:  dbModel.InputPrice,
			OutputPrice: dbModel.OutputPrice,
			CachePrice:  dbModel.CachePrice,
			Address:     dbModel.Address,
			Status:      dbModel.Status,
			ApiStyles:   strings.Split(dbModel.ApiStyles, ","),
			Providers:   providers,
		})
	}
	return modelInfos, nil

}

// readModelInfo2Redis 从Redis读取模型信息
func (n *NodeHttpService) readModelInfo2Redis(nodeId int64) ([]ModelInfo, error) {
	n.logger.Info("readModelInfo2Redis called", zap.Int64("nodeName", nodeId))

	// 检查Key是否存在
	infoKey := constants.NodeInfoKey(nodeId)
	exists, err := n.rds.Exists(n.ctx, infoKey)
	if err != nil && !errors.Is(err, redis.Nil) {
		n.logger.Error("检查Redis Key失败", zap.Error(err), zap.Int64("nodeName", nodeId))
		return nil, err
	}

	if !exists {
		n.logger.Info("Redis中不存在该节点的模型信息", zap.Int64("nodeId", nodeId))
		return []ModelInfo{}, nil
	}

	// 获取所有字段
	modelDataMap, err := n.rds.HGetAll(n.ctx, infoKey)
	if err != nil {
		n.logger.Error("从Redis读取模型信息失败", zap.Error(err), zap.Int64("nodeId", nodeId))
		return nil, err
	}
	var modelInfos []ModelInfo
	for modelId, modelData := range modelDataMap {
		var modelInfo ModelInfo
		err := json.Unmarshal(modelData, &modelInfo)
		if err != nil {
			n.logger.Error("模型信息反序列化失败",
				zap.Error(err),
				zap.String("modelId", modelId),
				zap.Int64("nodeName", nodeId))
			continue
		}
		modelInfos = append(modelInfos, modelInfo)
	}

	n.logger.Info("从Redis读取模型信息成功",
		zap.Int64("nodeName", nodeId),
		zap.Int("modelCount", len(modelInfos)))

	return modelInfos, nil
}

// IsRunning 检查服务是否正在运行
func (n *NodeHttpService) IsRunning() bool {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.isRunning
}

// GetExpireNotifyChannel 获取过期通知通道（用于测试或外部监控）
func (n *NodeHttpService) GetExpireNotifyChannel() <-chan string {
	return n.expireNotifyCh
}

// SendExpireNotification 手动发送过期通知（用于测试）
func (n *NodeHttpService) SendExpireNotification(key string) error {
	if !n.IsRunning() {
		return fmt.Errorf("服务未运行")
	}

	select {
	case n.expireNotifyCh <- key:
		n.logger.Debug("手动发送过期通知成功", zap.String("key", key))
		return nil
	default:
		return fmt.Errorf("过期通知通道已满")
	}
}

// unRegisterNodes 模型节点注销 - 批量注销多个节点的模型
func (n *NodeHttpService) unRegisterNodes(nodeId int64) (bool, error) {
	n.logger.Info("llmUnregisterNodes called", zap.Int64("node id", nodeId))
	nodeAccessKey := constants.NodeAccessModelsKey(nodeId)
	exists, err := n.rds.Exists(n.ctx, nodeAccessKey)
	if err != nil && !errors.Is(err, redis.Nil) {
		n.logger.Error("检查Redis Key失败", zap.Error(err), zap.Int64("nodeId", nodeId))
		return false, err
	}

	if !exists {
		n.logger.Info("节点不存在，无需注销", zap.Int64("nodeId", nodeId))
		return false, nil // 节点不存在，直接返回
	}

	// 删除Redis中的数据
	_, err = n.rds.Del(n.ctx, nodeAccessKey)
	if err != nil {
		n.logger.Error("删除Redis中的模型信息失败", zap.Error(err), zap.Int64("nodeId", nodeId))
		return false, err
	}

	// 更新数据库中的模型状态为offline
	err = n.updateModelStatusToOffline(nodeId)
	if err != nil {
		n.logger.Error("更新模型状态为offline失败", zap.Error(err), zap.Int64("nodeId", nodeId))
		return false, err
	}

	return true, nil
}

// NodeUpsetAccessKeyAndSecurityKey 添加节点用户的AccessKey和SecurityKey
func (n *NodeHttpService) NodeUpsetAccessKeyAndSecurityKey(nodeInfo NodeInfo) (*models.Nodes, error) {
	bean := &models.Nodes{
		OwnerId: nodeInfo.NodeUserId,
		Name:    nodeInfo.Name,
	}
	session := n.dao.NewSession()
	ok, err := session.FindOne(bean)
	if err != nil {
		n.logger.Error("检查节点是否存在失败",
			zap.Error(err),
			zap.Int64("nodeUserId", nodeInfo.NodeUserId))
		return nil, err
	}
	if !ok {
		n.logger.Warn("节点不存在，无法添加AccessKey和SecurityKey",
			zap.Int64("nodeUserId", nodeInfo.NodeUserId))
		return nil, fmt.Errorf("节点不存在")
	}
	bean.AccessKey = nodeInfo.AccessKey
	bean.SecurityKey = nodeInfo.SecretKey
	bean.LastupdateAt = time.Now().Unix()
	_, err = session.UpdateById(bean.Id, bean)
	if err != nil {
		n.logger.Error("更新节点AccessKey和SecurityKey失败",
			zap.Error(err),
			zap.Int64("nodeUserId", nodeInfo.NodeUserId))
		return nil, err
	}
	n.logger.Info("节点AccessKey和SecurityKey更新成功",
		zap.Int64("nodeUserId", nodeInfo.NodeUserId))
	return bean, nil
}

// NodeCheckin 在线Node检查
func (n *NodeHttpService) NodeCheckin(nodeId int64) (bool, error) {
	n.logger.Info("NodeCheckin called", zap.Int64("nodeId", nodeId))
	if nodeId <= 0 {
		n.logger.Error("无效的nodeId", zap.Int64("nodeId", nodeId))
		return false, fmt.Errorf("无效的nodeId")
	}
	// 检查nodeId是否注册过
	key := constants.NodeAccessModelsKey(nodeId)
	if ok, err := n.rds.Exists(n.ctx, key); (err != nil && !errors.Is(err, redis.Nil)) || !ok {
		n.logger.Error("节点未注册", zap.Error(err), zap.Int64("nodeId", nodeId))
		return false, err
	}
	return true, nil
}

// getNodeIdModelsInfo 获取节点的模型信息
func (n *NodeHttpService) getNodeIdModelsInfo(nodeId int64) ([]ModelInfo, error) {
	n.logger.Info("GetNodeIdModelsInfo called", zap.Int64("nodeName", nodeId))

	// 首先尝试从Redis读取
	//modelInfos, err := n.readModelInfo2Redis(nodeId)
	//if err != nil {
	//	n.logger.Warn("从Redis读取模型信息失败，尝试从数据库读取",
	//		zap.Error(err),
	//		zap.Int64("nodeId", nodeId))
	//} else if len(modelInfos) > 0 {
	//	n.logger.Info("从Redis成功读取模型信息",
	//		zap.Int64("nodeId", nodeId),
	//		zap.Int("modelCount", len(modelInfos)))
	//	return modelInfos, nil
	//}

	// Redis中没有数据，从数据库读取
	modelInfos, err := n.readModelInfo2DB(nodeId)
	if err != nil {
		n.logger.Error("从数据库读取模型信息失败",
			zap.Error(err),
			zap.Int64("nodeId", nodeId))
		return nil, err
	}

	// 将数据库中的数据缓存到Redis
	if len(modelInfos) > 0 {
		success, err := n.writeModelInfo2Redis(nodeId, modelInfos)
		if !success || err != nil {
			n.logger.Warn("缓存模型信息到Redis失败",
				zap.Error(err),
				zap.Int64("nodeId", nodeId))
		}
	}

	n.logger.Info("获取模型信息成功",
		zap.Int64("nodeId", nodeId),
		zap.Int("modelCount", len(modelInfos)))

	//for _, modelInfo := range modelInfos {
	//	modelInfos = append(modelInfos, modelInfo)
	//}

	return modelInfos, nil
}

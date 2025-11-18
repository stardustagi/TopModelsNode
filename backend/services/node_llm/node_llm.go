package node_llm

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/stardustagi/TopLib/libs/databases"
	"github.com/stardustagi/TopLib/libs/redis"
	"github.com/stardustagi/TopLib/libs/uuid"
	"github.com/stardustagi/TopModelsLogin/backend/services/system"
	message "github.com/stardustagi/TopModelsNode/backend/services/nats"
	"github.com/stardustagi/TopModelsNode/constants"
	"github.com/stardustagi/TopModelsNode/models"
	"go.uber.org/zap"
)

// NodeLogin 节点登录
func (n *NodeHttpService) NodeLogin(ak, once, pwd string, nodeUser *models.NodeUsers) (string, string, error) {
	n.logger.Info("LLMNodeUserCheck called", zap.String("nodeUserLogin", nodeUser.Email))

	session := n.dao.NewSession()
	ok, err := session.FindOne(nodeUser)
	if err != nil {
		n.logger.Error("节点用户登录失败", zap.Error(err), zap.String("email", nodeUser.Email))
		return "", "", err
	}
	if !ok {
		n.logger.Error("用户不存在或登录失败", zap.String("email", nodeUser.Email))
		return "", "", fmt.Errorf("用户不存在或登录失败: %s", nodeUser.Email)
	}
	// 验证密码
	vEmail, err := system.NodeUserMailDecodeToken(pwd, nodeUser.Password, nodeUser.Salt)
	if err != nil || vEmail != nodeUser.Email {
		n.logger.Error("节点用户登录失败，密码错误", zap.String("email", nodeUser.Email), zap.Error(err))
		return "", "", errors.New("节点用户登录失败，密码错误")
	}
	n.logger.Info("节点用户登录成功", zap.String("email", nodeUser.Email))

	// 生成Token
	nodeId, jwtToken, err := n.generateNodeUserToken(ak, once, nodeUser.Id)
	if err != nil {
		n.logger.Error("节点用户Token生成失败", zap.String("email", nodeUser.Email), zap.Error(err))
		return "", "", err
	}
	return nodeId, jwtToken, nil
}

// WriteModelInfo2Db 保存模型信息到数据库
func (n *NodeHttpService) WriteModelInfo2Db(nodeId, address string, modelsInfo []ModelInfo) (bool, error) {
	n.logger.Info("SaveModelInfo2Db called")
	// 数据库会话
	transactionSession := databases.NewSessionWrapper(n.dao)
	upsetModelsInfo := make([]*models.ModelsInfo, 0, len(modelsInfo))
	upsetProvider := make([]*models.ModelsProvider, 0)
	for _, modelInfo := range modelsInfo {
		model := &models.ModelsInfo{
			ModelId:     modelInfo.ID,
			NodeId:      nodeId,
			Name:        modelInfo.Name,
			ApiVersion:  modelInfo.APIVersion,
			DeployName:  modelInfo.DeployName,
			InputPrice:  modelInfo.InputPrice,
			OutputPrice: modelInfo.OutputPrice,
			CachePrice:  modelInfo.CachePrice,
			Status:      modelInfo.Status,
			Address:     address,
			LastUpdate:  time.Now().Unix(),
		}
		upsetModelsInfo = append(upsetModelsInfo, model)

		for _, provider := range modelInfo.Providers {
			providerModel := &models.ModelsProvider{
				ModelId:     model.ModelId,
				Type:        provider.Type,
				Name:        provider.Name,
				Endpoint:    provider.Endpoint,
				ApiType:     provider.APIType,
				ApiKeys:     strings.Join(provider.APIKeys, ","),
				ModelName:   provider.ModelName,
				NodeId:      nodeId,
				ProviderId:  provider.ID,
				InputPrice:  provider.InputPrice,
				OutputPrice: provider.OutputPrice,
				CachePrice:  provider.CachePrice,
				LastUpdate:  time.Now().Unix(),
			}
			upsetProvider = append(upsetProvider, providerModel)
		}
	}
	trErr := transactionSession.Execute(func(session databases.SessionDao) error {
		// 逐个插入模型信息
		for _, model := range upsetModelsInfo {
			_, err := session.Upsert(
				&models.ModelsInfo{ModelId: model.ModelId, NodeId: model.NodeId},
				model)
			if err != nil {
				n.logger.Error("保存模型信息到数据库失败", zap.Error(err))
				return err
			}
		}
		return nil
	}).Execute(func(session databases.SessionDao) error {
		// 逐个插入提供商信息
		for _, provider := range upsetProvider {
			_, err := session.Upsert(&models.ModelsProvider{ModelId: provider.ModelId, NodeId: provider.NodeId}, provider)
			if err != nil {
				n.logger.Error("保存模型提供商信息到数据库失败", zap.Error(err))
				return err
			}
		}
		return nil
	}).CommitAndClose()

	if trErr != nil {
		n.logger.Error("数据库事务执行失败", zap.Error(trErr))
		return false, trErr
	}

	// 假设保存成功
	n.logger.Info("模型信息保存到数据库成功")
	return true, nil
}

// WriteModelInfo2Redis 保存模型信息到Redis
func (n *NodeHttpService) WriteModelInfo2Redis(nodeId, address string, modelsInfo []ModelInfo) (bool, error) {
	n.logger.Info("SaveModelInfo2Redis called", zap.String("nodeId", nodeId))

	// 写入node address
	if err := n.rds.HSet(n.ctx, nodeId, "address", []byte(address)); err != nil {
		n.logger.Error("写入节点地址到Redis失败", zap.Error(err), zap.String("nodeId", nodeId), zap.String("address", address))
		return false, err
	}

	// 写入redis
	for _, model := range modelsInfo {
		modelKey := fmt.Sprintf("%s", model.ID)
		modelData, err := json.Marshal(model)
		if err != nil {
			n.logger.Error("模型序列化失败", zap.Error(err), zap.String("modelId", model.ID))
			return false, err
		}
		err = n.rds.HSet(n.ctx, nodeId, modelKey, modelData)
		if err != nil {
			n.logger.Error("写入Redis失败", zap.Error(err), zap.String("modelId", model.ID))
			return false, err
		}
	}
	_ = n.rds.Expire(n.ctx, nodeId, constants.ModelsNodeExpireTimeString)

	n.logger.Info("模型信息保存到Redis成功", zap.String("nodeId", nodeId))
	return true, nil
}

// ReadModelInfo2DB 从数据库读取模型信息
func (n *NodeHttpService) ReadModelInfo2DB(nodeId string) (string, []ModelInfo, error) {
	n.logger.Info("ReadModelInfo2DB called", zap.String("nodeId", nodeId))

	// 查询模型信息
	var modelsInfoList []*models.ModelsInfo
	var nodeAddress string
	err := n.dao.FindMany(&modelsInfoList, "id ASC", &models.ModelsInfo{NodeId: nodeId})
	if err != nil {
		n.logger.Error("查询模型信息失败", zap.Error(err), zap.String("nodeId", nodeId))
		return "", nil, err
	}

	// 查询提供商信息
	var providersMap = make(map[string][]*models.ModelsProvider)
	for _, modelInfo := range modelsInfoList {
		var providers []*models.ModelsProvider
		err := n.dao.FindMany(&providers, "id ASC", &models.ModelsProvider{
			ModelId: modelInfo.ModelId,
			NodeId:  nodeId,
		})
		if err != nil {
			n.logger.Error("查询提供商信息失败", zap.Error(err),
				zap.String("modelId", modelInfo.ModelId),
				zap.String("nodeId", nodeId))
			continue
		}
		providersMap[modelInfo.ModelId] = providers
	}

	// 转换为ModelInfo结构
	var modelInfos []ModelInfo
	for _, dbModel := range modelsInfoList {
		modelInfo := ModelInfo{
			ID:          dbModel.ModelId,
			Name:        dbModel.Name,
			APIVersion:  dbModel.ApiVersion,
			DeployName:  dbModel.DeployName,
			InputPrice:  dbModel.InputPrice,
			OutputPrice: dbModel.OutputPrice,
			CachePrice:  dbModel.CachePrice,
			Status:      dbModel.Status,
			Providers:   make([]Provider, 0),
		}
		if nodeAddress == "" {
			nodeAddress = dbModel.Address
		}
		// 添加提供商信息
		if providers, exists := providersMap[dbModel.ModelId]; exists {
			for _, dbProvider := range providers {
				provider := Provider{
					ID:          dbProvider.ProviderId,
					Type:        dbProvider.Type,
					Name:        dbProvider.Name,
					Endpoint:    dbProvider.Endpoint,
					APIType:     dbProvider.ApiType,
					ModelName:   dbProvider.ModelName,
					InputPrice:  dbProvider.InputPrice,
					OutputPrice: dbProvider.OutputPrice,
					CachePrice:  dbProvider.CachePrice,
				}

				// 解析API Keys
				if dbProvider.ApiKeys != "" {
					provider.APIKeys = strings.Split(dbProvider.ApiKeys, ",")
				}

				modelInfo.Providers = append(modelInfo.Providers, provider)
			}
		}

		modelInfos = append(modelInfos, modelInfo)
	}

	n.logger.Info("从数据库读取模型信息成功",
		zap.String("nodeId", nodeId),
		zap.Int("modelCount", len(modelInfos)))

	return nodeAddress, modelInfos, nil
}

// ReadModelInfo2Redis 从Redis读取模型信息
func (n *NodeHttpService) ReadModelInfo2Redis(nodeId string) (string, []ModelInfo, error) {
	n.logger.Info("ReadModelInfo2Redis called", zap.String("nodeId", nodeId))

	// key := fmt.Sprintf("%s:%s", constants.ModelsKeyPrefix, nodeId)

	// 检查Key是否存在
	exists, err := n.rds.Exists(n.ctx, nodeId)
	if err != nil && !errors.Is(err, redis.Nil) {
		n.logger.Error("检查Redis Key失败", zap.Error(err), zap.String("nodeId", nodeId))
		return "", nil, err
	}

	if !exists {
		n.logger.Info("Redis中不存在该节点的模型信息", zap.String("nodeId", nodeId))
		return "", []ModelInfo{}, nil
	}

	// 获取所有字段
	modelDataMap, err := n.rds.HGetAll(n.ctx, nodeId)
	if err != nil {
		n.logger.Error("从Redis读取模型信息失败", zap.Error(err), zap.String("nodeId", nodeId))
		return "", nil, err
	}
	address := ""
	var modelInfos []ModelInfo
	for modelId, modelData := range modelDataMap {
		if modelId == "address" {
			address = string(modelData)
			continue
		}
		var modelInfo ModelInfo
		err := json.Unmarshal(modelData, &modelInfo)
		if err != nil {
			n.logger.Error("模型信息反序列化失败",
				zap.Error(err),
				zap.String("modelId", modelId),
				zap.String("nodeId", nodeId))
			continue
		}
		modelInfos = append(modelInfos, modelInfo)
	}

	n.logger.Info("从Redis读取模型信息成功",
		zap.String("nodeId", nodeId),
		zap.Int("modelCount", len(modelInfos)))

	return address, modelInfos, nil
}

// GetModelInfo 获取模型信息（优先从Redis读取，Redis不存在则从数据库读取并缓存到Redis）
func (n *NodeHttpService) GetModelInfo(nodeId string) (string, []ModelInfo, error) {
	n.logger.Info("GetModelInfo called", zap.String("nodeId", nodeId))
	var nodeAddress string
	// 首先尝试从Redis读取
	address, modelInfos, err := n.ReadModelInfo2Redis(nodeId)
	if err != nil {
		n.logger.Warn("从Redis读取模型信息失败，尝试从数据库读取",
			zap.Error(err),
			zap.String("nodeId", nodeId))
	} else if len(modelInfos) > 0 {
		n.logger.Info("从Redis成功读取模型信息",
			zap.String("nodeId", nodeId),
			zap.Int("modelCount", len(modelInfos)))
		return address, modelInfos, nil
	}

	// Redis中没有数据，从数据库读取
	nodeAddress, modelInfos, err = n.ReadModelInfo2DB(nodeId)
	if err != nil {
		n.logger.Error("从数据库读取模型信息失败",
			zap.Error(err),
			zap.String("nodeId", nodeId))
		return "", nil, err
	}

	// 将数据库中的数据缓存到Redis
	if len(modelInfos) > 0 {
		// key := fmt.Sprintf("%s:%s", constants.ModelsKeyPrefix, nodeId)
		success, err := n.WriteModelInfo2Redis(nodeAddress, nodeId, modelInfos)
		if !success || err != nil {
			n.logger.Warn("缓存模型信息到Redis失败",
				zap.Error(err),
				zap.String("nodeId", nodeId))
		}
	}

	n.logger.Info("获取模型信息成功",
		zap.String("nodeId", nodeId),
		zap.Int("modelCount", len(modelInfos)))

	return nodeAddress, modelInfos, nil
}

// DeleteModelInfo 删除模型信息（同时删除Redis和数据库中的数据）
func (n *NodeHttpService) DeleteModelInfo(nodeId string) error {
	n.logger.Info("DeleteModelInfo called", zap.String("nodeId", nodeId))

	// 删除Redis中的数据
	_, err := n.rds.Del(n.ctx, nodeId)
	if err != nil {
		n.logger.Error("删除Redis中的模型信息失败", zap.Error(err), zap.String("nodeId", nodeId))
	}

	// 删除数据库中的数据
	transactionSession := databases.NewSessionWrapper(n.dao)
	err = transactionSession.Execute(func(session databases.SessionDao) error {
		// 删除提供商信息
		_, err := session.Delete(&models.ModelsProvider{NodeId: nodeId})
		if err != nil {
			n.logger.Error("删除数据库中的提供商信息失败", zap.Error(err), zap.String("nodeId", nodeId))
			return err
		}
		return nil
	}).Execute(func(session databases.SessionDao) error {
		// 删除模型信息
		_, err := session.Delete(&models.ModelsInfo{NodeId: nodeId})
		if err != nil {
			n.logger.Error("删除数据库中的模型信息失败", zap.Error(err), zap.String("nodeId", nodeId))
			return err
		}
		return nil
	}).Commit()

	if err != nil {
		n.logger.Error("删除模型信息事务执行失败", zap.Error(err), zap.String("nodeId", nodeId))
		return err
	}

	n.logger.Info("删除模型信息成功", zap.String("nodeId", nodeId))
	return nil
}

// GetAllNodesModelInfo 获取所有节点的模型信息
func (n *NodeHttpService) GetAllNodesModelInfo(page databases.Pageable) ([][]map[string]interface{}, error) {
	n.logger.Info("GetAllNodesModelInfo called")

	session := n.dao.NewSession()
	defer session.Close()
	// 查询所有模型信息
	//var result []*responses.ModelWithProviderResp
	//err := session.Query(&result, "call GetModelsWithProviders(?,?,?)", page.Limit(), page.Skip(), page.Sort())
	result, err := session.CallProcedure("GetModelsWithProviders", page.Limit(), page.Skip(), page.Sort())
	if err != nil {
		n.logger.Error("查询所有模型信息失败", zap.Error(err))
		return nil, err
	}

	n.logger.Info("获取所有节点模型信息成功", zap.Int("nodeCount", len(result)))
	return result, nil
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

// LLMUnregisterNodes 模型节点注销 - 批量注销多个节点的模型
func (n *NodeHttpService) UnRegisterNodes(nodeId string) (bool, error) {
	n.logger.Info("llmUnregisterNodes called", zap.Int("node id", len(nodeId)))

	exists, err := n.rds.Exists(n.ctx, constants.NodeAccessModelsKey(nodeId))
	if err != nil && !errors.Is(err, redis.Nil) {
		n.logger.Error("检查Redis Key失败", zap.Error(err), zap.String("nodeId", nodeId))
		return false, err
	}

	if !exists {
		n.logger.Info("节点不存在，无需注销", zap.String("nodeId", nodeId))
		return false, nil // 节点不存在，直接返回
	}

	// 删除Redis中的数据
	_, err = n.rds.Del(n.ctx, nodeId)
	if err != nil {
		n.logger.Error("删除Redis中的模型信息失败", zap.Error(err), zap.String("nodeId", nodeId))
		return false, err
	}

	// 更新数据库中的模型状态为offline
	err = n.updateModelStatusToOffline(nodeId)
	if err != nil {
		n.logger.Error("更新模型状态为offline失败", zap.Error(err), zap.String("nodeId", nodeId))
		return false, err
	}

	return true, nil
}

// LLMKeepAliveNodes 模型节点心跳 - 批量更新多个节点的心跳状态
func (n *NodeHttpService) KeepAliveNodes(nodeId, address string, nodeModels []ModelInfo) (map[string]bool, error) {
	if len(nodeModels) == 0 {
		n.logger.Warn("节点心跳调用时，模型列表为空", zap.String("nodeId", nodeId))
		return map[string]bool{nodeId: false}, fmt.Errorf("模型列表不能为空")
	}
	n.logger.Info("LLMKeepAliveNodes called", zap.Int("modelCount", len(nodeModels)))

	results := make(map[string]bool)

	// 设置模型状态为请求中的状态（保持一致）
	for i := range nodeModels {
		if nodeModels[i].Status == "" {
			nodeModels[i].Status = "online"
		}
	}

	// 更新Redis TTL为1小时
	nodeAccessKey := constants.NodeAccessModelsKey(nodeId)
	err := n.rds.Expire(n.ctx, nodeAccessKey, constants.ModelsNodeExpireTimeString)
	if err != nil {
		n.logger.Error("更新Redis TTL失败", zap.Error(err), zap.String("nodeId", nodeId))
	}

	// 重新设置Redis数据
	success, err := n.WriteModelInfo2Redis(nodeId, address, nodeModels)
	if err != nil || !success {
		n.logger.Error("更新Redis模型信息失败", zap.Error(err), zap.String("nodeId", nodeId))
		results[nodeId] = false
		return results, err
	}

	// 更新数据库中的数据状态
	success, err = n.WriteModelInfo2Db(nodeId, address, nodeModels)
	if err != nil || !success {
		n.logger.Error("更新数据库模型信息失败", zap.Error(err), zap.String("nodeId", nodeId))
		results[nodeId] = false
		return results, err
	}

	n.logger.Info("节点心跳更新成功",
		zap.String("nodeId", nodeId),
		zap.Int("modelCount", len(nodeModels)))
	results[nodeId] = true

	return results, nil
}

// LLMUserAddAccessKeyAndSecurityKey 添加节点用户的AccessKey和SecurityKey
func (n *NodeHttpService) UserAddAccessKeyAndSecurityKey(nodeUserId int64) (bool, error) {
	accessKey := uuid.GenString(16)
	secretKey := uuid.GenString(16)
	key := &models.NodeKeys{
		NodeUserId:  nodeUserId,
		AccessKey:   accessKey,
		SecurityKey: secretKey,
	}
	session := n.dao.NewSession()
	ok, err := session.Exists(key)
	if err != nil {
		n.logger.Debug("数据库查询错误", zap.Error(err), zap.Int64("node user", nodeUserId))
		return false, err
	}
	if ok {
		n.logger.Info("节点用户已存在AccessKey和SecurityKey，请勿重复", zap.Int64("node user", nodeUserId),
			zap.String("accessKey", accessKey),
			zap.String("securityKey", secretKey))
		return false, fmt.Errorf("节点用户已存在AccessKey和SecurityKey，请勿重复")
	}
	key.LastupdateAt = time.Now().Unix()
	_, err = session.InsertOne(key)
	if err != nil {
		n.logger.Debug("数据库插入错误", zap.Error(err), zap.Int64("node user", nodeUserId))
		return false, err
	}
	return true, nil
}

// LLMNodeBillingUsage 节点计费上报
func (n *NodeHttpService) NodeBillingUsage(nodeId string, usage []LLMUsageReport) error {
	n.logger.Info("LLMNodeBillingUsage called", zap.String("nodeId", nodeId), zap.Int("usageCount", len(usage)))
	if len(usage) == 0 {
		n.logger.Warn("节点计费上报调用时，使用量列表为空", zap.String("nodeId", nodeId))
		return fmt.Errorf("使用量列表不能为空")
	}
	if nodeId == "" {
		n.logger.Warn("节点计费上报调用时，节点ID为空")
		return fmt.Errorf("节点ID不能为空")
	}
	// 检查nodeId是否注册过
	key := constants.NodeAccessModelsKey(nodeId)
	if ok, err := n.rds.Exists(n.ctx, key); (err != nil && !errors.Is(err, redis.Nil)) || !ok {
		n.logger.Error("节点未注册", zap.Error(err), zap.String("nodeId", nodeId))
		return err
	}
	// 保存到数据库
	session := n.dao.NewSession()
	defer session.Close()
	for _, u := range usage {
		// 无效模型和私有模型不记录
		if u.ModelID == "" || u.IsPrivate == 1 {
			n.logger.Warn("使用量报告中包含无效条目，跳过", zap.String("nodeId", nodeId), zap.Any("usage", u))
			continue
		}
		stream := 0
		if u.Stream {
			stream = 1
		}
		record := &models.LlmUsageReport{
			Id:                        u.ID,
			NodeId:                    nodeId,
			ModelId:                   u.ModelID,
			ActualModel:               u.ActualModel,
			Provider:                  u.Provider,
			ActualProvider:            u.ActualProvider,
			Caller:                    u.Caller,
			CallerKey:                 u.CallerKey,
			ClientVersion:             u.ClientVersion,
			TokenUsageInputTokens:     u.TokenUsage.InputTokens,
			TokenUsageOutputTokens:    u.TokenUsage.OutputTokens,
			TokenUsageCachedTokens:    u.TokenUsage.CachedTokens,
			TokenUsageReasoningTokens: u.TokenUsage.ReasoningTokens,
			TokenUsageTokensPerSec:    u.TokenUsage.TokensPerSec,
			TokenUsageLatency:         u.TokenUsage.Latency,
			AgentVersion:              u.AgentVersion,
			Stream:                    stream,
		}
		tbName := record.GetSliceName(u.ID)
		if ok, err := session.Native().IsTableExist(tbName); err != nil || !ok {
			err = session.Native().Table(tbName).CreateTable(&record)
			if err != nil {
				n.logger.Error("创建分表失败", zap.Error(err), zap.String("table", tbName))
				return err
			}
		}
		_, err := session.Native().Table(tbName).Insert(&record)
		if err != nil {
			n.logger.Error("插入使用量报告失败", zap.Error(err), zap.String("table", tbName), zap.Any("usage", u))
			return err
		}
	}
	now := time.Now().Unix()
	for _, us := range usage {
		us.CreatedAt = now
	}
	// 发布到nats
	byteString, err := json.Marshal(usage)
	if err != nil {
		n.logger.Error("使用量报告序列化失败", zap.Error(err), zap.String("nodeId", nodeId))
		return err
	}
	msgSrv := message.GetNatsQueueInstance()
	if ok := msgSrv.PublisherStreamAsync("billing.nodeUsage", byteString); !ok {
		n.logger.Error("使用量报告发布到NATS失败", zap.String("nodeId", nodeId))
		return fmt.Errorf("使用量报告发布到NATS失败")
	}
	return nil
}

// NodeCheckin 在线Node检查
func (n *NodeHttpService) NodeCheckin(nodeId string) (bool, error) {
	n.logger.Info("NodeCheckin called", zap.String("nodeId", nodeId))
	if nodeId == "" {
		n.logger.Warn("节点ID为空")
		return false, fmt.Errorf("节点ID不能为空")
	}
	// 检查nodeId是否注册过
	key := constants.NodeAccessModelsKey(nodeId)
	if ok, err := n.rds.Exists(n.ctx, key); (err != nil && !errors.Is(err, redis.Nil)) || !ok {
		n.logger.Error("节点未注册", zap.Error(err), zap.String("nodeId", nodeId))
		return false, err
	}
	return true, nil
}

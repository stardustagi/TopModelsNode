package node_llm

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stardustagi/TopLib/libs/databases"
	"github.com/stardustagi/TopLib/libs/logs"
	"github.com/stardustagi/TopLib/libs/redis"
	"github.com/stardustagi/TopLib/libs/server"
	"github.com/stardustagi/TopLib/libs/uuid"
	"github.com/stardustagi/TopLib/protocol"
	"github.com/stardustagi/TopModelsNode/backend"
	"github.com/stardustagi/TopModelsNode/constants"
	"github.com/stardustagi/TopModelsNode/models"
	"github.com/stardustagi/TopModelsNode/protocol/requests"
	"github.com/stardustagi/TopModelsNode/protocol/responses"

	"go.uber.org/zap"
)

type NodeHttpService struct {
	logger      *zap.Logger
	ctx         context.Context
	cancelCtx   context.CancelFunc
	dao         databases.BaseDao
	rds         redis.RedisCli
	mu          sync.RWMutex // 读写锁保护usersJwt
	userInfosMu sync.RWMutex // 读写保护UserInfos
	app         *backend.Application

	// redis相关
	notifyRds      redis.RedisCli // 专门用于消息通知的Redis连接
	isRunning      bool
	expireNotifyCh chan string
	stopCh         chan struct{}
	wg             sync.WaitGroup
}

var (
	nodeHttpServiceInstance *NodeHttpService
	nodeHttpServiceOnce     sync.Once
)

// GetNodeHttpServiceInstance  获取 HTTP 用户服务实例
func GetNodeHttpServiceInstance() *NodeHttpService {
	nodeHttpServiceOnce.Do(func() {
		nodeHttpServiceInstance = NewNodeHttpService()
	})
	return nodeHttpServiceInstance
}

// NewNodeHttpService  创建新的 HTTP 用户服务
func NewNodeHttpService() *NodeHttpService {
	ctx, cancel := context.WithCancel(context.Background())
	return &NodeHttpService{
		logger:    logs.GetLogger("UsersHttpService"),
		ctx:       ctx,
		cancelCtx: cancel,
		dao:       databases.GetDao(),
		rds: redis.NewRedisView(redis.GetRedisDb(),
			constants.NodeUserKeyPrefix,
			logs.GetLogger("NodeUserRedis")),
		notifyRds: redis.NewRedisView(redis.GetRedisDb(),
			"",
			logs.GetLogger("notifyRedisView")), // 专门用于消息通知，不设置前缀

	}
}

func (n *NodeHttpService) Start(app *backend.Application) {
	if app == nil {
		panic("请设置后端应用")
	}
	n.app = app
	n.initialization()
	n.logger.Info("Starting UsersHttpService...")

	// redis key过期相关服务
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.isRunning {
		n.logger.Info("UsersHttpService already running")
		return
	}

	n.logger.Info("启动key 过期服务")
	n.wg.Add(1)
	go n.startRedisKeyExpireListener()
	n.isRunning = true
	n.logger.Info("Redis key过期启动完成")
}

func (n *NodeHttpService) Stop() {
	n.logger.Info("Stopping UsersHttpService...")
	n.cancelCtx()
	n.mu.Lock()
	defer n.mu.Unlock()
	if !n.isRunning {
		n.logger.Info("UsersHttpService already stopped")
		return
	}
	close(n.stopCh)
	n.wg.Wait()

	n.stopCh = make(chan struct{})
	n.isRunning = false
	n.logger.Info("UsersHttpService stopped.")
}

func (n *NodeHttpService) initialization() {
	n.app.AddGroup("node/llm", server.Request(), backend.NodeAccess())
}

// AddModelsInfos 添加模型信息
// godoc
// 添加模型信息
// @Summary 添加模型信息
// @Description 添加模型信息
// @Tags node/llm
// @Accept json
// @Produce json
// @Param request body requests.AddModelsInfoRequest true "请求参数"
// @Success 200 {object} responses.DefaultResponse
// @Failure 400 {object} responses.DefaultResponse
// @Failure 500 {object} responses.DefaultResponse
// @Router /node/llm/addModelsInfos [post]
func (n *NodeHttpService) AddModelsInfos(ctx echo.Context,
	req requests.AddModelsInfoRequest,
	resp responses.DefaultResponse) error {
	session := n.dao.NewSession()
	defer session.Close()
	// 1. 检查模型是否存在
	modelsInfos := make([]*models.ModelsInfo, len(req.ModelsInfos))
	for _, modelInfo := range req.ModelsInfos {
		m := &models.ModelsInfo{
			Name:       modelInfo.Name,
			ApiVersion: modelInfo.ApiVersion,
		}
		ok, err := session.FindOne(m)
		if err != nil {
			n.logger.Info("databases error：", zap.Error(err))
			continue
		}
		if ok {
			continue
		}
		// 不存在的模型添加
		m.DeployName = modelInfo.DeployName
		m.InputPrice = modelInfo.InputPrice
		m.OutputPrice = modelInfo.OutputPrice
		m.CachePrice = modelInfo.CachePrice
		m.Status = modelInfo.Status
		m.IsPrivate = modelInfo.IsPrivate
		m.OwnerId = modelInfo.OwnerId
		m.Address = modelInfo.Address
		m.LastUpdate = time.Now().Unix()
		num, err := session.InsertOne(m)
		if err != nil {
			n.logger.Info("databases error：", zap.Error(err))
			continue
		}
		modelsInfos = append(modelsInfos, m)
		n.logger.Info("insert success", zap.Int64("id", num))
	}
	return protocol.Response(ctx, nil, modelsInfos)
}

// AddModelsProvider 添加模型供应商
// godoc
// @Summary 添加模型供应商
// @Description 添加模型供应商
// @Tags NodeLLM
// @Accept json
// @Produce json
// @Param request body requests.AddModelsProviderInfoRequest true "请求参数"
// @Success 200 {object} responses.DefaultResponse
// @Router /node/llm/addModelsProvider [post]
func (n *NodeHttpService) AddModelsProvider(ctx echo.Context,
	req requests.AddModelsProviderInfoRequest,
	resp responses.DefaultResponse) error {
	n.logger.Info("AddModelsProvider is Called", zap.Any("req", req))
	session := n.dao.NewSession()
	defer session.Close()
	// 1. 检查模型服务商是否存在
	modelsProviders := make([]*models.ModelsProvider, len(req.ModelsProviderInfo))
	for _, v := range req.ModelsProviderInfo {
		modelsProvider := &models.ModelsProvider{
			Name:     v.Name,
			Type:     v.Type,
			Endpoint: v.Endpoint,
		}
		ok, err := session.FindOne(modelsProvider)
		if err != nil {
			n.logger.Info("databases error：", zap.Error(err))
			continue
		}
		if ok {
			continue
		}
		// 不存在的模型添加
		modelsProvider.ApiType = v.ApiType
		modelsProvider.ModelName = v.ModelName
		modelsProvider.InputPrice = v.InputPrice
		modelsProvider.OutputPrice = v.OutputPrice
		modelsProvider.CachePrice = v.CachePrice
		modelsProvider.ApiKeys = v.ApiKeys
		modelsProvider.LastUpdate = time.Now().Unix()
		if _, err := session.InsertOne(modelsProvider); err != nil {
			n.logger.Info("databases error:", zap.Error(err))
			continue
		}
		modelsProviders = append(modelsProviders, modelsProvider)
	}
	return protocol.Response(ctx, nil, modelsProviders)
}

// MapModelsProviderInfoToNode 映射模型和模型供应商
// godoc
// @Summary 映射模型和模型供应商
// @Description 映射模型和模型供应商
// @Tags NodeLLM
// @Accept json
// @Produce json
// @Param request body requests.MapModelsProviderInfoToNodeRequest true "请求参数"
// @Success 200 {object} responses.DefaultResponse
// @Router /node/llm/mapModelsProviderInfoToNode [post]
func (n *NodeHttpService) MapModelsProviderInfoToNode(ctx echo.Context,
	req requests.MapModelsProviderInfoToNodeRequest,
	resp responses.DefaultResponse) error {
	n.logger.Info("call MapModelsProviderInfoToNode", zap.Any("req", req))
	session := n.dao.NewSession()
	defer session.Close()
	// 检查Node和mode的合法师
	if !(n.checkModelIdExists(req.ModelId) && n.checkNodeIdExists(req.NodeId)) {
		return protocol.Response(ctx, nil, "node or model not exists")
	}
	// 设置条件
	where := &models.NodeModelsInfoMaps{
		NodeId:  req.NodeId,
		ModelId: req.ModelId,
	}
	count := 0
	for _, providerIds := range req.ProviderIds {
		nodeModelsInfoMap := &models.NodeModelsInfoMaps{
			NodeId:          req.NodeId,
			ModelId:         req.ModelId,
			ModelProviderId: providerIds,
			CreatedAt:       time.Now().Unix(),
			UpdatedAt:       time.Now().Unix(),
		}
		if _, err := session.Upsert(where, nodeModelsInfoMap); err != nil {
			n.logger.Info("databases error:", zap.Error(err))
			continue
		} else {
			count++
		}
	}
	return protocol.Response(ctx, nil, fmt.Sprintf("update success %d", count))
}

// UpsetNodeInfos 添加/更新模型信息
// godoc
// @Summary 添加/更新模型信息
// @Description 添加/更新模型信息
// @Tags node
// @Accept json
// @Produce json
// @Param request body requests.UpsetNodeInfoRequest true "请求参数"
// @Success 200 {object} responses.DefaultResponse
// @Router /node/llm/upsetNodeInfos [post]
func (n *NodeHttpService) UpsetNodeInfos(ctx echo.Context,
	req requests.UpsetNodeInfoRequest,
	resp responses.DefaultResponse) error {
	n.logger.Info("call UpsetNodeInfos", zap.Any("req", req))
	session := n.dao.NewSession()
	defer session.Close()
	if req.Code == "" {
		req.Code = uuid.GenString(12)
	}
	where := &models.Nodes{
		Ids:        req.Code,
		NodeUserId: req.NodeUserId,
	}
	node := &models.Nodes{
		Ids:          req.Code,
		NodeUserId:   req.NodeUserId,
		LastupdateAt: time.Now().Unix(),
		Domain:       req.Domain,
	}
	if _, err := session.Upsert(where, node); err != nil {
		n.logger.Info("databases error:", zap.Error(err))
		return protocol.Response(ctx, constants.ErrInternalServer.AppendErrors(err), nil)
	}
	return protocol.Response(ctx, nil, "success")
}

// ListNodeInfos 获取模型信息
// godoc
// @Summary 获取模型信息
// @Description 获取模型信息
// @Tags node
// @Accept json
// @Produce json
// @Param request body requests.ListNodeInfoRequest true "请求参数"
// @Success 200 {object} responses.DefaultResponse
// @Router /node/llm/listNodeInfos [post]
func (n *NodeHttpService) ListNodeInfos(ctx echo.Context,
	req requests.ListNodeInfoRequest,
	resp responses.DefaultResponse) error {
	n.logger.Info("List Node Infos call", zap.Any("req", req))
	session := n.dao.NewSession()
	defer session.Close()
	// 默认排序
	if req.PageInfo.Sort == "" {
		req.PageInfo.Sort = "node_id asc"
	}
	result, err := session.CallProcedure("ListNodeUserNodeInfos",
		req.NodeUserId, req.PageInfo.Skip, req.PageInfo.Limit, req.PageInfo.Sort)
	if err != nil {
		n.logger.Error("ListNodeUserNodeInfos error:", zap.Error(err))
		return protocol.Response(ctx, constants.ErrInternalServer.AppendErrors(err), nil)
	}

	return protocol.Response(ctx, nil, result)
}

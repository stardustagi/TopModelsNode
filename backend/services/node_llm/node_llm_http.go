package node_llm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stardustagi/TopLib/libs/databases"
	"github.com/stardustagi/TopLib/libs/logs"
	"github.com/stardustagi/TopLib/libs/redis"
	"github.com/stardustagi/TopLib/protocol"
	"github.com/stardustagi/TopModelsNode/backend"
	message "github.com/stardustagi/TopModelsNode/backend/services/nats"
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
		stopCh:    make(chan struct{}),
		rds: redis.NewRedisView(redis.GetRedisDb(),
			constants.NodeKeyPrefix,
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
	nodeUserId, err := n.getNodeIdFromContext(ctx)
	if err != nil {
		n.logger.Error("ListNodeInfos get nodeUserId failed", zap.Error(err))
		return protocol.Response(ctx, constants.ErrAuthFailed.AppendErrors(err), nil)
	}
	// 默认排序
	if req.PageInfo.Sort == "" {
		req.PageInfo.Sort = "node_id asc"
	}
	result, err := session.CallProcedure("ListNodeUserNodeInfos",
		nodeUserId, req.PageInfo.Skip, req.PageInfo.Limit, req.PageInfo.Sort)
	if err != nil {
		n.logger.Error("ListNodeUserNodeInfos error:", zap.Error(err))
		return protocol.Response(ctx, constants.ErrInternalServer.AppendErrors(err), nil)
	}

	return protocol.Response(ctx, nil, result)
}

// NodeLogin 节点登录
// @Summary 节点登录
// @Description 节点登录
// @Tags node
// @Accept json
// @Produce json
// @Param request body requests.NodeLoginReq true "请求参数"
// @Success 200 {object} responses.NodeLoginResp
// @Router /node/public/nodeLogin [post]
func (n *NodeHttpService) NodeLogin(ctx echo.Context, req requests.NodeLoginReq, resp responses.NodeLoginResp) error {
	n.logger.Info("NodeLogin called", zap.String("nodeUserLogin", req.Mail))

	nodeUser := &models.NodeUsers{
		Email: req.Mail,
	}
	session := n.dao.NewSession()
	ok, err := session.FindOne(nodeUser)
	if err != nil {
		n.logger.Error("节点用户登录失败", zap.Error(err), zap.String("email", nodeUser.Email))
		return protocol.Response(ctx, constants.ErrInternalServer.AppendErrors(err), nil)
	}
	if !ok {
		n.logger.Error("用户不存在或登录失败", zap.String("email", nodeUser.Email))
		return protocol.Response(ctx, constants.ErrNotDataSet, nil)
	}
	// 验证密码
	vEmail, err := n.nodeUserMailDecodeToken(req.Password, nodeUser.Password, nodeUser.Salt)
	if err != nil || vEmail != nodeUser.Email {
		n.logger.Error("节点用户登录失败，密码错误", zap.String("email", nodeUser.Email), zap.Error(err))
		return protocol.Response(ctx, constants.ErrAuthFailed, nil)
	}
	n.logger.Info("节点用户登录成功", zap.String("email", nodeUser.Email))

	// 生成Token
	nodeInfo, jwtToken, err := n.generateNodeLoginToken(req.AccessToken, req.Once, nodeUser.Id)
	if err != nil {
		n.logger.Error("节点用户Token生成失败", zap.String("email", nodeUser.Email), zap.Error(err))
		return protocol.Response(ctx, constants.ErrAuthFailed.AppendErrors(err), nil)
	}
	// 查找 name
	resp.NodeId = nodeInfo.Id
	resp.NodeName = nodeInfo.Name
	resp.Jwt = jwtToken
	resp.Once = req.Once
	resp.AccessKey = req.AccessToken
	resp.Address = nodeInfo.Domain
	modelsConfigs, err := n.getNodeIdModelsInfo(nodeInfo.Id)
	if err != nil {
		n.logger.Error("获取节点模型配置失败", zap.String("nodeName", nodeInfo.Name), zap.Error(err))
		return protocol.Response(ctx, constants.ErrInternalServer.AppendErrors(err), nil)
	}
	resp.Config = modelsConfigs
	ctx.Response().Header().Set("nodeId", fmt.Sprintf("%d", nodeInfo.Id))
	ctx.Response().Header().Set("jwt", jwtToken)
	ctx.Response().Header().Set("accessKey", req.AccessToken)
	ctx.Response().Header().Set("once", req.Once)
	return protocol.Response(ctx, nil, resp)
}

func (n *NodeHttpService) KeepLive(ctx echo.Context, req requests.NodeKeepLiveReq, resp responses.DefaultResponse) error {
	n.logger.Info("Node KeepLive called", zap.Any("nodeId", req))
	nodeId, err := n.getNodeIdFromContext(ctx)
	if err != nil {
		n.logger.Error("Node KeepLive get nodeId failed", zap.Error(err))
		return protocol.Response(ctx, constants.ErrAuthFailed.AppendErrors(err), nil)
	}
	if nodeId != req.NodeId {
		n.logger.Error("Node KeepLive nodeId mismatch", zap.Int64("nodeIdFromContext", nodeId), zap.Int64("nodeIdFromReq", req.NodeId))
		return protocol.Response(ctx, constants.ErrAuthFailed.AppendErrors(fmt.Errorf("节点ID不匹配")), nil)
	}
	// 更新心跳
	keepLiveDatas := make([]NodeKeepLiveInfo, len(req.Info))
	for i, info := range req.Info {
		keepLiveDatas[i] = NodeKeepLiveInfo{
			ModelId: info.ModelId,
			Metrics: ModelMetrics{
				Latency:     info.Metrics.Latency,
				HealthScore: info.Metrics.HealthScore,
			},
			ExpireTime: info.ExpireTime,
			KeepLive:   time.Now().Unix(),
		}
	}
	err = n.updateNodeKeepLive(nodeId, keepLiveDatas)
	if err != nil {
		n.logger.Error("Node KeepLive update failed", zap.Error(err))
		return protocol.Response(ctx, constants.ErrInternalServer.AppendErrors(err), nil)
	}
	return protocol.Response(ctx, nil, "keep live success")
}

// NodeBillingUsage 节点计费上报
// @Summary 节点计费上报
// @Description 节点计费上报
// @Tags node
// @Accept json
// @Produce json
// @Param request body requests.NodeReportUsageReq true "请求参数"
// @Success 200 {object} responses.DefaultResponse
// @Router /node/llm/billingUsage [post]
func (n *NodeHttpService) NodeBillingUsage(ctx echo.Context, req requests.NodeReportUsageReq, resp responses.DefaultResponse) error {
	n.logger.Info("LLMNodeBillingUsage called",
		zap.Int64("nodeId", req.NodeId),
		zap.Int("usageCount", len(req.Report)))
	if len(req.Report) == 0 {
		n.logger.Warn("节点计费上报调用时，使用量列表为空", zap.Int64("nodeId", req.NodeId))
		return protocol.Response(ctx,
			constants.ErrInternalServer.AppendErrors(fmt.Errorf("使用量列表不能为空")),
			nil)
	}
	if req.NodeId <= 0 {
		n.logger.Warn("节点计费上报调用时，节点ID为空")
		return protocol.Response(ctx,
			constants.ErrInternalServer.AppendErrors(fmt.Errorf("节点ID不能为空")),
			nil)
	}
	// 检查nodeId是否注册过
	key := constants.NodeAccessModelsKey(req.NodeId)
	if ok, err := n.rds.Exists(n.ctx, key); (err != nil && !errors.Is(err, redis.Nil)) || !ok {
		n.logger.Error("节点未注册", zap.Error(err), zap.Int64("nodeId", req.NodeId))
		return protocol.Response(ctx,
			constants.ErrInternalServer.AppendErrors(fmt.Errorf("节点未注册")),
			nil)
	}
	// 保存到数据库
	session := n.dao.NewSession()
	defer session.Close()
	for _, u := range req.Report {
		// 无效模型和私有模型不记录
		if u.ModelID <= 0 || u.IsPrivate == 1 {
			n.logger.Warn("使用量报告中包含无效条目，跳过", zap.Int64("nodeId", req.NodeId), zap.Any("usage", u))
			continue
		}
		stream := 0
		if u.Stream {
			stream = 1
		}
		usage, err := json.Marshal(u.TokenUsage)
		if err != nil {
			n.logger.Error("使用量数据错亶在，跳过", zap.Int64("nodeId", req.NodeId), zap.Any("usage", u), zap.Error(err))
			continue
		}
		record := &models.LlmUsageReport{
			Id:             u.ID,
			NodeId:         req.NodeId,
			ModelId:        u.ModelID,
			ActualModel:    u.ActualModel,
			Provider:       u.Provider,
			ActualProvider: u.ActualProvider,
			Caller:         u.Caller,
			CallerKey:      u.CallerKey,
			ClientVersion:  u.ClientVersion,
			AgentVersion:   u.AgentVersion,
			Stream:         stream,
			Usage:          string(usage),
		}
		tbName := record.GetSliceName(u.ID, 10)
		if ok, err := session.Native().IsTableExist(tbName); err != nil || !ok {
			err = session.Native().Table(tbName).CreateTable(&record)
			if err != nil {
				n.logger.Error("创建分表失败", zap.Error(err), zap.String("table", tbName))
				return protocol.Response(ctx,
					constants.ErrInternalServer.AppendErrors(fmt.Errorf("创建分表失败: %v", err)),
					nil)
			}
		}
		_, err = session.Native().Table(tbName).Insert(&record)
		if err != nil {
			n.logger.Error("插入使用量报告失败", zap.Error(err), zap.String("table", tbName), zap.Any("usage", u))
			return protocol.Response(ctx,
				constants.ErrInternalServer.AppendErrors(err),
				nil)
		}
	}
	now := time.Now().Unix()
	for _, us := range req.Report {
		us.CreatedAt = now
	}
	// 发布到nats
	byteString, err := json.Marshal(req.Report)
	if err != nil {
		n.logger.Error("使用量报告序列化失败", zap.Error(err), zap.Int64("nodeId", req.NodeId))
		return protocol.Response(ctx,
			constants.ErrInternalServer.AppendErrors(err),
			nil)
	}
	msgSrv := message.GetNatsQueueInstance()
	if ok := msgSrv.PublisherStreamAsync("billing.nodeUsage", byteString); !ok {
		n.logger.Error("使用量报告发布到NATS失败", zap.Int64("nodeId", req.NodeId))
		return protocol.Response(ctx,
			constants.ErrInternalServer.AppendErrors(fmt.Errorf("使用量报告发布失败")),
			nil)
	}
	return protocol.Response(ctx, nil, "success")
}

// NodeUnregister LLM节点用户注销
// godoc
// @Summary LLM节点用户注销
// @Description LLM节点用户注销
// @Tags LLM模型节点管理
// @Accept json
// @Produce json
// @Param nodeUserId query int true "节点用户ID"
// @Param request body requests.NodeUnRegisterReq true "节点用户注销请求"
// @Success 200 {object} responses.DefaultResponse "成功"
// @Failure 400 {object} responses.DefaultResponse "参数错误"
// @Failure 500 {object} responses.DefaultResponse "服务器错误"
// @Router /node/llm/nodeUnRegister [post]
func (n *NodeHttpService) NodeUnregister(c echo.Context, req requests.NodeUnRegisterReq, resp responses.DefaultResponse) error {
	logger := logs.GetLogger("NodeUnregister call")
	_nodeUserId := c.QueryParams().Get("id")
	nodeUserId, err := strconv.ParseInt(_nodeUserId, 10, 64)
	if err != nil {
		logger.Error("Failed to parse node user id", zap.String("nodeUserId", _nodeUserId), zap.Error(err))
		return protocol.Response(c, constants.ErrInvalidParams, nil)
	}
	logger.Info("LLM node user unregister requested", zap.Any("request", req))

	ok, nodeId, err := n.ownerNodeCheck(nodeUserId, req.Name)
	if err != nil {
		logger.Error("Failed to get node user owner node", zap.Error(err))
		return protocol.Response(c, constants.ErrInternalServer.AppendErrors(err), nil)
	}
	if !ok {
		logger.Error("Node user does not own this node",
			zap.Int64("nodeUserId", nodeUserId), zap.String("nodeName", req.Name))
		return protocol.Response(c, constants.ErrNodeUserNotOwnNode, nil)
	}
	// Unregister node user
	if ok, err := n.unRegisterNodes(nodeId); err != nil || !ok {
		logger.Error("Failed to unregister node user", zap.Error(err), zap.String("", req.Mail))
		return protocol.Response(c, constants.ErrInternalServer, nil)
	}

	logger.Info("Successfully unregistered node user", zap.String("mail", req.Mail))
	return protocol.Response(c, nil, "注销成功")
}

// CheckUserBalance 检查用户钱包
// @Summary 检查用户钱包
// @Description 检查用户钱包
// @Tag 节点服务
// @Accept json
// @Produce json
// @Param request body requests.UserBalanceReq true "请求参数"
// @Success 200 {object} responses.UserBalanceResp
// @Failure 400 {object} responses.DefaultResponse "参数错误"
// @Failure 500 {object} responses.DefaultResponse "服务器错误"
// @Router /node/llm/checkUserBalance [post]
func (n *NodeHttpService) CheckUserBalance(ctx echo.Context,
	req requests.UserBalanceReq, resp responses.UserBalanceResp) error {
	n.logger.Info("CheckUserBalance", zap.Any("request", req))
	nodeId, err := n.getNodeIdFromContext(ctx)
	if err != nil {
		n.logger.Error("invalid node id", zap.Error(err))
		return protocol.Response(ctx, constants.ErrInvalidParams.AppendErrors(err), nil)
	}
	ok, err := n.NodeCheckin(nodeId)
	if err != nil {
		n.logger.Error("Node check fail", zap.Error(err))
		return protocol.Response(ctx, constants.ErrInternalServer.AppendErrors(err), nil)
	}
	if !ok {
		err = fmt.Errorf("node user does not own node %d", nodeId)
		n.logger.Error("Node check fail", zap.Error(err))
		return protocol.Response(ctx, constants.ErrNodeUserNotOwnNode, nil)
	}
	wallet, err := n.getUserWalletBalance(req.UserID, req.WalletType)
	if err != nil {
		return protocol.Response(ctx, constants.ErrGetUserBalance.AppendErrors(err), nil)
	}
	resp.Balance = wallet.Balance
	resp.WalletType = wallet.WalletType
	resp.WalletAddress = wallet.WalletAddress

	return protocol.Response(ctx, nil, resp)
}

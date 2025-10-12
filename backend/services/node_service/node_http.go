package node_service

import (
	"context"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/stardustagi/TopLib/libs/databases"
	"github.com/stardustagi/TopLib/libs/logs"
	"github.com/stardustagi/TopLib/libs/redis"
	"github.com/stardustagi/TopLib/libs/server"
	"github.com/stardustagi/TopLib/protocol"
	"github.com/stardustagi/TopModelsLogin/constants"
	"github.com/stardustagi/TopModelsLogin/protocol/responses"
	"github.com/stardustagi/TopModelsNode/backend"
	"github.com/stardustagi/TopModelsNode/protocol/requests"

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
}

var (
	usersHttpServiceInstance *NodeHttpService
	usersHttpServiceOnce     sync.Once
)

// GetUsersHttpServiceInstance 获取 HTTP 用户服务实例
func GetUsersHttpServiceInstance() *NodeHttpService {
	usersHttpServiceOnce.Do(func() {
		usersHttpServiceInstance = NewUsersHttpService()
	})
	return usersHttpServiceInstance
}

// NewUsersHttpService 创建新的 HTTP 用户服务
func NewUsersHttpService() *NodeHttpService {
	ctx, cancel := context.WithCancel(context.Background())
	return &NodeHttpService{
		logger:    logs.GetLogger("UsersHttpService"),
		ctx:       ctx,
		cancelCtx: cancel,
		dao:       databases.GetDao(),
		rds: redis.NewRedisView(redis.GetRedisDb(),
			constants.NodeUserKeyPrefix,
			logs.GetLogger("NodeUserRedis")),
	}
}

func (n *NodeHttpService) Start(app *backend.Application) {
	if app == nil {
		panic("请设置后端应用")
	}
	n.app = app
	n.initialization()
	n.logger.Info("Starting UsersHttpService...")
}

func (n *NodeHttpService) Stop() {
	n.logger.Info("Stopping UsersHttpService...")
	n.cancelCtx()
	n.logger.Info("UsersHttpService stopped.")
}

func (n *NodeHttpService) initialization() {
	n.app.AddGroup("node", server.Request())
	n.app.AddPostHandler("node", server.NewHandler(
		"register",
		[]string{"Node"},
		n.RegisterNode,
	))
}

func (n *NodeHttpService) RegisterNode(c echo.Context, req requests.DefaultRequest, resp responses.DefaultResponse) error {
	return protocol.Response(c, nil, "待实现")
}

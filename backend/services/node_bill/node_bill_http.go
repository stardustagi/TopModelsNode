package node_bill

import (
	"context"
	"sync"

	"github.com/stardustagi/TopLib/libs/databases"
	"github.com/stardustagi/TopLib/libs/logs"
	"github.com/stardustagi/TopLib/libs/redis"
	"github.com/stardustagi/TopModelsNode/backend"
	"github.com/stardustagi/TopModelsNode/constants"
	"go.uber.org/zap"
)

type NodeHttpBillService struct {
	logger    *zap.Logger
	ctx       context.Context
	cancelCtx context.CancelFunc
	dao       databases.BaseDao
	rds       redis.RedisCli
	mu        sync.RWMutex
	app       *backend.Application
}

var (
	nodeHttpBillServiceInstance *NodeHttpBillService
	nodeHttpBillServiceOnce     sync.Once
)

// GetNodeHttpBillServiceInstance 获取账单服务实例
func GetNodeHttpBillServiceInstance() *NodeHttpBillService {
	nodeHttpBillServiceOnce.Do(func() {
		nodeHttpBillServiceInstance = NewNodeHttpBillService()
	})
	return nodeHttpBillServiceInstance
}

// NewNodeHttpBillService 创建新的账单服务
func NewNodeHttpBillService() *NodeHttpBillService {
	ctx, cancel := context.WithCancel(context.Background())
	return &NodeHttpBillService{
		logger:    logs.GetLogger("NodeHttpBillService"),
		ctx:       ctx,
		cancelCtx: cancel,
		dao:       databases.GetDao(),
		rds: redis.NewRedisView(redis.GetRedisDb(),
			constants.NodeKeyPrefix,
			logs.GetLogger("NodeBillRedis")),
	}
}

func (n *NodeHttpBillService) Start(app *backend.Application) {
	if app == nil {
		panic("请设置后端应用")
	}
	n.app = app
	n.initialization()
	n.logger.Info("Starting NodeHttpBillService...")
}

func (n *NodeHttpBillService) Stop() {
	n.logger.Info("Stopping NodeHttpBillService...")
	n.cancelCtx()
	n.logger.Info("NodeHttpBillService stopped.")
}

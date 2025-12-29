package backend

import (
	"context"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/stardustagi/TopLib/libs/databases"
	"github.com/stardustagi/TopLib/libs/logs"
	"github.com/stardustagi/TopLib/libs/server"
	"github.com/stardustagi/TopLib/utils"
	"github.com/stardustagi/TopModelsNode/models"
	"go.uber.org/zap"
)

type Application struct {
	ctx      context.Context
	backend  *server.Backend
	logger   *zap.Logger
	config   server.HttpServerConfig
	wsConfig server.HttpWebSocketConfig
	manager  server.IClientManager
	upgrader websocket.Upgrader
}

func NewApplication(configBytes, wsConfigBytes []byte) *Application {
	config, err := utils.Bytes2Struct[server.HttpServerConfig](configBytes)
	if err != nil {
		panic(err)
	}
	wsConfig, err := utils.Bytes2Struct[server.HttpWebSocketConfig](wsConfigBytes)
	if err != nil {
		panic(err)
	}
	b, err := server.NewBackend(configBytes)
	return &Application{
		ctx:      context.Background(),
		config:   config,
		wsConfig: wsConfig,
		logger:   logs.GetLogger("HttpBackend"),
		manager:  server.NewClientManager(logs.GetLogger("clientManager")),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  wsConfig.ReadBufferSize,
			WriteBufferSize: wsConfig.WriteBufferSize,
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for simplicity
			},
		},
		backend: b,
	}
}

func (h *Application) Start() {
	h.logger.Info("Starting HttpBackend")
	// 同步数据库
	h.syncDatabaseSchema()
	go h.manager.Start()
	go func() {
		if err := h.backend.Start(); err != nil {
			h.logger.Error("backend.Start error", zap.Error(err))
		}
	}()
}

func (h *Application) Stop() {
	h.logger.Info("Stopping HttpBackend")
	h.manager.Stop()
	h.backend.Stop()
}

func (h *Application) AddGroup(group string, middleware ...echo.MiddlewareFunc) {
	h.backend.AddGroup(group, middleware...)
}

func (h *Application) AddPostHandler(group string, handler server.IHandler) {
	h.backend.AddPostHandler(group, handler)
}

func (h *Application) AddGetHandler(group string, handler server.IHandler) {
	h.backend.AddGetHandler(group, handler)
}

func (h *Application) AddNativeHandler(method, path string, handler echo.HandlerFunc) {
	h.backend.AddNativeHandler(method, path, handler)
}

func (h *Application) HandleWebSocket() {
	//ws := server.NewHandler(
	//	"ws",
	//	[]string{"websocket"},
	//	func(ctx echo.Context, req requests.DefaultWsRequest, resp responses.DefaultWsResponse) error {
	//		conn, err := h.upgrader.Upgrade(ctx.Response(), ctx.Request(), nil)
	//		if err != nil {
	//			return err
	//		}
	//		userId := ctx.Request().Header.Get("User-Id")
	//		if userId == "" {
	//			return echo.NewHTTPError(http.StatusUnauthorized, "User ID is required")
	//		}
	//		sessionId := uuid.GetUuidString()
	//		logger := logs.GetLogger("websocketClient")
	//		handler := wshandler.NewLLMModelServiceHandler()
	//		client := server.NewClient(
	//			userId,
	//			sessionId,
	//			conn,
	//			codec.NewJsonCodec(),
	//			logger,
	//			ctx.Request().Context(),
	//			handler,
	//			h.manager,
	//		)
	//		h.manager.RegisterClient(client)
	//		go client.Listen()
	//		return nil
	//	},
	//)
	//h.backend.AddHandler("GET", "/ws", ws)
	//h.logger.Info("WebSocket handler registered")
}

func (h *Application) syncDatabaseSchema() {
	h.logger.Info("Syncing database schema...")
	modelList := []interface{}{
		&models.Company{},
		&models.InvitationCode{},
		&models.LlmUsageReport{},
		&models.ModelsInfo{},
		&models.ModelsProvider{},
		&models.NodeModelsInfoMaps{},
		&models.NodeUsers{},
		&models.Nodes{},
		&models.SystemConfig{},
		&models.UserAgentTokens{},
		&models.UserApiKeys{},
		&models.UserConsumeRecord{},
		&models.UserModelsInfos{},
		&models.UserPayLog{},
		&models.UserWallet{},
		&models.Users{},
		&models.UsersKey{},
		&models.UserConsumeDetailText{},
		&models.UserConsumeDetailImage{},
		&models.UserConsumeDetailVideo{},
		&models.ModelsTieredPricing{},
	}
	dbDao := databases.GetDao()
	if err := dbDao.Native().Sync2(modelList...); err != nil {
		h.logger.Error("Database schema sync failed", zap.Error(err))
		panic("Database schema sync failed: " + err.Error())
	}
	h.logger.Info("Database schema synced successfully")
}

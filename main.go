package main

import (
	"github.com/stardustagi/TopLib/libs/conf"
	"github.com/stardustagi/TopLib/libs/databases"
	"github.com/stardustagi/TopLib/libs/logs"
	"github.com/stardustagi/TopLib/libs/redis"
	"github.com/stardustagi/TopModelsNode/backend"
	message "github.com/stardustagi/TopModelsNode/backend/services/nats"
	bill "github.com/stardustagi/TopModelsNode/backend/services/node_bill"
	llm "github.com/stardustagi/TopModelsNode/backend/services/node_llm"
	users "github.com/stardustagi/TopModelsNode/backend/services/node_users"
	"github.com/stardustagi/TopModelsNode/constants"
	_ "github.com/stardustagi/TopModelsNode/docs"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title TopModelsNode Backend Service
// @version 1.0
// @description This is the backend services for TopModelsLogin, providing user management and LLM
// @host localhost:8081
// @BasePath /api
func main() {
	conf.Init()
	loggerConfig := conf.Get("logger")
	constants.Init()
	logs.Init(loggerConfig)
	logger := logs.GetLogger("main")
	logger.Info("Init logs")

	// 初始化数据库
	_, _ = databases.Init(conf.Get("mysql"))
	logger.Info("Init mysql")

	// 初始化redis
	_, _ = redis.Init(conf.Get("redis"))
	logger.Info("Init redis")

	// 初始化nats
	nats := message.Init(conf.Get("nats"))
	nats.Start()
	logger.Info("Init nats")

	// 初始化backend
	app := backend.NewApplication(conf.Get("websrv"), conf.Get("websocket"))
	app.AddNativeHandler("GET", "/swagger/*", echoSwagger.WrapHandler)
	nodeService := llm.GetNodeHttpServiceInstance()
	nodeService.Start(app)
	defer nodeService.Stop()
	userService := users.GetNodeUsersHttpServiceInstance()
	userService.Start(app)
	nodeBill := bill.GetNodeHttpBillServiceInstance()
	nodeBill.Start(app)
	defer userService.Stop()

	app.Start()
	app.Stop()
}

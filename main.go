package main

import (
	"github.com/stardustagi/TopLib/libs/conf"
	"github.com/stardustagi/TopLib/libs/databases"
	"github.com/stardustagi/TopLib/libs/logs"
	"github.com/stardustagi/TopLib/libs/redis"
	"github.com/stardustagi/TopModelsNode/backend"
	llm "github.com/stardustagi/TopModelsNode/backend/services/node_llm"
	users "github.com/stardustagi/TopModelsNode/backend/services/node_users"
	"github.com/stardustagi/TopModelsNode/constants"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title TopModelsLogin Backend Service
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
	_, _ = databases.Init(conf.Get("mysql"))
	logger.Info("Init mysql")
	_, _ = redis.Init(conf.Get("redis"))
	logger.Info("Init redis")
	app := backend.NewApplication(conf.Get("websrv"), conf.Get("websocket"))
	app.AddNativeHandler("GET", "/swagger/*", echoSwagger.WrapHandler)
	nodeService := llm.GetNodeHttpServiceInstance()
	nodeService.Start(app)
	defer nodeService.Stop()
	userService := users.GetNodeUsersHttpServiceInstance()
	userService.Start(app)
	defer userService.Stop()

	app.Start()
	app.Stop()
}

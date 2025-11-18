package message

import (
	"os"

	"github.com/stardustagi/TopLib/libs/conf"
	"github.com/stardustagi/TopLib/libs/databases"
	"github.com/stardustagi/TopLib/libs/logs"
	tnats "github.com/stardustagi/TopLib/libs/nats"
	"github.com/stardustagi/TopLib/libs/redis"
)

func InitTestConfig() {
	err := os.Setenv("runConfig", "../../../config/prod.toml")
	if err != nil {
		panic("Failed to set environment variable: " + err.Error())
	}
	conf.Init()
	loggerConfig := conf.Get("logger")
	logs.Init(loggerConfig)
	logger := logs.GetLogger("main")
	logger.Info("Init logs")
	_, err = databases.Init(conf.Get("mysql"))
	if err != nil {
		panic("Failed to init database: " + err.Error())
	}
	logger.Info("Init mysql")
	_, err = redis.Init(conf.Get("redis"))
	if err != nil {
		panic("Failed to init redis: " + err.Error())
	}
	tnats.Init(conf.Get("nats"))
	logger.Info("Init nats")
}

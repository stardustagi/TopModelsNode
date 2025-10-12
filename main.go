package TopModelsNode

import (
	"github.com/stardustagi/TopLib/libs/conf"
	"github.com/stardustagi/TopLib/libs/databases"
	"github.com/stardustagi/TopLib/libs/logs"
	"github.com/stardustagi/TopLib/libs/redis"
	message "github.com/stardustagi/TopModelsLogin/backend/services/nats"
	"github.com/stardustagi/TopModelsLogin/constants"
)

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
	message.Init(conf.Get("nats"))
}

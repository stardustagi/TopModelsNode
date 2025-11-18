package message

import (
	"github.com/nats-io/nats.go"
	"github.com/stardustagi/TopLib/libs/databases"
	"github.com/stardustagi/TopLib/libs/logs"
	"go.uber.org/zap"
)

type MProcess struct {
	dbCon  databases.BaseDao
	logger *zap.Logger
}

func NewMProcess() *MProcess {
	dbInterface := databases.GetDao()
	if dbInterface == nil {
		return nil
	}
	return &MProcess{
		dbCon:  dbInterface,
		logger: logs.GetLogger("MProcess"), // 初始化logger
	}
}

// ProcessTokenUsage 得到订阅消息后的处理
func (pt *MProcess) ProcessTokenUsage(msg *nats.Msg) {
	if pt.logger != nil {
		pt.logger.Info("Received message on ", zap.String("subject", msg.Subject))
		pt.logger.Info("revice : ", zap.String("data", string(msg.Data)))
	}
}

// ProcessUsageReport 得到订阅消息后的处理
func (pt *MProcess) ProcessUsageReport(msg *nats.Msg) {
	if pt.logger != nil {
		pt.logger.Info("Received message on ", zap.String("subject", msg.Subject))
		pt.logger.Info("revice : ", zap.String("data", string(msg.Data)))
	}
}

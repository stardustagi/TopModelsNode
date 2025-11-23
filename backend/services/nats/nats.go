package message

import (
	"sync"
	"time"

	"github.com/stardustagi/TopLib/libs/logs"
	"github.com/stardustagi/TopLib/libs/nats"
	"go.uber.org/zap"
)

type NatsQueue struct {
	url       string
	conn      *nats.NatsConnection
	config    *nats.NatsConfig
	mp        *MProcess
	isRun     bool
	method    string
	csName    string
	pbName    string
	useStream bool
	logger    *zap.Logger
}

var (
	natsQueueInstance IMessage
	once              sync.Once
)

func Init(config []byte) IMessage {
	once.Do(func() {
		nats.Init(config)
		t := "publisher"
		natsQueueInstance = NewNatsQueue(t)
	})
	return natsQueueInstance
}

func GetNatsQueueInstance() IMessage {
	if natsQueueInstance == nil {
		panic("natsQueueInstance is nil, please call Init first")
	}
	return natsQueueInstance
}

func NewNatsQueue(t string) IMessage {
	natsManager := nats.GetNatsManager()
	natsCon, ok := natsManager.GetClient(t)
	if !ok {
		panic("nats connection not found")
	}
	config := natsCon.GetConfig()
	mp := NewMProcess()
	natsQueue := &NatsQueue{url: config.Url,
		conn:      natsCon,
		isRun:     false,
		mp:        mp,
		useStream: config.UseStream,
		config:    config,
		logger:    logs.GetLogger("natsQueue"),
	}
	switch config.Type {
	case "consumer":
		natsQueue.csName = config.Name
	case "publisher":
		natsQueue.pbName = config.Name
	}
	return natsQueue
}

func (n *NatsQueue) Start() {
	go n.conn.Start()
}

func (n *NatsQueue) Stop() {
	n.logger.Info("close nats connection!")
	err := n.conn.StopAllSubscriptions()
	if err != nil {
		n.logger.Error("Failed to stop all subscriptions", zap.Error(err))
	}
	n.conn.Stop()
}

func (n *NatsQueue) Publish(subject string, msg []byte) bool {
	if err := n.conn.Publish(subject, msg); err != nil {
		n.logger.Error("Failed to publish message", zap.String("subject", subject), zap.Error(err))
		return false
	}
	n.logger.Info("Published message", zap.String("subject", subject), zap.ByteString("msg", msg))
	return true
}

func (n *NatsQueue) PublisherStreamAsync(subject string, msg []byte) bool {
	if err := n.conn.PublishAsync(subject, msg); err != nil {
		n.logger.Error("Failed to publish async message", zap.String("subject", subject), zap.Error(err))
		return false
	}
	n.logger.Info("Published async message", zap.String("subject", subject), zap.ByteString("msg", msg))
	return true
}

func (n *NatsQueue) AddSubscriptionWithName(subject string, handler func(msg *nats.Msg)) bool {
	err := n.subscribe(subject, handler)
	if err != nil {
		n.logger.Error("subject 订阅失败: ", zap.String("subject", subject), zap.Error(err))
		return false
	}
	return true
}

func (n *NatsQueue) subscribe(subject string, handler func(msg *nats.Msg)) error {
	// 如果csName为空，使用默认值
	name := "default"
	if n.config.Type == "publisher" {
		name = n.pbName
	} else {
		name = n.csName
	}

	// 直接调用而不是使用goroutine，确保订阅立即建立
	err := n.conn.StartSubscription(subject, name, handler)
	if err != nil {
		n.logger.Error("Failed to start subscription", zap.String("subject", subject), zap.Error(err))
		return err
	}

	// 添加延迟确保订阅完全建立
	time.Sleep(1 * time.Second)
	return nil
}

func (n *NatsQueue) Unsubscribe(subject string) error {
	return n.conn.StopSubscription(subject)
}

func (n *NatsQueue) UnsubscribeAll() error {
	return n.conn.StopAllSubscriptions()
}

func (n *NatsQueue) GetNatsConn() *nats.NatsConnection {
	return n.conn
}

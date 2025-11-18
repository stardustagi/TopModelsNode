package message

import (
	"github.com/stardustagi/TopLib/libs/nats"
)

type IMessage interface {
	Start()
	Stop()
	Publish(subject string, msg []byte) bool
	PublisherStreamAsync(subject string, msg []byte) bool
	AddSubscriptionWithName(subject string, handler func(msg *nats.Msg)) bool
	Unsubscribe(sub string) error
	UnsubscribeAll() error
	GetNatsConn() *nats.NatsConnection
}

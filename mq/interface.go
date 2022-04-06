package mq

import (
	"context"
	"sync"
	"time"

	"github.com/glory-go/glory/config"
	"github.com/glory-go/glory/log"
)

var (
	mqTypeMap sync.Map
)

type MQMsgHandler func(context.Context, []byte) error

type MQServiceFactory func(model config.ChannelType, rawConfig map[string]string) (MQService, error)

type MQService interface {
	Connect() error
	Send(topic string, msg []byte) (msgID string, err error)
	DelaySend(topic string, msg []byte, handleTime time.Time) (msgID string, err error)
	Publish(topic string, msg []byte) (msgID string, err error)
	RegisterHandler(topic string, handler MQMsgHandler)
}

func RegisterMQType(mqType string, mqFactory MQServiceFactory) {
	_, ok := mqTypeMap.LoadOrStore(mqType, mqFactory)
	if ok {
		log.Warnf("mq type [%s] has already been registered, now replace earlier one")
	}
}

func getMQFactory(mqType string) (MQServiceFactory, bool) {
	v, ok := mqTypeMap.Load(mqType)
	if !ok || v == nil {
		return nil, false
	}
	factory, ok := v.(MQServiceFactory)
	if !ok || factory == nil {
		return nil, false
	}
	return factory, true
}

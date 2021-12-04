package mq

import (
	"sync"
	"time"

	"github.com/glory-go/glory/log"
)

var (
	mqTypeMap sync.Map
)

type MQMsgHandler func([]byte) error

type MQServiceFactory func(rawConfig map[string]string) (MQService, error)

type MQService interface {
	Connect() error // 根据配置进行mq的连接操作
	Send(topic string, msg []byte) (msgID string, err error)
	DelaySend(topic string, msg []byte, handleTime time.Time) (msgID string, err error)
	RegisterHandler(topic string, handler MQMsgHandler) error
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

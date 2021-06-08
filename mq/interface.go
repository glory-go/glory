package mq

import (
	"time"

	"github.com/glory-go/glory/config"
)

// TODO: 抽象出收到的消息的通用方法
type MQMsgHandler func([]byte) error

type MQService interface {
	loadConfig(conf *config.MQConfig) error
	// 返回MsgID
	Send(topic string, msg []byte) (string, error)
	DelaySend(topic string, msg []byte, handleTime time.Time) (string, error)
	// Handler 函数的调用会阻塞进程
	Handler(topic string, handler MQMsgHandler) error
}

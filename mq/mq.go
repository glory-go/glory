package mq

import (
	"errors"
	"fmt"

	"github.com/glory-go/glory/config"
	"github.com/glory-go/glory/log"
	"github.com/streadway/amqp"
)

func init() {
	defaultRabbitMQHandler = newRabbitMQDefaultHandler()
	defaultRabbitMQHandler.setup(config.GlobalServerConf.MQConfig)
}

type RabbitMQHandler struct {
	mqServicesMap map[string]*RabbitMQService
}

var defaultRabbitMQHandler *RabbitMQHandler

func (mh *RabbitMQHandler) setup(conf map[string]*config.MQConfig) {
	for k, v := range conf {
		tempService := newRabbitMQService()
		tempService.loadConfig(*v)
		if err := tempService.connect(); err != nil {
			log.Error("opendb with key = ", k, " error = ", err)
			continue
		}
		mh.mqServicesMap[k] = tempService
	}
}

func newRabbitMQDefaultHandler() *RabbitMQHandler {
	return &RabbitMQHandler{
		mqServicesMap: make(map[string]*RabbitMQService),
	}
}

func NewRabbitMQClient(rbmqServiceName string) (*amqp.Connection, error) {
	service, ok := defaultRabbitMQHandler.mqServicesMap[rbmqServiceName]
	if !ok {
		log.Error("rabbitmq service name = ", rbmqServiceName, " not registered in config")
		return nil, errors.New("rabbitmq service name = " + rbmqServiceName + " not registered in config")
	}
	return service.conn, nil
}

func GetService(rbmqServiceName string) (*RabbitMQService, bool) {
	s, ok := defaultRabbitMQHandler.mqServicesMap[rbmqServiceName]
	return s, ok
}

func StartOnMQMsgHandler(rbmqServiceName, channelName string, hanlder MQMsgHandler) error {
	service, ok := defaultRabbitMQHandler.mqServicesMap[rbmqServiceName]
	if !ok {
		log.Error("rabbitmq service name = ", rbmqServiceName, " not registered in config")
		return errors.New("rabbitmq service name = " + rbmqServiceName + " not registered in config")
	}
	if err := service.startOnMsgHandler(channelName, hanlder); err != nil {
		log.Error("rabbitmq service send to ", channelName, " queue err = ", err)
		return err
	}
	return nil
}

func GetMQInstance(name string) MQService {
	srv, ok := mqInstance[name]
	if !ok {
		panic(fmt.Sprintf("mq with name %s not found", name))
	}
	return srv
}

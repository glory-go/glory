package mq

import (
	"errors"
	"fmt"

	"github.com/glory-go/glory/config"
	"github.com/glory-go/glory/log"
	"github.com/streadway/amqp"
)

// SendMQMessage 直接发送一条消息
func SendMQMessage(rbmqServiceName, chName string, msgData []byte) error {
	service, ok := defaultRabbitMQHandler.mqServicesMap[rbmqServiceName]
	if !ok {
		log.Error("rabbitmq service name = ", rbmqServiceName, "not registered in config")
		return errors.New("rabbitmq service name = " + rbmqServiceName + "not registered in config")
	}
	if err := service.send(chName, msgData, nil); err != nil {
		log.Error("rabbitmq service send to ", chName, " queue err = ", err)
		return err
	}
	return nil
}

// SendDelayMsg 发送延迟消息，单位为ms
func SendDelayMsg(rbmqServiceName, chName string, msgData []byte, time int64) error {
	service, ok := defaultRabbitMQHandler.mqServicesMap[rbmqServiceName]
	if !ok {
		log.Error("rabbitmq service name = ", rbmqServiceName, "not registered in config")
		return errors.New("rabbitmq service name = " + rbmqServiceName + "not registered in config")
	}
	if err := service.send(chName, msgData, &time); err != nil {
		log.Error("rabbitmq service send delay msg to ", chName, " queue err = ", err)
		return err
	}
	return nil
}

func (ms *RabbitMQService) send(chName string, msg []byte, delayTime *int64) error {
	chConfig, ok := ms.conf.Channels[chName]
	if !ok {
		return fmt.Errorf("channel %v not found", chName)
	}
	// 检查格式
	if chConfig.Type == config.Delay || chConfig.Type == config.PubSub {
		if chConfig.ExchangeName == "" {
			return fmt.Errorf("invalid channel %v, delay queue must has exchange name", chName)
		}
	}
	// 建立连接
	ch, err := ms.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// 初始化queue
	q, err := ch.QueueDeclare(
		chConfig.QueueName, // name
		true,               // durable
		false,              // delete when unused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		log.Errorf("init channel queue %v for %v meets error: %v", chConfig.QueueName, chName, err)
		return err
	}

	// 初始化exchange
	if chConfig.ExchangeName != "" {
		args := amqp.Table{}
		// TODO: 未考虑direct模式下应该如何定义
		exchangeType := "fanout"
		if chConfig.Type == config.Delay {
			args["x-delayed-type"] = "direct"
			exchangeType = "x-delayed-message"
		}
		if err := ch.ExchangeDeclare(
			chConfig.ExchangeName, // name
			exchangeType,          // type
			true,                  // durable
			false,                 // auto-deleted
			false,                 // internal
			false,                 // no-wait
			args,                  // arguments
		); err != nil {
			log.Errorf("init channel exchange %v for %v meets error", chConfig.ExchangeName, chName)
			return err
		}
	}

	headers := amqp.Table{}
	routingKey := q.Name
	if chConfig.Type == config.Delay {
		if delayTime == nil {
			tmp := int64(0)
			delayTime = &tmp
		}
		headers["x-delay"] = *delayTime
		routingKey = ""
	}
	if err = ch.Publish(
		chConfig.ExchangeName, // exchange
		routingKey,            // routing key
		false,                 // mandatory
		false,                 // immediate
		amqp.Publishing{
			Headers:     headers,
			ContentType: "text/plain",
			Body:        msg,
		}); err != nil {
		return err
	}
	return nil
}

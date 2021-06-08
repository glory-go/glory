package mq

import (
	"fmt"

	"github.com/glory-go/glory/log"

	"github.com/streadway/amqp"

	"github.com/glory-go/glory/config"
	_ "github.com/go-sql-driver/mysql"
)

// RabbitMQService 保存多个redis的库
type RabbitMQService struct {
	db   map[string]*amqp.Delivery
	conf config.MQConfig
	conn *amqp.Connection
	ch   map[string]*amqp.Channel
}

func newRabbitMQService() *RabbitMQService {
	return &RabbitMQService{
		db: make(map[string]*amqp.Delivery, 8),
		ch: make(map[string]*amqp.Channel),
	}
}

func (ms *RabbitMQService) loadConfig(conf config.MQConfig) {
	ms.conf = conf
}

func (ms *RabbitMQService) connect() error {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s", ms.conf.Username, ms.conf.Password, ms.conf.Host, ms.conf.Port)
	conn, err := amqp.Dial(url)
	if err != nil {
		return err
	}
	ms.conn = conn
	return nil
}

func (ms *RabbitMQService) startOnMsgHandler(channelName string, handler MQMsgHandler) error {
	chConfig, ok := ms.conf.Channels[channelName]
	if !ok {
		return fmt.Errorf("channel %v not found", channelName)
	}
	// 检查格式
	if chConfig.Type == config.Delay || chConfig.Type == config.PubSub {
		if chConfig.ExchangeName == "" {
			return fmt.Errorf("invalid channel %v, delay queue must has exchange name", channelName)
		}
	}
	// 建立连接
	ch, err := ms.conn.Channel()
	if err != nil {
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
			log.Errorf("init channel exchange %v for %v meets error", chConfig.ExchangeName, channelName)
			return err
		}
	}

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
		log.Errorf("init channel queue %v for %v meets error: %v", chConfig.QueueName, channelName, err)
		return err
	}
	if chConfig.ExchangeName != "" {
		if err = ch.QueueBind(
			q.Name,                // queue name
			"",                    // routing key
			chConfig.ExchangeName, // exchange
			false,
			nil,
		); err != nil {
			log.Errorf("queue bind %v for %v meets error", chConfig.QueueName, channelName)
			return err
		}
	}
	// 消费消息
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return err
	}
	go func() {
		defer ch.Close()
		defer func() {
			if e := recover(); e != nil {
				log.Errorf("mq loop runtime err = %v", e)
			}
		}()
		for msg := range msgs {
			log.Debug("queue ", chConfig.QueueName, " receive msg ", string(msg.Body))
			if err := handler(msg.Body); err != nil {
				log.Warn("mq handler throw an error = ", err)
				if ms.conf.AutoACK {
					msg.Ack(false)
				}
			} else {
				msg.Ack(false)
			}
		}
	}()
	return nil
}

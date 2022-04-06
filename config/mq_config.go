package config

import (
	"fmt"

	"github.com/glory-go/glory/tools"
)

type ChannelType string

const (
	Direct ChannelType = "direct"
	PubSub ChannelType = "pub_sub"
	Delay  ChannelType = "delay"
)

type MQConfig struct {
	Type         string `yaml:"type"` // 使用的MQ类型
	ConfigSource string `yaml:"config_source"`

	Config map[string]string `yaml:"config"` // 由具体的mq实现决定其内容如何解析
}

type RabbitMQChannelConfig struct {
	QueueName    string      `yaml:"queue"`    // 可为空，代表生成一个no-durable的队列，名字由系统给定
	ExchangeName string      `yaml:"exchange"` // 可为空，若为空则不会使用exchange，而是往queue中直接发送；pubsub和delay必须指定
	Type         ChannelType `yaml:"type"`     // direct/pub_sub/delay
}

func (s *MQConfig) checkAndFix() {
	if err := tools.ReadFromEnvIfNeed(s); err != nil {
		fmt.Println("warn: RedisConfig checkAndFix failed with err = ", err)
	}
}

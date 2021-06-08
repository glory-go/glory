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
	// common
	Type string `yaml:"type"` // 使用的MQ类型

	ConfigSource string                           `yaml:"config_source"`
	Host         string                           `yaml:"host"`
	Port         string                           `yaml:"port"`
	Username     string                           `yaml:"username"`
	Password     string                           `yaml:"password"`
	AutoACK      bool                             `yaml:"auto_ack"`
	Channels     map[string]RabbitMQChannelConfig `yaml:"channels"` // key是这个队列我们指定的名称

	// aliyun
	Endpoint        string `yaml:"endpoint"`
	AccessKey       string `yaml:"access_key"`
	SecretKey       string `yaml:"secret_key"`
	InstanceID      string `yaml:"instance_id"`
	ConsumerGroupID string `yaml:"consumer_group_id"`
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

package config

import (
	"fmt"

	"github.com/glory-go/glory/tools"
)

type ChannelType string

const (
	Direct ChannelType = "normal"
	PubSub ChannelType = "pub_sub"
)

type MQConfig struct {
	Type         string      `yaml:"type"` // 使用的MQ类型
	Mod          ChannelType `yaml:"mod"`
	ConfigSource string      `yaml:"config_source"`

	Config map[string]string `yaml:"config"` // 由具体的mq实现决定其内容如何解析
}

func (s *MQConfig) checkAndFix() {
	if err := tools.ReadFromEnvIfNeed(s); err != nil {
		fmt.Println("warn: RedisConfig checkAndFix failed with err = ", err)
	}
}

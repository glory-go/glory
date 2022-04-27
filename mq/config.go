package mq

import (
	"fmt"
)

import (
	"github.com/glory-go/glory/config"
	"github.com/glory-go/glory/tools"
)

var mqInstance = map[string]MQService{}

func loadConfig(mqConfigs map[string]*config.MQConfig) {
	for name, config := range mqConfigs {
		if config.ConfigSource == "env" {
			config.Config = tools.ReadMapConfigFromEnv(config.Config)
		}
		factory, ok := getMQFactory(config.Type)
		if !ok {
			panic(fmt.Sprintf("mq with type %s not register", config.Type))
		}
		instance, err := factory(config.Config)
		if err != nil {
			panic(fmt.Sprintf("fail to load config for mq %s, err is %v", name, err))
		}
		mqInstance[name] = instance
	}
}

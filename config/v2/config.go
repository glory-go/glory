package config

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"
)

var (
	ConfigInstance *viper.Viper
	builderMap     sync.Map
)

type ComponentBuilder func(name string, conf map[string]interface{}) error

func RegisterConfig(name string, builder ComponentBuilder) error {
	funcName := "RegisterConfig"
	fmt.Printf("[%s] start to parse config for [%s]", funcName, name)

	if _, ok := builderMap.Load(name); ok {
		return fmt.Errorf("builder with name %v already registered", name)
	}
	builderMap.Store(name, builder)
	return nil
}

func InitModules(modules ...string) error {
	if !inited {
		return fmt.Errorf("config not inited yet")
	}
	for _, module := range modules {
		builderI, ok := builderMap.Load(module)
		if !ok {
			return fmt.Errorf("builder for config %v not registered", module)
		}
		builder, ok := builderI.(ComponentBuilder)
		if !ok {
			return fmt.Errorf("builder %v got invalid type: %T", module, builderI)
		}

		// 将配置中从config_center中获取的部分进行解析
		conf := make(map[string]map[string]interface{})
		if err := ConfigInstance.UnmarshalKey(module, &conf); err != nil {
			return err
		}
		for name, subConf := range conf {
			if err := ReplaceStringValueFromConfigCenter(subConf); err != nil {
				return err
			}
			if err := builder(name, subConf); err != nil {
				return err
			}
		}
	}
	return nil
}

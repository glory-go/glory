package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var (
	ConfigInstance *viper.Viper
)

type ComponentBuilder func(name string, conf map[string]interface{}) error

func RegisterConfig(name string, builder ComponentBuilder) error {
	funcName := "RegisterConfig"
	fmt.Printf("[%s] start to parse config for [%s]", funcName, name)
	// 将配置中从config_center中获取的部分进行解析
	conf := make(map[string]map[string]interface{})
	if err := ConfigInstance.UnmarshalKey(name, &conf); err != nil {
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
	return nil
}

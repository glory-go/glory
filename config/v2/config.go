package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var (
	ConfigInstance *viper.Viper
)

type ComponentBuilder func(conf map[string]interface{}) error

func RegisterConfig(name string, builder ComponentBuilder) error {
	funcName := "RegisterConfig"
	fmt.Printf("[%s] start to parse config for [%s]", funcName, name)
	// 将配置中从config_center中获取的部分进行解析
	conf := ConfigInstance.GetStringMap(name)
	if err := ReplaceStringValueFromConfigCenter(conf); err != nil {
		return err
	}
	return builder(conf)
}

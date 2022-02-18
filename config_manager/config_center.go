package configmanager

import (
	"fmt"

	"github.com/spf13/viper"
)

var (
	configCenter              ConfigCenter
	configCenterViperInstance *viper.Viper
)

func RegisterConfigBuilder(builder ConfigCenterBuilder) error {
	if configCenter != nil {
		return fmt.Errorf("only allow one config center builder")
	}
	// 初始化配置中心
	var err error
	configCenter, err = builder(configCenterViperInstance)
	if err != nil {
		return err
	}

	return nil
}

func GetConfigCenter() (ConfigCenter, error) {
	if configCenter == nil {
		return configCenter, fmt.Errorf("config center not register")
	}
	return configCenter, nil
}

package configmanager

import (
	"fmt"

	"github.com/spf13/viper"
)

var (
	configCenter              ConfigCenter
	configCenterViperInstance *viper.Viper
)

func RegisterConfigBuilder(configCenterName string, builder ConfigCenterBuilder) error {
	if configCenter != nil {
		return fmt.Errorf("only allow one config center builder")
	}
	config := configCenterViperInstance.GetStringMapString(configCenterName)
	// TODO: 从env中替换内容
	// 初始化配置中心
	var err error
	configCenter, err = builder(config)
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

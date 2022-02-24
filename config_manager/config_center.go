package configmanager

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"
)

var (
	configCenterMap           sync.Map
	configCenterViperInstance *viper.Viper
)

func RegisterConfigBuilder(configCenterName string, builder ConfigCenterBuilder) error {
	if _, ok := configCenterMap.Load(configCenterName); ok {
		return fmt.Errorf("config center %s already registered", configCenterName)
	}
	config := configCenterViperInstance.GetStringMapString(configCenterName)
	ReadFromEnvIfNeed(config)
	// 初始化配置中心
	var err error
	configCenter, err := builder(config)
	if err != nil {
		return err
	}
	configCenterMap.Store(configCenterName, configCenter)

	return nil
}

func GetConfigCenter(configCenterName string) (ConfigCenter, error) {
	configCenter, ok := configCenterMap.Load(configCenterName)
	if !ok || configCenter == nil {
		return nil, fmt.Errorf("config center %s not registered", configCenterName)
	}
	return configCenter.(ConfigCenter), nil
}

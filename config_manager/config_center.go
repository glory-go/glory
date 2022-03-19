package configmanager

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"
)

var (
	configCenterMap           sync.Map
	configCenterBuilderMap    sync.Map
	configCenterViperInstance *viper.Viper
)

// ConfigCenterBuilder实现了配置中心的初始化工作，conf文件需以yaml形式提供
type ConfigCenterBuilder func(conf map[string]string) (ConfigCenter, error)

func RegisterConfigBuilder(configCenterName string, builder ConfigCenterBuilder) error {
	funcName := "RegisterConfig"
	fmt.Printf("[%s] start to register config center [%s]", funcName, configCenterName)

	if _, ok := configCenterBuilderMap.Load(configCenterName); ok {
		return fmt.Errorf("config center %s already registered", configCenterName)
	}
	configCenterBuilderMap.Store(configCenterName, builder)

	return nil
}

func InitConfigCenter(configCenterName ...string) error {
	if !inited {
		return fmt.Errorf("config center not inited")
	}
	for _, name := range configCenterName {
		builderI, ok := configCenterBuilderMap.Load(configCenterName)
		if !ok {
			return fmt.Errorf("builder for config center %s not registered", name)
		}

		builder, ok := builderI.(ConfigCenterBuilder)
		if !ok {
			return fmt.Errorf("invalid builder type for config center %s, type: %T", name, builderI)
		}
		config := configCenterViperInstance.GetStringMapString(name)
		ReadFromEnvIfNeed(config)
		// 初始化配置中心
		var err error
		configCenter, err := builder(config)
		if err != nil {
			return err
		}
		configCenterMap.Store(configCenterName, configCenter)
	}
	return nil
}

func GetConfigCenter(configCenterName string) (ConfigCenter, error) {
	if !inited {
		return nil, fmt.Errorf("config center not register yet")
	}
	configCenter, ok := configCenterMap.Load(configCenterName)
	if !ok || configCenter == nil {
		return nil, fmt.Errorf("config center %s not registered", configCenterName)
	}
	return configCenter.(ConfigCenter), nil
}

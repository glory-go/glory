package configmanager

type ConfigCenter interface {
	// LoadConfig 从配置中心获取配置
	LoadConfig(key string) (string, error)
	// SyncConfig 持续从配置中心获取配置，当配置发生变更时，将会修改value指向的值
	SyncConfig(key string, value *string) error
}

type ConfigCenterBuilder func(map[string]string) (ConfigCenter, error)

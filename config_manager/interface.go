package configmanager

type ConfigCenter interface {
	// LoadConfig 从配置中心获取配置
	LoadConfig(key, group string) (string, error)
	// SyncConfig 持续从配置中心获取配置，当配置发生变更时，将会修改value指向的值
	SyncConfig(key, group string, value *string, cancel <-chan struct{}) error
}

// ConfigCenterBuilder实现了配置中心的初始化工作，conf文件需以yaml形式提供
type ConfigCenterBuilder func(conf map[string]string) (ConfigCenter, error)

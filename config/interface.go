package config

type ConfigCenter interface {
	Name() string
	// Init 为配置中心需实现的初始化函数，该函数只会被调用一次
	Init(config map[string]interface{}) error
}

type Component interface {
	Name() string
	// Init 为组件提供方需实现的初始化函数，该函数只会被调用一次
	Init(config map[string]interface{}) error
}

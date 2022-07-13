package config

//go:generate mockgen -source interface.go -destination mock/interface.go

type ConfigCenter interface {
	Name() string
	// Init 为配置中心需实现的初始化函数，该函数只会被调用一次
	Init(config map[string]any) error
	// Get 从配置中心读取配置，params为待传入的参数，一般来源于用户的配置文件
	Get(params ...any) (string, error)
}

type Component interface {
	Name() string
	// Init 为组件提供方需实现的初始化函数，该函数只会被调用一次
	Init(config map[string]any) error
}

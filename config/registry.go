package config

import "sync"

var (
	configCenterRegistry sync.Map
	componentRegistry    sync.Map
)

// 配置中心相关方法

func RegisterConfigCenter(center ConfigCenter) {}

func iterConfigRegistry(func(name string, center ConfigCenter) error) {}

// 组件相关方法

func RegisterComponent(center Component) {}

func iterComponentRegistry(func(name string, component Component) error) {}

package config

import (
	"fmt"
	"sync"
)

var (
	configCenterRegistry sync.Map
	componentRegistry    sync.Map
)

/*** 配置中心相关方法 ***/

func RegisterConfigCenter(center ConfigCenter) {
	if center == nil {
		panic("register nil config center")
	}
	if keepedConfigCenterName.Contains(center.Name()) {
		panic(fmt.Sprintf("config center name [%s] was kept by glory", center.Name()))
	}
	_, ok := configCenterRegistry.LoadOrStore(center.Name(), center)
	if ok {
		panic(fmt.Sprintf("config center [%s] already register before", center.Name()))
	}
}

// registerInnerConfigCenter 用于给glory内部实现进行配置中心的注册
func registerInnerConfigCenter(centers ...ConfigCenter) {
	for _, center := range centers {
		configCenterRegistry.LoadOrStore(center.Name(), center)
	}
}

func GetConfigCenter(name string) (ConfigCenter, error) {
	val, ok := configCenterRegistry.Load(name)
	if !ok || val == nil {
		return nil, fmt.Errorf("config center %s not register", name)
	}
	center := val.(ConfigCenter)
	return center, nil
}

func iterConfigRegistry(f func(name string, center ConfigCenter) error) {
	configCenterRegistry.Range(func(key, value any) bool {
		if err := f(key.(string), value.(ConfigCenter)); err != nil {
			panic(err)
		}
		return true
	})
}

/*** 组件相关方法 ***/

func RegisterComponent(component Component) {
	if component == nil {
		panic("register nil component")
	}
	if keepedComponentName.Contains(component.Name()) {
		panic(fmt.Sprintf("component name [%s] was kept by glory", component.Name()))
	}
	// 初始化配置中心，确保配置内容完成了解析
	Init()
	componentRegistry.Store(component.Name(), component)
	rawConfig := make(map[string]any)
	getConfig(component.Name(), &rawConfig)
	if err := component.Init(rawConfig); err != nil {
		panic(err)
	}
}

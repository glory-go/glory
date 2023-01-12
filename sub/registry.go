package sub

import (
	"log"

	"github.com/mitchellh/mapstructure"
)

func (s *subSrv) RegisterSubProvider(provider SubProvider) {
	if provider == nil {
		panic("register nil provider")
	}
	register()
	raw, ok := rawConf[provider.Name()]
	if !ok {
		log.Printf("config of sub %s not found", provider.Name())
		return
	}
	// 将配置转为map[string]any形式
	providerConfig := make(map[string]any)
	if err := mapstructure.Decode(raw, &providerConfig); err != nil {
		panic(err)
	}
	// 调用服务本身的初始化方法完成初始化
	if err := provider.Init(providerConfig); err != nil {
		panic(err)
	}
	s.subProviderRegistry.Store(provider.Name(), provider)
}

func (s *subSrv) iterSubProviderRegistry(f func(name string, provider SubProvider) error) error {
	var resErr error
	s.subProviderRegistry.Range(func(key, value any) bool {
		if err := f(key.(string), value.(SubProvider)); err != nil {
			resErr = err
			return false
		}
		return true
	})

	return resErr
}

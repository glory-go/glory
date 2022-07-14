package sub

import (
	"fmt"
	"sync"

	"github.com/mitchellh/mapstructure"
)

const (
	SubSrvName = "sub"
)

type subSrv struct {
	subProviderRegistry sync.Map
}

var (
	sub  *subSrv
	once sync.Once
)

func GetSub() *subSrv {
	once.Do(func() {
		sub = &subSrv{}
	})

	return sub
}

func (s *subSrv) GetSubProvider(name string) SubProvider {
	raw, ok := s.subProviderRegistry.Load(name)
	if !ok {
		return nil
	}
	provider, ok := raw.(SubProvider)
	if !ok {
		return nil
	}
	return provider
}

func (s *subSrv) Name() string { return SubSrvName }

func (s *subSrv) Init(config map[string]any) error {
	if err := s.iterSubProviderRegistry(func(name string, provider SubProvider) error {
		raw, ok := config[name]
		if !ok {
			return fmt.Errorf("config of sub %s not found", name)
		}
		// 将配置转为map[string]any形式
		providerConfig := make(map[string]any)
		if err := mapstructure.Decode(raw, &providerConfig); err != nil {
			return err
		}
		// 调用服务本身的初始化方法完成初始化
		if err := provider.Init(providerConfig); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (s *subSrv) Run() error {
	if err := s.iterSubProviderRegistry(func(name string, provider SubProvider) error {
		return provider.Run()
	}); err != nil {
		return err
	}

	return nil
}

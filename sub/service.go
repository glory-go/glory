package sub

import (
	"sync"
)

const (
	SubSrvName = "sub"
)

type subSrv struct {
	subProviderRegistry sync.Map
}

var (
	sub     *subSrv
	rawConf map[string]any

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
	rawConf = config
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

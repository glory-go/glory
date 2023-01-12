package service

import (
	"sync"
)

const (
	ServiceComponentName = "service"
)

var (
	srv     *serviceComponent
	rawConf map[string]any

	srvInitOnce sync.Once
)

func GetService() *serviceComponent {
	srvInitOnce.Do(func() {
		srv = &serviceComponent{
			serviceRegistry: sync.Map{},
		}
	})

	return srv
}

type serviceComponent struct {
	serviceRegistry sync.Map
}

func (s *serviceComponent) Name() string {
	return ServiceComponentName
}

func (s *serviceComponent) Init(config map[string]any) error {
	rawConf = config
	return nil
}

package service

import (
	"fmt"
	"sync"

	"github.com/mitchellh/mapstructure"
)

const (
	ServiceComponentName = "service"
)

var (
	srv         *serviceComponent
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
	if err := s.iterServiceRegistry(func(name string, srv Service) error {
		raw, ok := config[name]
		if !ok {
			return fmt.Errorf("config of service %s not found", name)
		}
		// 将配置转为map[string]any形式
		srvConfig := make(map[string]any)
		if err := mapstructure.Decode(raw, &srvConfig); err != nil {
			return err
		}
		// 调用服务本身的初始化方法完成初始化
		if err := srv.Init(srvConfig); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

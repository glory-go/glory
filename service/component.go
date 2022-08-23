package service

import (
	"log"
	"sync"

	mapset "github.com/deckarep/golang-set/v2"
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
			inited:          mapset.NewSet[string](),
		}
	})

	return srv
}

type serviceComponent struct {
	serviceRegistry sync.Map
	inited          mapset.Set[string] // 已经初始化的服务才会存储在这里
}

func (s *serviceComponent) Name() string {
	return ServiceComponentName
}

func (s *serviceComponent) Init(config map[string]any) error {
	if err := s.iterServiceRegistry(func(name string, srv Service) error {
		raw, ok := config[name]
		if !ok {
			log.Printf("config of service %s not found", name)
			return nil
		}
		s.inited.Add(name)
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

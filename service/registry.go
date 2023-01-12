package service

import (
	"log"

	"github.com/mitchellh/mapstructure"
)

func (s *serviceComponent) RegisterService(srv Service) {
	if srv == nil {
		panic("register nil service")
	}
	register()
	raw, ok := rawConf[srv.Name()]
	if !ok {
		log.Printf("config of service %s not found", srv.Name())
		return
	}
	// 将配置转为map[string]any形式
	srvConfig := make(map[string]any)
	if err := mapstructure.Decode(raw, &srvConfig); err != nil {
		panic(err)
	}
	// 调用服务本身的初始化方法完成初始化
	if err := srv.Init(srvConfig); err != nil {
		panic(err)
	}
	s.serviceRegistry.Store(srv.Name(), srv)
}

func (s *serviceComponent) iterServiceRegistry(f func(name string, srv Service) error) error {
	var resErr error
	s.serviceRegistry.Range(func(key, value any) bool {
		if err := f(key.(string), value.(Service)); err != nil {
			resErr = err
			return false
		}
		return true
	})

	return resErr
}

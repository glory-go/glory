package service

import (
	"sync"

	"github.com/glory-go/glory/v2/config"
)

var (
	registerOnce sync.Once
)

func register() {
	registerOnce.Do(func() {
		config.Init()
		config.RegisterComponent(GetService())
	})
}

func (s *serviceComponent) Run() {
	wg := sync.WaitGroup{}
	s.iterServiceRegistry(func(name string, srv Service) error {
		if !s.inited.Contains(name) { // 未初始化的服务不调用Run方法
			return nil
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := srv.Run(); err != nil {
				panic(err)
			}
		}()
		return nil
	})
	wg.Wait()
}

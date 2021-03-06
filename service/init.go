package service

import (
	"sync"

	"github.com/glory-go/glory/v2/config"
)

func init() {
	config.RegisterComponent(GetService())
}

func (s *serviceComponent) Run() {
	wg := sync.WaitGroup{}
	s.iterServiceRegistry(func(name string, srv Service) error {
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

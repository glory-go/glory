package service

import (
	"fmt"
)

func (s *serviceComponent) RegisterService(srv Service) {
	if srv == nil {
		panic("register nil service")
	}
	_, ok := s.serviceRegistry.LoadOrStore(srv.Name(), srv)
	if ok {
		panic(fmt.Sprintf("service [%s] already register before", srv.Name()))
	}
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

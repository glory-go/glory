package sub

import "fmt"

func (s *subSrv) RegisterSubProvider(srv SubProvider) {
	if srv == nil {
		panic("register nil service")
	}
	_, ok := s.subProviderRegistry.LoadOrStore(srv.Name(), srv)
	if ok {
		panic(fmt.Sprintf("service [%s] already register before", srv.Name()))
	}
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

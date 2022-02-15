package nebula

import (
	"github.com/glory-go/glory/config"
	"github.com/glory-go/glory/log"
)

func init() {
	defaultMysqlHandler = newNebulaHandler()
	defaultMysqlHandler.setup(config.GlobalServerConf.NebulaConfigs)
}

type NebulaHandler struct {
	services map[string]*NebulaService
}

var defaultMysqlHandler *NebulaHandler

func (mh *NebulaHandler) setup(conf map[string]*config.NebulaConfig) {
	for k, v := range conf {
		tempService := newNebulaService()
		if err := tempService.openDB(*v); err != nil {
			log.Errorf("opendb with key = %s, err = %s", k, err)
			continue
		}
		mh.services[k] = tempService
	}
}

func newNebulaHandler() *NebulaHandler {
	return &NebulaHandler{
		services: make(map[string]*NebulaService),
	}
}

func GetService(serviceName string) (*NebulaService, bool) {
	s, ok := defaultMysqlHandler.services[serviceName]
	return s, ok
}

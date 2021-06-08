package k8s

import (
	"github.com/glory-go/glory/config"
	"github.com/glory-go/glory/log"
)

func init() {
	defaultK8SHandler = newK8SHandler()
	defaultK8SHandler.setup(config.GlobalServerConf.K8SConfig)
}

type K8SHandler struct {
	redisServices map[string]*K8SService
}

var defaultK8SHandler *K8SHandler

func (mh *K8SHandler) setup(conf map[string]*config.K8SConfig) {
	for k, v := range conf {
		tempService := newK8SService()
		if err := tempService.openDB(*v); err != nil {
			log.Error("opendb with key = ", k, "err")
			continue
		}
		mh.redisServices[k] = tempService
	}
}

func newK8SHandler() *K8SHandler {
	return &K8SHandler{
		redisServices: make(map[string]*K8SService),
	}
}

func GetService(k8sServiceName string) (*K8SService, bool) {
	s, ok := defaultK8SHandler.redisServices[k8sServiceName]
	return s, ok
}

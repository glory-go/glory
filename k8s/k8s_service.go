package k8s

import (
	"github.com/glory-go/glory/log"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/glory-go/glory/config"
)

// K8SService 保存该k8s的配置信息
type K8SService struct {
	conf   config.K8SConfig
	Client *kubernetes.Clientset
}

func (ms *K8SService) GetNamespace() string {
	return ms.conf.Namespace
}

func (ms *K8SService) GetDefaultUsername() string {
	return ms.conf.DefaultUsername
}
func (ms *K8SService) GetPVC(name string) string {
	return ms.conf.PVC[name]
}

func newK8SService() *K8SService {
	return &K8SService{}
}

func (ms *K8SService) loadConfig(conf config.K8SConfig) error {
	ms.conf = conf
	return nil
}

func (ms *K8SService) openDB(conf config.K8SConfig) error {
	if err := ms.loadConfig(conf); err != nil {
		log.Error("opendb error with err = ", err)
		return err
	}
	// 连接到k8s
	connConfig := &rest.Config{}
	if conf.ConnectType == config.InCluster {
		var err error
		connConfig, err = rest.InClusterConfig()
		if err != nil {
			log.Error(err)
			return err
		}
	} else if conf.ConnectType == config.ConfigFile {
		var err error
		connConfig, err = clientcmd.BuildConfigFromFlags("", conf.ConfigFile)
		if err != nil {
			log.Error(err)
			return err
		}
	}
	clientset, err := kubernetes.NewForConfig(connConfig)
	if err != nil {
		log.Error(err)
		return err
	}
	ms.Client = clientset
	return nil
}

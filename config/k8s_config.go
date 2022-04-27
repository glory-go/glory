package config

import (
	"fmt"
)

import (
	"github.com/glory-go/glory/tools"
)

type (
	ConnType    string
	K8sProvider string
)

const (
	InCluster  ConnType = "incluster"
	ConfigFile ConnType = "config"

	DefaultProvider K8sProvider = "default"
	AliyunProvider  K8sProvider = "aliyun"
)

type K8SConfig struct {
	ConfigSource    string            `yaml:"config_source"`
	Namespace       string            `yaml:"namespace"`
	DefaultUsername string            `yaml:"username"`
	ConnectType     ConnType          `yaml:"conn_type"`
	ConfigFile      string            `yaml:"config_path"`
	PVC             map[string]string `yaml:"pvc"`
	Provider        K8sProvider       `yaml:"provider"`
}

func (s *K8SConfig) checkAndFix() {
	if err := tools.ReadFromEnvIfNeed(s); err != nil {
		fmt.Println("warn: K8SConfig checkAndFix failed with err = ", err)
	}
}

package config

import (
	"fmt"

	"github.com/glory-go/glory/tools"
)

type ConfigCenterConfig struct {
	OrgName    string `yaml:"org_name"` // 可选 classroom|ide|children|goonline: goonline为公共服务，比如前端数据上报
	ServerName string `yaml:"server_name"`

	ConfigSource string `yaml:"config_source"`
	Name         string `yaml:"name"` // 配置中心配置名 目前只支持nacos

	// aliCloud config
	EndPoint     string `yaml:"end_point"`
	NamespaceID  string `yaml:"namespace_id"`
	AccessKeyID  string `yaml:"access_key_id"`
	AccessSecret string `yaml:"access_secret"`

	// env
	EnvNamespaceIDMap map[string]string `yaml:"env_map"` // get from environment
}

func (s *ConfigCenterConfig) checkAndFix() {
	if err := tools.ReadFromEnvIfNeed(s); err != nil {
		fmt.Println("warn: MetricsConfig checkAndFix failed with err = ", err)
	}

	if s.OrgName == "" {
		panic("please add your service org_name in config file from: classroom|ide|children|goonline")
	}

	if s.ServerName == "" {
		panic("please add your server_name in config file!")
	}
}

func (s *ConfigCenterConfig) GetAppKey() string {
	return s.OrgName + "_" + s.ServerName
}

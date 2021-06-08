package config

import (
	"fmt"

	"github.com/glory-go/glory/tools"
)

type OssConfig struct {
	ConfigSource  string `yaml:"config_source"`
	OssType       string `yaml:"oss_type"`
	Buckname      string `yaml:"buckname"`
	OssAccessKey  string `yaml:"oss_access_key"`
	OssSecretKey  string `yaml:"oss_secret_key"`
	OssDomainName string `yaml:"oss_domain_name"`
}

func (s *OssConfig) checkAndFix() {
	if err := tools.ReadFromEnvIfNeed(s); err != nil {
		fmt.Println("warn: OssConfig checkAndFix failed with err = ", err)
	}
}

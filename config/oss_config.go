package config

import (
	"fmt"

	"github.com/glory-go/glory/tools"
)

type OssConfig struct {
	ConfigSource string `yaml:"config_source"`
	OssType      string `yaml:"oss_type"` // qiniu\aliyun
	// 通用配置
	Region        string `yaml:"region"`
	Buckname      string `yaml:"buckname"`
	OssDomainName string `yaml:"oss_domain_name"`
	OssAccessKey  string `yaml:"oss_access_key"`
	OssSecretKey  string `yaml:"oss_secret_key"`
	// 阿里云
	RoleArn  string `yaml:"role_arn"`
	Endpoint string `yaml:"endpoint"`
}

func (s *OssConfig) checkAndFix() {
	if err := tools.ReadFromEnvIfNeed(s); err != nil {
		fmt.Println("warn: OssConfig checkAndFix failed with err = ", err)
	}
}

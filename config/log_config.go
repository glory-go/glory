package config

import (
	"fmt"

	"github.com/glory-go/glory/tools"
)

type LogConfig struct {
	ConfigSource string `yaml:"config_source"`
	LogType      string `yaml:"log_type"`
	FilePath     string `yaml:"file_path"`
	LogLevel     string `yaml:"level"`
	ElasticAddr  string `yaml:"elastic_addr"`

	// aliyun sls config
	ProjectName    string `yaml:"project_name"`
	AccessKeyID    string `yaml:"access_key_id"`
	AccessSecret   string `yaml:"access_secret"`
	EndPoint       string `yaml:"endpoint"`
	LogStoreName   string `yaml:"log_store_name"`
	UploadInterval int    `yaml:"interval"`
}

func (s *LogConfig) checkAndFix() {
	if err := tools.ReadFromEnvIfNeed(s); err != nil {
		fmt.Println("warn: LogConfig checkAndFix failed with err = ", err)
	}
	if s.UploadInterval <= 0 {
		s.UploadInterval = 5
	}
	if s.LogLevel == "" {
		s.LogLevel = "debug"
	}
	if s.LogType == "" {
		s.LogLevel = "console"
	}

	if GlobalServerConf.AliCloudCommonConfig == nil {
		return
	}
	if secret := GlobalServerConf.AliCloudCommonConfig.AccessSecret; s.AccessSecret == "" && secret != "" {
		s.AccessSecret = secret
	}
	if id := GlobalServerConf.AliCloudCommonConfig.AccessKeyID; s.AccessSecret == "" && id != "" {
		s.AccessKeyID = id
	}
}

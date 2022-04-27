package config

import (
	"fmt"
)

import (
	"github.com/glory-go/glory/tools"
)

// AliCloudCommonConfig store global alicloud common auth data, which makes user to defined once
type AliCloudCommonConfig struct {
	ConfigSource  string            `yaml:"config_source"`
	AccessKeyID   string            `yaml:"access_key_id"`
	AccessSecret  string            `yaml:"access_secret"`
	SMSTemplateID map[string]string `yaml:"sms_template_id"`
	SMSSignName   map[string]string `yaml:"sms_sign_name"`
}

func (s *AliCloudCommonConfig) checkAndFix() {
	if err := tools.ReadFromEnvIfNeed(s); err != nil {
		fmt.Println("warn: AliCloudCommonConfig checkAndFix failed with err = ", err)
	}
}

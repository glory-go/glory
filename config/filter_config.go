package config

import (
	"fmt"
)

import (
	"github.com/glory-go/glory/tools"
)

type FilterConfig struct {
	ConfigSource string `yaml:"config_source"`

	//FilterName now only support chain|jaeger
	FilterName string `yaml:"filter_name"`
	// chain filter config
	SubFiltersKey []string `yaml:"sub_filters"`

	// jaeger filter config
	JaegerType   string  `yaml:"jaeger_type"` // self (default 自建jaeger) or alicloud(阿里云jaeger服务)
	SamplerType  string  `yaml:"jaeger_config_type"`
	SamplerParam float64 `yaml:"jaeger_config_param"`
	Address      string  `yaml:"jaeger_address"`

	// jaeger aliyun sender config
	AliyunToken1 string `yaml:"aliyun_token_1"`
	AliyunToken2 string `yaml:"aliyun_token_2"`
}

func (s *FilterConfig) checkAndFix() {
	if err := tools.ReadFromEnvIfNeed(s); err != nil {
		fmt.Println("warn: JaegerConfig checkAndFix failed with err = ", err)
	}
	if s.SamplerType == "" {
		s.SamplerType = "const"
	}

	if s.SamplerParam == 0 {
		s.SamplerParam = 1
	}
}

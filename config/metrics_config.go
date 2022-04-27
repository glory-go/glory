package config

import (
	"fmt"
)

import (
	"github.com/glory-go/glory/tools"
)

type MetricsConfig struct {
	ConfigSource string `yaml:"config_source"`
	MetricsType  string `yaml:"metrics_type"`
	ActionType   string `yaml:"action_type"`
	ClientPort   string `yaml:"client_port"`
	ClientPath   string `yaml:"client_path"`
	GateWayHost  string `yaml:"gateway_host"`
	GateWayPort  string `yaml:"gateway_port"`
	JobName      string `yaml:"job_name"` // 数据上报job_name，默认为配置文件中的server_name
}

func (s *MetricsConfig) checkAndFix() {
	if err := tools.ReadFromEnvIfNeed(s); err != nil {
		fmt.Println("warn: MetricsConfig checkAndFix failed with err = ", err)
	}
}

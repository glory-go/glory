package config

import (
	"fmt"

	"github.com/glory-go/glory/tools"
)

type ServiceConfig struct {
	ConfigSource string   `yaml:"config_source"`
	Protocol     string   `yaml:"protocol"`
	RegistryKey  string   `yaml:"registry_key"`
	Port         int      `yaml:"port"`
	ServiceID    string   `yaml:"service_id"`
	FiltersKey   []string `yaml:"filters_key"`
}

func (s *ServiceConfig) checkAndFix() {
	if err := tools.ReadFromEnvIfNeed(s); err != nil {
		fmt.Println("warn: ServiceConfig checkAndFix failed with err = ", err)
	}
}

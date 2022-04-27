package config

import (
	"fmt"
)

import (
	"github.com/glory-go/glory/tools"
)

type RedisConfig struct {
	ConfigSource string `yaml:"config_source"`
	Host         string `yaml:"host"`
	Port         string `yaml:"port"`
	Password     string `yaml:"password"`
}

func (s *RedisConfig) checkAndFix() {
	if err := tools.ReadFromEnvIfNeed(s); err != nil {
		fmt.Println("warn: RedisConfig checkAndFix failed with err = ", err)
	}
}

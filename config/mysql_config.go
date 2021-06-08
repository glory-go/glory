package config

import (
	"fmt"

	"github.com/glory-go/glory/tools"
)

type MysqlConfig struct {
	ConfigSource string `yaml:"config_source"`
	Host         string `yaml:"host"`
	Port         string `yaml:"port"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	DBName       string `yaml:"dbname"`
}

func (s *MysqlConfig) checkAndFix() {
	if err := tools.ReadFromEnvIfNeed(s); err != nil {
		fmt.Println("warn: MysqlConfig checkAndFix failed with err = ", err)
	}
}

package config

import (
	"fmt"

	"github.com/glory-go/glory/tools"
)

type MongoDBConfig struct {
	ConfigSource   string `yaml:"config_source"`
	Host           string `yaml:"host"`
	Port           string `yaml:"port"`
	Username       string `yaml:"username"`
	Password       string `yaml:"password"`
	DBName         string `yaml:"dbname"`
	CollectionName string `yaml:"collectionName"`
}

func (s *MongoDBConfig) checkAndFix() {
	if err := tools.ReadFromEnvIfNeed(s); err != nil {
		fmt.Println("warn: MongoDBConfig checkAndFix failed with err = ", err)
	}
}

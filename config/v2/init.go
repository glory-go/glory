package config

import (
	"os"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

var once sync.Once

// value of env keys can be changed from environment

// EnvKeyGloryConfigPath is absolute/relate path to glory.yaml
const EnvKeyGloryConfigPath = "GLORY_CONFIG_PATH" // default val is "config/glory.yaml"

// EnvKeyGloryEnv if is set to dev,then:
// 1. choose config center with namespace dev
// 2. choose config path like "config/glory_dev.yaml
const EnvKeyGloryEnv = "GLORY_ENV" //

const DefaultConfigPath = "config/glory.yaml"

func GetConfigPath() string {
	configPath := ""
	env := os.Getenv(EnvKeyGloryEnv)

	configFilePath := DefaultConfigPath
	if os.Getenv(EnvKeyGloryConfigPath) != "" {
		configFilePath = os.Getenv(EnvKeyGloryConfigPath)
	}
	prefix := strings.Split(configFilePath, ".yaml")
	// prefix == ["config/glory", ""]
	if len(prefix) != 2 {
		panic("Invalid config file path = " + configFilePath)
	}
	// get target env yaml file
	if env != "" {
		configPath = prefix[0] + "_" + env + ".yaml"
	} else {
		configPath = configFilePath
	}
	return configPath
}

func loadFileConfig() {
	configPath := GetConfigPath()

	ConfigInstance = viper.New()
	ConfigInstance.AddConfigPath(configPath)
}

func Init() {
	once.Do(func() {
		loadFileConfig()
	})
}

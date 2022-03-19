package configmanager

import (
	"log"
	"os"
	"strings"
	"sync"

	"github.com/glory-go/glory/common"
	"github.com/spf13/viper"
)

var (
	once   sync.Once
	inited bool
)

func Init() {
	once.Do(func() {
		loadConfigCenterConfig()
		inited = true
	})
}

func GetConfigCenterPath() string {
	configPath := ""
	env := os.Getenv(common.EnvKeyGloryEnv)

	configFilePath := DefaultConfigCenterConfigPath
	if os.Getenv(EnvKeyGloryConfigCenterPath) != "" {
		configFilePath = os.Getenv(EnvKeyGloryConfigCenterPath)
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

func loadConfigCenterConfig() {
	path := GetConfigCenterPath()

	viperInstance := viper.GetViper()
	viperInstance.SetConfigFile(path)
	if err := viperInstance.ReadInConfig(); err != nil {
		log.Panic("load config center meets error", err)
	}
	configCenterViperInstance = viperInstance
}

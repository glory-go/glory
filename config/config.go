package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

// value of env keys can be changed from environment

// EnvKeyGloryConfigPath is absolute/relate path to glory.yaml
const EnvKeyGloryConfigPath = "GLORY_CONFIG_PATH" // default val is "config/glory.yaml"

// EnvKeyGloryEnv if is set to dev,then:
// 1. choose config center with namespace dev
// 2. choose config path like "config/glory_dev.yaml
const EnvKeyGloryEnv = "GLORY_ENV" //

const DefaultConfigPath = "config/glory.yaml"

var GlobalServerConf *ServerConfig

type config interface {
	checkAndFix()
}

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

func loadFileConfig() error {
	conf := NewServerConfig()
	configPath := GetConfigPath()

	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Printf("error: yamlFile get error= %v\n", err)
		return errors.New("yamlFile.Get err ")
	}

	err = yaml.Unmarshal(yamlFile, conf)
	if err != nil {
		fmt.Printf("yamlFile Unmarshal err: %v\n", err)
		return err
	}
	// make sure when checkAndFix procedure, config can access the whole config and do refactor
	GlobalServerConf = conf
	conf.checkAndFix()
	// store the final config
	GlobalServerConf = conf
	return nil
}

func init() {
	if err := loadFileConfig(); err != nil {
		fmt.Printf("load conf from file failed with err =  %s, try to load from default config\n", err)
		loadDefaultConfig()
	}
	// no config file
}

func loadDefaultConfig() {
	defaultSrvConf := NewServerConfig()
	defaultSrvConf.LogConfigs["default-log"] = &LogConfig{
		LogType:  "console",
		LogLevel: "debug",
	}
	defaultSrvConf.OrgName = "default_org"
	defaultSrvConf.ServerName = "default_server"
	GlobalServerConf = defaultSrvConf
}

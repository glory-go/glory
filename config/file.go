package config

import (
	"os"
	"strings"
)

import (
	perrors "github.com/pkg/errors"

	"gopkg.in/yaml.v3"
)

// value of env keys can be changed from environment
// EnvKeyGloryConfigPath is absolute/relate path to glory.yaml
const EnvKeyGloryConfigPath = "GLORY_CONFIG_PATH" // default val is "config/glory.yaml"

// EnvKeyGloryEnv if is set to dev,then:
// 1. choose config center with namespace dev
// 2. choose config path like "config/glory_dev.yaml
const EnvKeyGloryEnv = "GLORY_ENV" //

// EnvKeyGloryConfigCenterConfigPath is absolute/relate path to glory_config_center.yaml
const EnvKeyGloryConfigCenterConfigPath = "GLORY_CONFIG_CENTER_CONFIG_PATH"

const DefaultConfigPath = "config/glory.yaml"
const DefaultConfigCenterConfigPath = "config/config_center.yaml"

func GetGloryEnv() string {
	return os.Getenv(EnvKeyGloryEnv)
}

func GetConfigPath() string {
	configPath := ""
	env := GetGloryEnv()

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

func loadProperty(splitedConfigName []string, index int, tempConfigMap Config, configStructPtr interface{}) error {
	subConfig, ok := tempConfigMap[splitedConfigName[index]]
	if !ok {
		return perrors.Errorf("property %s's key %s not found", splitedConfigName, splitedConfigName[index])
	}
	if index+1 == len(splitedConfigName) {
		targetConfigByte, err := yaml.Marshal(subConfig)
		if err != nil {
			return perrors.Errorf("property %s's key %s invalid, error = %s", splitedConfigName, splitedConfigName[index], err)
		}
		err = yaml.Unmarshal(targetConfigByte, configStructPtr)
		if err != nil {
			return perrors.Errorf("property %s's key %s doesn't match type %+v, error = %s", splitedConfigName, splitedConfigName[index], configStructPtr, err)
		}
		return nil
	}
	subMap, ok := subConfig.(Config)
	if !ok {
		return perrors.Errorf("property %s's key %s of config is not map[string]string, which is %+v", splitedConfigName,
			splitedConfigName[index], subConfig)
	}
	return loadProperty(splitedConfigName, index+1, subMap, configStructPtr)
}

package config

import (
	"io/ioutil"
	"strings"
)

import (
	"github.com/fatih/color"

	"gopkg.in/yaml.v3"
)

type Config map[string]interface{}

var config Config

func SetConfig(yamlBytes []byte) error {
	return yaml.Unmarshal(yamlBytes, &config)
}

func Load() error {
	configPath := GetConfigPath()

	color.Blue("[Config] Load config file from %s", configPath)
	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		color.Red("Load glory config file failed. %v\n The load procedure is continue\n", err)
		return nil
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		color.Red("yamlFile Unmarshal err: %v\n", err)
		return err
	}
	parseConfigSource(config)
	return nil
}

// LoadConfigByPrefix prefix is a.b.c, configStructPtr is interface ptr
func LoadConfigByPrefix(prefix string, configStructPtr interface{}) error {
	if configStructPtr == nil {
		return nil
	}
	configProperties := strings.Split(prefix, ".")
	return loadProperty(configProperties, 0, config, configStructPtr)
}

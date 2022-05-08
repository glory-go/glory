package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

import (
	"gopkg.in/yaml.v3"
)

type Config map[string]interface{}

var config Config

func SetConfig(yamlBytes []byte) error {
	return yaml.Unmarshal(yamlBytes, &config)
}

func Load() error {
	configPath := GetConfigPath()

	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Printf("error: yamlFile get error= %v\n", err)
		return errors.New("yamlFile.Get err ")
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		fmt.Printf("yamlFile Unmarshal err: %v\n", err)
		return err
	}
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

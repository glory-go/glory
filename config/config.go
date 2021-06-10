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

// EnvKeyGloryConfigCenterConfigPath is absolute/relate path to glory_config_center.yaml
const EnvKeyGloryConfigCenterConfigPath = "GLORY_CONFIG_CENTER_CONFIG_PATH"

const DefaultConfigPath = "config/glory.yaml"
const DefaultConfigCenterConfigPath = "config/config_center.yaml"

var GlobalServerConf *ServerConfig

type config interface {
	checkAndFix()
}

func loadFileConfig() error {
	conf := NewServerConfig()
	env := os.Getenv(EnvKeyGloryEnv)

	configFilePath := DefaultConfigPath
	if os.Getenv(EnvKeyGloryConfigPath) != "" {
		configFilePath = os.Getenv(EnvKeyGloryConfigPath)
	}

	var yamlFile []byte
	var err error
	prefix := strings.Split(configFilePath, ".yaml")
	// prefix == ["config/glory", ""]
	if len(prefix) != 2 {
		panic("Invalid config file path = " + configFilePath)
	}
	// get target env yaml file
	if env != "" {
		yamlFile, err = ioutil.ReadFile(prefix[0] + "_" + env + ".yaml")
	} else {
		yamlFile, err = ioutil.ReadFile(configFilePath)
	}
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

// oadConfigCenterConfig load config center's config from file, default is config/config_center_config.yaml
func loadConfigCenterConfig() bool {
	configCenterConfigFilePath := DefaultConfigCenterConfigPath
	if path := os.Getenv(EnvKeyGloryConfigCenterConfigPath); path != "" {
		configCenterConfigFilePath = path
	}
	yamlFile, err := ioutil.ReadFile(configCenterConfigFilePath)
	if err != nil {
		fmt.Println("config center info: can't load config center config at " + DefaultConfigCenterConfigPath)
		return false
	}
	configCenterConfig := ConfigCenterConfig{}
	err = yaml.Unmarshal(yamlFile, &configCenterConfig)
	if err != nil {
		fmt.Printf("error: Unmarshal config center config err: %v\n", err)
		return false
	}
	configCenterConfig.checkAndFix()
	if configCenterConfig.Name != "nacos" && configCenterConfig.Name != "" {
		fmt.Println("error: config center name ", configCenterConfig.Name, "not support")
		return false
	}

	// new config center
	env := os.Getenv(EnvKeyGloryEnv)
	if env == "" {
		panic("please config GLORY_ENV environment variable")
	}
	nacosConfigCenter := newNacosConfigCenter(env)
	if err := nacosConfigCenter.Conn(&configCenterConfig); err != nil {
		fmt.Println("error: config center", configCenterConfig.Name, "Conn err = ", err)
		return false
	}
	if cfg, err := nacosConfigCenter.LoadConfig(); err != nil {
		fmt.Println("error: config center ", configCenterConfig.Name, "LoadConfig err = ", err)
		return false
	} else {
		GlobalServerConf = cfg
		fmt.Println("error: load config from remote config center successful")
		return true
	}

}

func init() {
	if ok := loadConfigCenterConfig(); ok {
		return
	}
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

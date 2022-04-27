package config

import (
	"log"
)

import (
	"github.com/spf13/viper"
)

var viperInstance *viper.Viper

func init() {
	viperInstance = viper.GetViper()
	viperInstance.SetConfigFile(GetConfigPath())
	if err := viper.ReadInConfig(); err != nil {
		log.Println("Fatal: load config meets error", err)
	}
}

func GetViperConfig() *viper.Viper {
	return viperInstance
}

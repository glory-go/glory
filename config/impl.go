package config

import (
	"sync"

	"github.com/spf13/viper"
)

var (
	configData      *viper.Viper
	
	getConfigOnce   sync.Once
	parseConfigOnce sync.Once
)

func GetConfig(path string) *viper.Viper {
	getConfigOnce.Do(func() {})

	return configData
}

func ParseConfig(config *viper.Viper) {
	parseConfigOnce.Do(func() {})
}

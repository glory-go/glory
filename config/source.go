package config

import (
	"os"
)

import (
	"github.com/fatih/color"
)

const ConfigSourceKey = "_glory_config_source"
const ConfigSourceEnvFlag = "env"

func parseConfigSource(config Config) {
	envFlag := false
	if source, ok := config[ConfigSourceKey]; ok {
		if sourceStr, okStr := source.(string); okStr && sourceStr == ConfigSourceEnvFlag {
			color.Blue("[Config] %s under %v is set to %s, try to load from env", ConfigSourceKey, config, ConfigSourceEnvFlag)
			envFlag = true
		}
	}
	for k, v := range config {
		if val, ok := v.(string); ok {
			if envFlag {
				if envVal := os.Getenv(val); envVal != "" {
					config[k] = envVal
				} else {
					color.Blue("[Config] Try to load %s from env failed", val)
				}
			}
		} else if subConfig, ok := v.(Config); ok {
			parseConfigSource(subConfig)
		}
	}
}

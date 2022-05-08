package config

import (
	"os"
)

const ConfigSourceKey = "_glory_config_source"
const ConfigSourceEnvFlag = "env"

func parseConfigSource(config Config) {
	envFlag := false
	if source, ok := config[ConfigSourceKey]; ok {
		if sourceStr, okStr := source.(string); okStr && sourceStr == ConfigSourceEnvFlag {
			envFlag = true
		}
	}
	for k, v := range config {
		if val, ok := v.(string); ok {
			if envFlag {
				if envVal := os.Getenv(val); envVal != "" {
					config[k] = envVal
				}
			}
		} else if subConfig, ok := v.(Config); ok {
			parseConfigSource(subConfig)
		}
	}
}

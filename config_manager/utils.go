package configmanager

import "os"

func ReadFromEnvIfNeed(conf map[string]string) {
	if conf["config_source"] != "env" {
		return
	}
	for k, v := range conf {
		if data := os.Getenv(v); data != "" {
			conf[k] = data
		}
	}
}

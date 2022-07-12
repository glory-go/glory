package config

func Init() {
	configPath := GetConfigPath()
	config := GetConfig(configPath)
	
}
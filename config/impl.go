package config

import (
	"github.com/spf13/viper"
)

var (
	configData *viper.Viper
)

// loadRawConfig 从文件中加载原始配置信息
func loadRawConfig(path string) {
}

// convertConfigFromEnv 将配置中需要从环境变量获取的内容进行读取和替换
func convertConfigFromEnv() {

}

// convertConfigFromConfigCenter 从配置中心中读取并更新配置
func convertConfigFromConfigCenter() {}

func getConfig(key string, raw interface{}) {
	err := configData.UnmarshalKey(key, raw)
	if err != nil {
		panic(err)
	}
}

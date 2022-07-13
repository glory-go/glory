package config

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/viper"
)

var (
	configData *viper.Viper
)

// loadRawConfig 从文件中加载原始配置信息
func loadRawConfig(reader io.Reader) {
	configData = viper.GetViper()
	configData.SetConfigType("YAML")
	if err := configData.ReadConfig(reader); err != nil {
		panic(err)
	}
}

// convertConfigFromEnv 将配置中需要从环境变量获取的内容进行读取和替换
func convertConfigFromEnv() {
	keys := configData.AllKeys()
	for _, key := range keys {
		match := placeHolderRegexp.FindStringSubmatch(configData.GetString(key))
		if len(match) < 3 {
			configData.Set(key, configData.Get(key)) // 这行不加的话，貌似没有set过的值都会被忽略掉
			continue
		}
		// 获取配置中心实例
		center, err := GetConfigCenter(match[1])
		if err != nil {
			panic(err)
		}
		// 读取参数
		params := make([]any, 0)
		if err := jsoniter.UnmarshalFromString(match[2], &params); err != nil {
			panic(err)
		}
		// 从配置中心读取配置
		data, err := center.Get(params...)
		if err != nil {
			panic(err)
		}
		configData.Set(key, data)
	}
}

// convertConfigFromConfigCenter 从配置中心中读取并更新配置
func convertConfigFromConfigCenter() {}

func getConfig(key string, raw any) {
	err := configData.UnmarshalKey(key, raw)
	if err != nil {
		panic(err)
	}
}

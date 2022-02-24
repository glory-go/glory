package config

import (
	"strings"

	configmanager "github.com/glory-go/glory/config_manager"
)

const (
	ConfigSourceKey  = "config_source"
	GroupKeySplitter = "$$"
)

func ReplaceStringValueFromConfigCenter(conf map[string]interface{}) error {
	configSourceInter, ok := conf[ConfigSourceKey]
	if !ok || configSourceInter == nil {
		return nil
	}
	configSource, ok := configSourceInter.(string)
	if !ok || configSource == "" {
		return nil
	}
	tmp, err := replaceMapValueFromConfigCenter(configSource, conf)
	if err != nil {
		return err
	}
	conf = tmp
	return nil
}

func replaceMapValueFromConfigCenter(configSource string, conf map[string]interface{}) (map[string]interface{}, error) {
	tmpVal := make(map[string]interface{})
	for k, v := range conf {
		var (
			err error
			val interface{}
		)
		switch v.(type) {
		case map[string]interface{}:
			val, err = replaceMapValueFromConfigCenter(configSource, v.(map[string]interface{}))
		case string:
			val, err = replaceStringValueFromConfigCenter(configSource, v.(string))
		default:
			continue
		}
		if err != nil {
			return nil, err
		}
		tmpVal[k] = val
	}
	return tmpVal, nil
}

func replaceStringValueFromConfigCenter(configSource, rawVal string) (string, error) {
	center, err := configmanager.GetConfigCenter(configSource)
	if err != nil {
		return "", err
	}
	group, key := getGroupAndKey(rawVal)
	return center.LoadConfig(key, group)
}

func getGroupAndKey(val string) (group, key string) {
	result := strings.SplitN(val, GroupKeySplitter, 2)
	if len(result) <= 1 {
		return val, ""
	}
	return result[0], result[1]
}

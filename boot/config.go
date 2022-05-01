package boot

import (
	"strings"
)

import (
	perrors "github.com/pkg/errors"
)

var userConfig map[interface{}]interface{}

func implConfig(configName string) (string, error) {
	configProperties := strings.Split(configName, ".")
	return loadProperty(configProperties, 0, userConfig)
}

func loadProperty(splitedConfigName []string, index int, tempConfigMap map[interface{}]interface{}) (string, error) {
	subConfig, ok := tempConfigMap[splitedConfigName[index]]
	if !ok {
		return "", perrors.Errorf("property %s's key %s not found", splitedConfigName, splitedConfigName[index])
	}
	if index+1 == len(splitedConfigName) {
		target, ok := subConfig.(string)
		if !ok {
			return "", perrors.Errorf("property %s's key %s of config is not string, which is %+v", splitedConfigName,
				splitedConfigName[index], subConfig)
		}
		return target, nil
	}
	subMap, ok := subConfig.(map[interface{}]interface{})
	if !ok {
		return "", perrors.Errorf("property %s's key %s of config is not map[string]string, which is %+v", splitedConfigName,
			splitedConfigName[index], subConfig)
	}
	return loadProperty(splitedConfigName, index+1, subMap)
}

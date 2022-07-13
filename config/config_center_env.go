package config

import (
	"errors"
	"fmt"
	"os"
)

const (
	EnvConfigCenterName = "env"
)

var (
	EnvConfigCenter envConfigCenter
)

type envConfigCenter struct{}

func GetEnvConfigCenter() ConfigCenter {
	return EnvConfigCenter
}

func (envConfigCenter) Name() string { return EnvConfigCenterName }

func (envConfigCenter) Init(config map[string]any) error { return nil }

func (envConfigCenter) Get(params ...any) (string, error) {
	if len(params) == 0 {
		return "", errors.New("envConfigCenter: no enough params")
	}
	key, ok := params[0].(string)
	if !ok {
		return "", fmt.Errorf("envConfigCenter: param$0 [%v] not a string", params[0])
	}
	return os.Getenv(key), nil
}

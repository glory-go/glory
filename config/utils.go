package config

import (
	"fmt"
	"os"
	"strings"
)

func ChangeDefaultConfigPath(path string) {
	// 检查是否以.yaml或yml结尾
	if !(strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml")) {
		panic(fmt.Sprintf("path %s not has yaml suffix", path))
	}
	defaultConfigPath = path
}

func GetConfigPath() string {
	path := defaultConfigPath
	// 获取env信息
	env := strings.Trim(os.Getenv(GLORY_ENV), " ")
	if env == "" {
		return path
	}
	// 注入环境信息
	ptrIdx := strings.LastIndex(path, ".")
	path = path[:ptrIdx] + "." + env + path[ptrIdx:]

	return path
}


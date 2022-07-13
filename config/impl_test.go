package config

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_loadRawConfig(t *testing.T) {
	// 构造初始配置文件内容
	data := strings.NewReader(`
a: test
`)
	// 读取配置并测试
	loadRawConfig(data)
	assert.NotNil(t, configData)
	assert.Equal(t, "test", configData.GetString("a"))
}

func Test_convertConfigFromEnv(t *testing.T) {
	// 构造初始配置文件内容
	data := strings.NewReader(`
a: 1
b: $env{["b_env"]}
c:
  c1: $env{["c1_env"]}
`)
	os.Setenv("b_env", "b_test")
	os.Setenv("c1_env", "c1_test")
	// 读取配置
	loadRawConfig(data)
	// 测试环境变量
	convertConfigFromEnv()
	assert.Equal(t, "1", configData.GetString("a"))
	assert.Equal(t, "b_test", configData.GetString("b"))
	assert.Equal(t, "c1_test", configData.GetStringMapString("c")["c1"])
}

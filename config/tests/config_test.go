package tests

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/glory-go/glory/v2/config"
	mockconfig "github.com/glory-go/glory/v2/config/mock"
	"github.com/golang/mock/gomock"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func Test_Config(t *testing.T) {
	setUp()
	ctrl := gomock.NewController(t)

	/** 注册配置中心 **/
	configCenter := mockconfig.NewMockConfigCenter(ctrl)
	configCenter.EXPECT().Name().Return("mock_config_center").MinTimes(1)
	configCenter.EXPECT().Init(gomock.Any()).DoAndReturn(func(config map[string]any) error {
		// 与替换环境变量值后配置中心的配置应保持一致
		assert.Equal(t, "test", config["a1"])
		assert.Equal(t, "a2_test", config["a2"])
		target := make([]int, 0)
		mapstructure.Decode(config["a3"], &target)
		assert.Equal(t, []int{1, 2, 3}, target)

		return nil
	}).Times(1)
	configCenter.EXPECT().Get(gomock.Len(3)). // 配置文件中定义的参数长度
							DoAndReturn(func(params ...any) (string, error) {
			// 原封不动返回第二个参数的字符串形式
			assert.Len(t, params, 3)
			assert.Equal(t, "b3_param1", params[0])
			assert.Equal(t, float64(1), params[1])
			assert.Equal(t, true, params[2])

			return fmt.Sprint(params[1]), nil
		}).Times(1) // 配置文件中用到这个配置中心的次数
	config.RegisterConfigCenter(configCenter)

	/** 注册组件 **/
	component := mockconfig.NewMockComponent(ctrl)
	component.EXPECT().Name().Return("mock_component").MinTimes(1)
	component.EXPECT().Init(gomock.Any()).DoAndReturn(func(config map[string]any) error {
		assert.Len(t, config, 3)
		assert.Equal(t, "test_b1", config["b1"])
		assert.Equal(t, "b2_test", config["b2"])
		assert.Equal(t, "1", config["b3"])

		return nil
	}).Times(1)
	config.RegisterComponent(component)

	/** 初始化 **/
	config.Init()
}

func setUp() {
	path := fmt.Sprintf("/tmp/test_config_%d.yaml", time.Now().Unix())
	config.ChangeDefaultConfigPath(path)
	content := `
config_center:
  mock_config_center:
    a1: test
    a2: $env{["a2_env"]}
    a3: [1, 2, 3]
mock_component:
  b1: test_b1
  b2: $env{["b2_env"]}
  b3: $mock_config_center{["b3_param1", 1, true]}
`
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	_, err = file.WriteString(content)
	if err != nil {
		panic(err)
	}
	os.Setenv("a2_env", "a2_test")
	os.Setenv("b2_env", "b2_test")
}

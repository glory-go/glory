package config

import (
	"fmt"
	"os"
	"sync"
)

var (
	inited   bool
	initOnce sync.Once
)

func init() {
	registerInnerConfigCenter(GetEnvConfigCenter())
}

func Init() {
	initOnce.Do(func() {
		defer func() {
			inited = true
		}()
		/* 初始化原始文件配置 */
		// 获取配置文件的地址
		configPath := GetConfigPath()
		// 加载文件中最原始的配置内容
		file, err := os.Open(configPath)
		if err != nil {
			panic(err)
		}
		loadRawConfig(file)
		// 替换环境变量中的内容
		convertConfigFromEnv()
		/* 初始化原始文件配置结束 */

		/* 初始化配置中心 */
		// 读取配置中心的配置
		rawConfig := make(map[string]map[string]any)
		getConfig(CONFIG_CENTER_KEY, &rawConfig)
		// 读取注册的配置中心，并初始化
		iterConfigRegistry(func(name string, center ConfigCenter) error {
			if skipInitConfigCenterName.Contains(name) {
				return nil
			}
			config, ok := rawConfig[name]
			if !ok {
				return fmt.Errorf("config of config center %s not found", name)
			}
			if err := center.Init(config); err != nil {
				return err
			}
			return nil
		})
		/* 初始化配置中心结束 */

		// 基于配置中心更新配置文件内容
		convertConfigFromConfigCenter()

		/* 初始化用户注册的组件 */
		iterComponentRegistry(func(name string, component Component) error {
			rawConfig := make(map[string]any)
			getConfig(name, &rawConfig)
			if err := component.Init(rawConfig); err != nil {
				return err
			}
			return nil
		})
		/* 初始化组件结束 */
	})
}

package logrus

import (
	"fmt"
	"sync"

	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

const (
	LogrusComponentName = "logrus"
)

type logrusComponent struct{}

var (
	component *logrusComponent
	once      sync.Once
)

func getLogrusComponent() *logrusComponent {
	once.Do(func() {
		component = &logrusComponent{}
	})
	return component
}

func (c *logrusComponent) Name() string { return LogrusComponentName }

func (c *logrusComponent) Init(conf map[string]any) error {
	rawLogrusConf, ok := conf[LogrusConfigKey]
	logrusConf := &logrusComponentConfig{}
	if ok {
		if err := mapstructure.Decode(rawLogrusConf, logrusConf); err != nil {
			return err
		}
	}
	// 初始化logrus
	// 逐一初始化用户定义的hooks
	rawHooksConf, ok := conf[HooksConfigKey]
	hooksConf := map[string]any{}
	if ok {
		if err := mapstructure.Decode(rawHooksConf, &hooksConf); err != nil {
			return err
		}
	}
	logrus.SetLevel(logrusConf.Level)
	for name, raw := range hooksConf {
		mapConf := make(map[string]any)
		if err := mapstructure.Decode(raw, &mapConf); err != nil {
			return err
		}
		if mapConf[HookTypeKey] == nil {
			return fmt.Errorf("hook with name [%s] not define type in config", name)
		}
		tp, _ := mapConf[HookTypeKey].(string)
		if tp == "" {
			return fmt.Errorf("hook with name [%s] not define type in config", name)
		}
		// 判断类型是否存在
		builderAny, ok := hooksRegistry.Load(tp)
		if !ok {
			return fmt.Errorf("hook builder with type [%s] not registered before", tp)
		}
		// 构建并添加日志
		builder := builderAny.(HookBuilder)
		hook, err := builder(mapConf)
		if err != nil {
			return err
		}
		logrus.AddHook(hook)
	}

	return nil
}

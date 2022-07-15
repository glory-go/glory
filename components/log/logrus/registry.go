package logrus

import (
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	hooksRegistry sync.Map
)

// HookBuilder 实现了hook的一个构建方法，glory将会根据用户配置调用已经注册的builder，从而完成hook的追加
type HookBuilder func(conf map[string]any) (logrus.Hook, error)

func RegisterHookBuilder(tp string, hook HookBuilder) {
	_, ok := hooksRegistry.LoadOrStore(tp, hook)
	if ok {
		panic(fmt.Sprintf("hook with type [%s] already registered before", tp))
	}
}

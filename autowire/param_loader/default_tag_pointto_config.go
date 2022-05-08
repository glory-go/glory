package param_loader

import (
	"fmt"
	"strings"
)

import (
	"github.com/pkg/errors"
)

import (
	"github.com/glory-go/glory/autowire"
	"github.com/glory-go/glory/autowire/util"
	"github.com/glory-go/glory/config"
)

type defaultTagPointToConfig struct {
}

func getDefaultTagPointToConfigPrefix(sd *autowire.StructDescriber, instanceName string) string {
	return fmt.Sprintf("autowire.%s.%s.%s.%s.param", sd.AutowireType(), util.GetStructName(sd.Interface), util.GetStructName(sd.Factory()), instanceName)
}

var defaultTagPointToConfigSingleton autowire.ParamLoader

func GetDefaultTagPointToConfigParamLoader() autowire.ParamLoader {
	if defaultTagPointToConfigSingleton == nil {
		defaultTagPointToConfigSingleton = &defaultTagPointToConfig{}
	}
	return defaultTagPointToConfigSingleton
}

/*
Load support load struct described like:
```go
normal.RegisterStructDescriber(&autowire.StructDescriber{
		Interface: new(Redis),
		Factory:   func() interface{}{
			return &Impl{}
		},
		ParamFactory: func() interface{}{
			return &Config{}
		},
		ConstructFunc: func(i interface{}, p interface{}) (interface{}, error) {
			return i, nil
		},
	})
}

type Config struct {
	Address  string
	Password string
	DB       string
}
```
with
Autowire type 'normal'
InterfaceName 'Redis'
StructName 'Impl'

from:

```yaml
extension:
  normal:
    Redis:
      Impl:
        param:
          address: 127.0.0.1
          password: xxx
          db: 0
```
*/
func (p *defaultTagPointToConfig) Load(sd *autowire.StructDescriber, fi *autowire.FieldInfo) (interface{}, error) {
	if fi == nil || sd == nil || sd.ParamFactory == nil {
		return nil, errors.New("not supported")
	}

	param := sd.ParamFactory()

	splitedTagValue := strings.Split(fi.TagValue, ",")
	if len(splitedTagValue) < 2 {
		return nil, errors.New("tag value not supported")
	}
	prefix := getDefaultTagPointToConfigPrefix(sd, splitedTagValue[1])
	if err := config.LoadConfigByPrefix(prefix, param); err != nil {
		return nil, err
	}
	return param, nil
}

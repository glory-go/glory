package param_loader

import (
	"fmt"
)

import (
	"github.com/pkg/errors"
)

import (
	"github.com/glory-go/glory/autowire"
	"github.com/glory-go/glory/autowire/util"
	"github.com/glory-go/glory/config"
)

type defaultConfig struct {
}

func getDefaultConfigPrefix(sd *autowire.StructDescriber) string {
	return fmt.Sprintf("autowire.%s.%s.%s.param", sd.AutowireType(), util.GetStructName(sd.Interface), util.GetStructName(sd.Factory()))
}

var defaultConfigParamLoaderSingleton autowire.ParamLoader

func GetDefaultConfigParamLoader() autowire.ParamLoader {
	if defaultConfigParamLoaderSingleton == nil {
		defaultConfigParamLoaderSingleton = &defaultConfig{}
	}
	return defaultConfigParamLoaderSingleton
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
autowire:
  normal:
    Redis:
      Impl:
        param:
          address: 127.0.0.1
          password: xxx
          db: 0
```
*/
func (p *defaultConfig) Load(sd *autowire.StructDescriber, fi *autowire.FieldInfo) (interface{}, error) {
	if sd == nil || sd.ParamFactory == nil {
		return nil, errors.New("not supporterd")
	}
	param := sd.ParamFactory()
	prefix := getDefaultConfigPrefix(sd)
	if err := config.LoadConfigByPrefix(prefix, param); err != nil {
		return nil, err
	}
	return param, nil
}

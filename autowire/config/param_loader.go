package config

import (
	"errors"
	"log"
	"strings"
)

import (
	"github.com/glory-go/glory/autowire"
	"github.com/glory-go/glory/config"
)

type paramLoader struct {
}

/*
Load support load config field like:
```go
Address  configString.ConfigString `config:"ConfigString,myConfig.myConfigSubPath.myConfigKey"`
```go

from:

```yaml
myConfig:
  myConfigSubPath:
      myConfigKey: myConfigValue
```
*/
func (p *paramLoader) Load(sd *autowire.StructDescriber, fi *autowire.FieldInfo) (interface{}, error) {
	if sd == nil || fi == nil || sd.ParamFactory == nil {
		return nil, errors.New("not supported")
	}
	splitedTagValue := strings.Split(fi.TagValue, ",")
	configPath := splitedTagValue[1]
	param := sd.ParamFactory()
	if err := config.LoadConfigByPrefix(configPath, param); err != nil {
		log.Println("load config path "+configPath+" error = ", err)
		// FIXME ignore config read error?
	}
	return param, nil
}

type Config struct {
	Address string
}

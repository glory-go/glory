package param_loader

import (
	"encoding/json"
	"log"
	"strings"
)

import (
	"github.com/pkg/errors"
)

import (
	"github.com/glory-go/glory/autowire"
)

type defaultTag struct {
}

var defaultTagParamLoaderSingleton autowire.ParamLoader

func GetDefaultTagParamLoader() autowire.ParamLoader {
	if defaultTagParamLoaderSingleton == nil {
		defaultTagParamLoaderSingleton = &defaultTag{}
	}
	return defaultTagParamLoaderSingleton
}

/*
Load support load param like:
```go
type Config struct {
	Address  string
	Password string
	DB       string
}
```

from field:

```go
NormalRedis  normalRedis.Redis  `normal:"Impl,address=127.0.0.1&password=xxx&db=0"`
```
*/
func (p *defaultTag) Load(sd *autowire.StructDescriber, fi *autowire.FieldInfo) (interface{}, error) {
	if sd == nil || fi == nil || sd.ParamFactory == nil {
		return nil, errors.New("not supported")
	}
	splitedTagValue := strings.Split(fi.TagValue, ",")
	if len(splitedTagValue) < 2 {
		return nil, errors.New("not supported")
	}
	kvs := strings.Split(splitedTagValue[1], "&")
	kvMaps := make(map[string]interface{})
	for _, kv := range kvs {
		splitedKV := strings.Split(kv, "=")
		if len(splitedKV) != 2 {
			return nil, errors.New("not supported")
		}
		kvMaps[splitedKV[0]] = splitedKV[1]
	}
	data, err := json.Marshal(kvMaps)
	if err != nil {
		log.Printf("error json marshal %s\n", err)
		return nil, err
	}
	param := sd.ParamFactory()
	if err := json.Unmarshal(data, param); err != nil {
		log.Printf("error jsonun marshal %s\n", err)
		return nil, err
	}
	return param, nil
}

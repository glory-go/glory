package param_loader

import (
	"github.com/glory-go/glory/autowire"
)

type defaultParamLoader struct {
	defaultConfigParamLoader     autowire.ParamLoader
	defaultTagParamLoader        autowire.ParamLoader
	defaultTagPointToParamLoader autowire.ParamLoader
}

var defaultParamLoaderSingleton autowire.ParamLoader

func GetDefaultParamLoader() autowire.ParamLoader {
	if defaultParamLoaderSingleton == nil {
		defaultParamLoaderSingleton = &defaultParamLoader{
			defaultConfigParamLoader:     GetDefaultConfigParamLoader(),
			defaultTagParamLoader:        GetDefaultTagParamLoader(),
			defaultTagPointToParamLoader: GetDefaultTagPointToConfigParamLoader(),
		}
	}
	return defaultParamLoaderSingleton
}

/*
Load try to load config from 3 types: ordered from harsh to loose

1. Try to use defaultTagPointToParamLoader to load from field tag
2. Try to use defaultTagParamLoader to load from config
3. Try to use defaultConfigParamLoader to load from config pointed by tag

It will return with error if both way are failed.
```
*/
func (d *defaultParamLoader) Load(sd *autowire.StructDescriber, fi *autowire.FieldInfo) (interface{}, error) {
	if param, err := d.defaultTagPointToParamLoader.Load(sd, fi); err == nil {
		return param, nil
	}
	if param, err := d.defaultTagParamLoader.Load(sd, fi); err == nil {
		return param, nil
	}
	return d.defaultConfigParamLoader.Load(sd, fi)
}

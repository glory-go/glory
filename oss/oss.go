package oss

import (
	"errors"
)

import (
	"github.com/glory-go/glory/config"
)

func init() {
	defaultOssHandler = newOssHandler()
	defaultOssHandler.setup(config.GlobalServerConf.OssConfigs)
}

var defaultOssHandler *ossHandler

type ossHandler struct {
	services map[string]OssService
}

func (oh *ossHandler) setup(configs map[string]*config.OssConfig) {
	for k, v := range configs {
		switch v.OssType {
		case "qiniu":
			oh.services[k] = newQiniuService()
		default:
			panic("unsupport Oss type, only support qiniu")
		}
		oh.services[k].loadConfig(config.GlobalServerConf.OssConfigs[k])
		oh.services[k].setup()
	}
}

func newOssHandler() *ossHandler {
	return &ossHandler{
		services: make(map[string]OssService),
	}
}

func GetOssService(ossServiceKey string) (OssService, error) {
	if service, ok := defaultOssHandler.services[ossServiceKey]; ok {
		return service, nil
	}
	return nil, errors.New("oss service key not exist")
}

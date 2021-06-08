package service

import (
	"context"

	"github.com/glory-go/glory/log"

	"github.com/glory-go/glory/common/invoker_impl"
	"github.com/glory-go/glory/plugin"

	"github.com/glory-go/glory/config"
)

type httpProtocolService struct {
	serviceBase
	httpProtocolServiceImpl interface{}
}

func NewHTTPProtocolService(name string, serviceProvider interface{}) *httpProtocolService {
	newHTTPProtocolService := &httpProtocolService{
		httpProtocolServiceImpl: serviceProvider,
	}
	newHTTPProtocolService.name = name
	newHTTPProtocolService.loadConfig(config.GlobalServerConf.ServiceConfigs[name])
	return newHTTPProtocolService
}

func (ts *httpProtocolService) Run(ctx context.Context) {
	protoc := plugin.GetProtocol(ts.conf.protocol, nil) // http protocol
	invoker := invoker_impl.NewInvokerFromProvider(ts.httpProtocolServiceImpl)
	if err := protoc.Export(ctx, invoker, ts.conf.addr); err != nil {
		log.Error("protoc Export error")
	}
}

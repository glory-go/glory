package service

import (
	"context"

	"github.com/glory-go/glory/log"

	"github.com/glory-go/glory/common/invoker_impl"
	"github.com/glory-go/glory/plugin"

	"github.com/glory-go/glory/config"
)

type tripleService struct {
	serviceBase
	tripleServiceImpl interface{}
}

func NewTripleService(name string, serviceProvider interface{}) *tripleService {
	newgrpcService := &tripleService{
		tripleServiceImpl: serviceProvider,
	}
	newgrpcService.name = name
	newgrpcService.loadConfig(config.GlobalServerConf.ServiceConfigs[name])
	return newgrpcService
}

func (ts *tripleService) Run(ctx context.Context) {
	protoc := plugin.GetProtocol(ts.conf.protocol, nil, ts.tripleServiceImpl) // glory protocol
	invoker := invoker_impl.NewInvokerFromProvider(ts.tripleServiceImpl)
	if err := protoc.Export(ctx, invoker, ts.conf.addr); err != nil {
		log.Error("protoc Export error")
	}
}

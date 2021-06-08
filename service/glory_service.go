package service

import (
	"context"

	"github.com/glory-go/glory/log"

	"github.com/glory-go/glory/common/invoker_impl"
	"github.com/glory-go/glory/plugin"

	"github.com/glory-go/glory/config"
)

type gloryService struct {
	serviceBase
	gloryServiceImpl interface{}
}

func NewGloryService(name string, serviceProvider interface{}) *gloryService {
	newgrpcService := &gloryService{
		gloryServiceImpl: serviceProvider,
	}
	newgrpcService.name = name
	newgrpcService.loadConfig(config.GlobalServerConf.ServiceConfigs[name])
	return newgrpcService
}

func (gs *gloryService) Run(ctx context.Context) {
	protoc := plugin.GetProtocol(gs.conf.protocol, nil) // glory protocol
	invoker := invoker_impl.NewInvokerFromProvider(gs.gloryServiceImpl)
	if err := protoc.Export(ctx, invoker, gs.conf.addr); err != nil {
		log.Error("protoc Export error")
	}
}

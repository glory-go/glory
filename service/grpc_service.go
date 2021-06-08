package service

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/glory-go/glory/filter/intercepter_impl"

	"github.com/glory-go/glory/config"

	"google.golang.org/grpc"
)

type GrpcService struct {
	serviceBase
	grpcServer *grpc.Server
}

func NewGrpcService(name string) *GrpcService {
	newgrpcService := &GrpcService{}
	newgrpcService.name = name
	newgrpcService.loadConfig(config.GlobalServerConf.ServiceConfigs[name])
	newgrpcService.setup()
	return newgrpcService
}

func (gs *GrpcService) setup() {
	gs.grpcServer = grpc.NewServer(getOptionFromFilter(gs.conf.filtersKey)...)
}

func (gs *GrpcService) Run(ctx context.Context) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", gs.conf.addr.Port))
	if err != nil {
		log.Fatalf("failed to listen grpc: %v", err)
	}
	fmt.Println("grpc start listening on", fmt.Sprintf(":%v", gs.conf.addr.Port))
	if err := gs.grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// GetGrpcServer 对用户暴露的接口
func (gs *GrpcService) GetGrpcServer() *grpc.Server {
	return gs.grpcServer
}

func getOptionFromFilter(filterKeys []string) []grpc.ServerOption {
	intercepter, err := intercepter_impl.NewDefaultGRPCIntercepter(filterKeys)
	if err != nil {
		panic(err)
	}
	return []grpc.ServerOption{
		grpc.UnaryInterceptor(intercepter.ServerIntercepterHandle),
	}
}

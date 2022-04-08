package service

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/glory-go/glory/service/middleware/common"
	"github.com/glory-go/glory/service/middleware/jaeger"

	"github.com/glory-go/glory/config"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

type GrpcService struct {
	serviceBase
	grpcServer *grpc.Server

	unaryMWs []grpc.UnaryServerInterceptor
}

func NewGrpcService(name string) *GrpcService {
	newgrpcService := &GrpcService{}
	newgrpcService.name = name
	newgrpcService.loadConfig(config.GlobalServerConf.ServiceConfigs[name])
	newgrpcService.setup()
	return newgrpcService
}

func (gs *GrpcService) setup() {
	gs.unaryMWs = make([]grpc.UnaryServerInterceptor, 0)
	gs.RegisterUnaryInterceptor(jaeger.UnaryServerMW())
	gs.RegisterUnaryInterceptor(common.GetUnaryServerMWs()...)
}

func (gs *GrpcService) RegisterUnaryInterceptor(mw ...grpc.UnaryServerInterceptor) {
	gs.unaryMWs = append(gs.unaryMWs, mw...)
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
	gs.grpcServer = grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(gs.unaryMWs...),
		),
	)
	return gs.grpcServer
}

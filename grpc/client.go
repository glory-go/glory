package grpc

import (
	"github.com/glory-go/glory/common"
	"github.com/glory-go/glory/config"
	_ "github.com/glory-go/glory/filter/filter_impl"
	"github.com/glory-go/glory/filter/intercepter_impl"
	_ "github.com/glory-go/glory/grpc/resolver"
	"github.com/glory-go/glory/log"
	mwcomm "github.com/glory-go/glory/service/middleware/common"
	"github.com/glory-go/glory/service/middleware/jaeger"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/balancer/roundrobin"
)

type GrpcClient struct {
	conn          *grpc.ClientConn
	serverID      string
	clientName    string
	targetAddress *common.Address
	schema        string
}

func (gc *GrpcClient) setTargetAddr(address *common.Address) {
	gc.targetAddress = address
}

func (gc *GrpcClient) setTargetServerID(serverID string) {
	gc.serverID = serverID
}

func (gc *GrpcClient) setSchema(schema string) {
	gc.schema = schema
}

func (gc *GrpcClient) setClientName(clientName string) {
	gc.clientName = clientName
}

func (gc *GrpcClient) setup(unaryMWs ...grpc.UnaryClientInterceptor) {
	var err error
	dialOption := []grpc.DialOption{grpc.WithInsecure()}
	// add client middlewares
	unaryMWs = append(unaryMWs, jaeger.UnaryClientMW())
	unaryMWs = append(unaryMWs, mwcomm.GetUnaryClientMWs()...)
	dialOption = append(dialOption,
		grpc.WithUnaryInterceptor(
			grpc_middleware.ChainUnaryClient(unaryMWs...),
		),
	)

	if gc.targetAddress != nil {
		// no need service discovery
		gc.conn, err = grpc.Dial(gc.targetAddress.GetUrl(), dialOption...)
	} else {
		dialOption = addDialOptionsWithLoadBalancer(dialOption)
		dialOption = addDialOptionsWithSchemaResolver(dialOption, gc.schema)
		gc.conn, err = grpc.Dial(gc.schema+":///"+gc.clientName, dialOption...)
	}
	if err != nil {
		log.Error(err)
	}
}

// GetConn 返回连接好的grpc.ClientConn指针，用于pb文件注册
func (gc *GrpcClient) GetConn() *grpc.ClientConn {
	return gc.conn
}

// addDialOptionsWithFilters 根据filter返回对应的DialOption
func addDialOptionsWithFilters(opts []grpc.DialOption, filterKeys []string) []grpc.DialOption {
	intercepter, err := intercepter_impl.NewDefaultGRPCIntercepter(filterKeys)
	if err != nil {
		panic(err)
	}
	return append(opts, grpc.WithUnaryInterceptor(intercepter.ClientIntercepterHandle))
}

//addDialOptionsWithLoadBalancer 根据opt返回对应的DialOption
func addDialOptionsWithLoadBalancer(opts []grpc.DialOption) []grpc.DialOption {
	return append(opts, grpc.WithDefaultServiceConfig(`{"LoadBalancingPolicy": "round_robin"}`)) //,
}

//addDialOptionsWithSchemaResolver 增加对应的DialOption
func addDialOptionsWithSchemaResolver(opts []grpc.DialOption, schema string) []grpc.DialOption {
	return append(opts, grpc.WithResolvers(NewResolverBuilder(schema))) //,
}

func NewGrpcClient(grpcClientName string, unaryMWs ...grpc.UnaryClientInterceptor) *GrpcClient {
	//var targetAddress []common.Address
	// get glory client config
	gloryClientConfig, ok := config.GlobalServerConf.ClientConfig[grpcClientName]
	if !ok {
		panic("glory serviceName " + grpcClientName + " in your source code not found in config file!")
	}

	grpcClient := &GrpcClient{}
	grpcClient.setClientName(grpcClientName)
	if gloryClientConfig.ServerAddress != "" {
		grpcClient.setTargetAddr(common.NewAddress(gloryClientConfig.ServerAddress))
	}
	grpcClient.setTargetServerID(gloryClientConfig.ServiceID)
	if regConf, ok := config.GlobalServerConf.RegistryConfig[gloryClientConfig.RegistryKey]; ok {
		grpcClient.setSchema(regConf.Service)
	}
	grpcClient.setup(unaryMWs...) // filter keys used by grpc client
	return grpcClient
}

func NewGrpcClientWithDynamicAddr(grpcClientName string, addr string, unaryMWs ...grpc.UnaryClientInterceptor) *GrpcClient {
	gloryClientConfig, ok := config.GlobalServerConf.ClientConfig[grpcClientName]
	if !ok {
		panic("glory serviceName " + grpcClientName + " in your source code not found in config file!")
	}

	grpcClient := &GrpcClient{}
	grpcClient.setClientName(grpcClientName)
	grpcClient.setTargetAddr(common.NewAddress(addr))
	grpcClient.setTargetServerID(gloryClientConfig.ServiceID)
	if regConf, ok := config.GlobalServerConf.RegistryConfig[gloryClientConfig.RegistryKey]; ok {
		if regConf.Service == "k8s" {
			grpcClient.setSchema("k8s")
		}
	}
	grpcClient.setup(unaryMWs...) // filter keys used by grpc client
	return grpcClient
}

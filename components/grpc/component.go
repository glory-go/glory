package grpc

import (
	"fmt"
	"sync"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	GRPCComponentName = "grpc"
)

type grpcComponent struct {
	config             map[string]*grpcConfig
	conns              map[string]*grpc.ClientConn
	options            map[string][]grpc.DialOption
	unaryInterceptors  map[string][]grpc.UnaryClientInterceptor
	streamInterceptors map[string][]grpc.StreamClientInterceptor
}

var (
	component *grpcComponent
	once      sync.Once
)

func getGrpcComponent() *grpcComponent {
	once.Do(func() {
		component = &grpcComponent{
			config:  map[string]*grpcConfig{},
			conns:   make(map[string]*grpc.ClientConn),
			options: make(map[string][]grpc.DialOption),
		}
	})
	return component
}

// WithOption 提供使用者指定grpc客户端初始化时使用的option的能力。name为空代表为所有的客户端端进行注册
// UnaryInterceptor和StreamInterceptor的注册，请使用WithXXXInterceptors方法，否则会导致运行时panic
func WithOptions(name string, options ...grpc.DialOption) {
	comp := getGrpcComponent()
	comp.options[name] = append(comp.options[name], options...)
}

// WithUnaryInterceptors 提供使用者指定grpc客户端初始化时使用的UnaryInterceptor的能力。name为空代表为所有的客户端端进行注册
func WithUnaryInterceptors(name string, incs ...grpc.UnaryClientInterceptor) {
	comp := getGrpcComponent()
	comp.unaryInterceptors[name] = append(comp.unaryInterceptors[name], incs...)
}

// WithStreamInterceptors 提供使用者指定grpc客户端初始化时使用的StreamInterceptor的能力。name为空代表为所有的客户端端进行注册
func WithStreamInterceptors(name string, incs ...grpc.StreamClientInterceptor) {
	comp := getGrpcComponent()
	comp.streamInterceptors[name] = append(comp.streamInterceptors[name], incs...)
}

func GetGRPCClient(name string) *grpc.ClientConn {
	return getGrpcComponent().conns[name]
}

func (c *grpcComponent) Name() string { return GRPCComponentName }

func (c *grpcComponent) Init(config map[string]any) error {
	for name := range config {
		raw := config[name]
		conf := &grpcConfig{}
		if err := mapstructure.Decode(raw, conf); err != nil {
			return err
		}
		// 检查配置
		if conf.Host == "" || conf.Port == 0 {
			return fmt.Errorf("grpcComponent/Init: invalid host %s or port %d for grpc client %s", conf.Host, conf.Port, name)
		}
		// 初始化客户端连接
		options := append(c.options[name], c.options[""]...)
		options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))
		unaryInterceptors := append(c.unaryInterceptors[name], c.unaryInterceptors[""]...)
		streamInterceptors := append(c.streamInterceptors[name], c.streamInterceptors[""]...)
		options = append(options,
			grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(streamInterceptors...)),
			grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(unaryInterceptors...)),
		)
		conn, err := grpc.Dial(fmt.Sprintf("%s:%d", conf.Host, conf.Port), options...)
		if err != nil {
			return err
		}
		c.config[name] = conf
		c.conns[name] = conn
	}

	return nil
}

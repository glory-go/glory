package grpc

import (
	"fmt"
	"sync"

	"github.com/mitchellh/mapstructure"
	"google.golang.org/grpc"
)

const (
	GRPCComponentName = "grpc"
)

type grpcComponent struct {
	config  map[string]*grpcConfig
	conns   map[string]*grpc.ClientConn
	options map[string][]grpc.DialOption
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
func WithOptions(name string, options ...grpc.DialOption) {
	comp := getGrpcComponent()
	comp.options[name] = append(comp.options[name], options...)
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
		options := append(c.options[""], c.options[name]...)
		conn, err := grpc.Dial(fmt.Sprintf("%s:%d", conf.Host, conf.Port), options...)
		if err != nil {
			return err
		}
		c.config[name] = conf
		c.conns[name] = conn
	}

	return nil
}

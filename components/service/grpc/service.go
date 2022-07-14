package grpc

import (
	"net"
	"sync"

	"github.com/mitchellh/mapstructure"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

const (
	GRPCServiceName = "grpc"
)

type grpcService struct {
	configs map[string]*grpcServiceConfig
	servers map[string]*grpc.Server
	options map[string][]grpc.ServerOption
	g       errgroup.Group
}

var (
	srv  *grpcService
	once sync.Once
)

func GetGRPCService() *grpcService {
	once.Do(func() {
		srv = &grpcService{
			configs: make(map[string]*grpcServiceConfig),
			servers: make(map[string]*grpc.Server),
			options: make(map[string][]grpc.ServerOption),
		}
	})

	return srv
}

// WithOption 提供使用者指定grpc服务初始化时使用的option的能力。name为空代表为所有的服务端进行注册
func (s *grpcService) WithOptions(name string, options ...grpc.ServerOption) {
	s.options[name] = append(s.options[name], options...)
}

func (s *grpcService) GetServer(name string) *grpc.Server {
	return s.servers[name]
}

func (s *grpcService) Name() string { return GRPCServiceName }

func (s *grpcService) Init(config map[string]interface{}) error {
	for k, v := range config {
		// 解析配置
		conf := &grpcServiceConfig{}
		if err := mapstructure.Decode(v, conf); err != nil {
			return err
		}
		s.configs[k] = conf
		// 初始化服务
		options := append(s.options[""], s.options[k]...)
		s.servers[k] = grpc.NewServer(options...)
	}
	return nil
}

func (s *grpcService) Run() error {
	for k := range s.servers {
		name := k
		s.g.Go(func() error {
			lis, err := net.Listen("tcp", s.configs[name].Addr)
			if err != nil {
				return err
			}
			return s.servers[name].Serve(lis)
		})
	}

	return s.g.Wait()
}

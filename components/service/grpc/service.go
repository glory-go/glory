package grpc

import (
	"fmt"
	"log"
	"net"
	"sync"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	GRPCServiceName = "grpc"
)

type grpcService struct {
	configs            map[string]*grpcServiceConfig
	servers            map[string]*grpc.Server
	options            map[string][]grpc.ServerOption
	unaryInterceptors  map[string][]grpc.UnaryServerInterceptor
	streamInterceptors map[string][]grpc.StreamServerInterceptor
	g                  errgroup.Group
}

var (
	srv  *grpcService
	once sync.Once
)

func getGRPCService() *grpcService {
	once.Do(func() {
		srv = &grpcService{
			configs:            make(map[string]*grpcServiceConfig),
			servers:            make(map[string]*grpc.Server),
			options:            make(map[string][]grpc.ServerOption),
			unaryInterceptors:  make(map[string][]grpc.UnaryServerInterceptor),
			streamInterceptors: make(map[string][]grpc.StreamServerInterceptor),
		}
	})

	return srv
}

// WithOption 提供使用者指定grpc服务初始化时使用的option的能力。name为空代表为所有的服务端进行注册.
// UnaryInterceptor和StreamInterceptor的注册，请使用WithXXXInterceptors方法，否则会导致运行时panic
func WithOptions(name string, options ...grpc.ServerOption) {
	s := getGRPCService()
	s.options[name] = append(s.options[name], options...)
}

// WithUnaryInterceptors 提供使用者指定grpc服务初始化时使用的UnaryInterceptor的能力。name为空代表为所有的服务端进行注册
func WithUnaryInterceptors(name string, incs ...grpc.UnaryServerInterceptor) {
	s := getGRPCService()
	s.unaryInterceptors[name] = append(s.unaryInterceptors[name], incs...)
}

// WithStreamInterceptors 提供使用者指定grpc服务初始化时使用的StreamInterceptor的能力。name为空代表为所有的服务端进行注册
func WithStreamInterceptors(name string, incs ...grpc.StreamServerInterceptor) {
	s := getGRPCService()
	s.streamInterceptors[name] = append(s.streamInterceptors[name], incs...)
}

func GetServer(name string) *grpc.Server {
	register()
	return getGRPCService().servers[name]
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
		options := append(s.options[k], s.options[""]...)
		unaryInterceptors := append(s.unaryInterceptors[k], s.unaryInterceptors[""]...)
		streamInterceptors := append(s.streamInterceptors[k], s.streamInterceptors[""]...)
		options = append(options,
			grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(streamInterceptors...)),
			grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(unaryInterceptors...)),
		)
		if conf.TLS {
			creds, err := credentials.NewServerTLSFromFile(conf.CertFile, conf.KeyFile)
			if err != nil {
				return fmt.Errorf("grpcService/Init: failed to create server TLS creds, %s", err)
			}
			options = append(options, grpc.Creds(creds))
		}
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
			log.Printf("grpc service %s listening at %s\n", name, lis.Addr().String())
			return s.servers[name].Serve(lis)
		})
	}

	return s.g.Wait()
}

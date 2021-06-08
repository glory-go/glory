package server

import (
	"context"
	"sync"

	"github.com/glory-go/glory/log"

	gostNet "github.com/dubbogo/gost/net"

	"github.com/glory-go/glory/common"
	"github.com/glory-go/glory/plugin"

	"github.com/glory-go/glory/config"

	"github.com/glory-go/glory/service"
)

type GloryServer interface {
	Run()
	LoadConfig()
	RegisterService(service service.Service)
}

type DefaultGloryServer struct {
	Services     map[string]service.Service
	GloryService map[string]interface{}
	GloryClient  map[string]interface{}
	ServerConfig *config.ServerConfig

	wg      sync.WaitGroup
	ctx     context.Context
	localIp string
}

var defaultGloryServer *DefaultGloryServer

// LoadConfig load configure file
func (s *DefaultGloryServer) LoadConfig() {
	s.ServerConfig = config.GlobalServerConf
}

// RegisterService Register other manual service
func (s *DefaultGloryServer) RegisterService(service service.Service) {
	s.Services[service.GetName()] = service
}

// Run run manual service and load and run auto serivce
func (s *DefaultGloryServer) Run() {
	// load auto service:
	for _, v := range s.Services {
		// 对于注册好的每个service，都要1 服务注册、 2 开启监听
		// service registry procedure
		// get export Addr
		lstnAddress := common.Address{
			Host: s.localIp,
			Port: v.GetPort(),
		}
		// glory_registry protocol do register
		registryKey := v.GetRegistryKey()
		if registryKey != "" {
			registryConfig, ok := config.GlobalServerConf.RegistryConfig[registryKey]
			if !ok {
				panic("serverConfig.RegistryKey = " + registryKey + " not defined in registry block")
			}
			registryProtoc := plugin.GetRegistry(registryConfig)
			if registryProtoc == nil {
				log.Errorf("get registry protocol failed with registryKey = %s", registryKey)
			} else {
				registryProtoc.Register(v.GetServiceID(), lstnAddress)
				go func() {
					select {
					case <-s.ctx.Done():
						registryProtoc.UnRegister(v.GetServiceID(), lstnAddress)
						return
					}
				}()
			}

		}
		v.SetListeningAddr(lstnAddress)

		s.wg.Add(1)
		go func(s service.Service, wg *sync.WaitGroup, ctx context.Context) {
			defer func() {
				if e := recover(); e != nil {
					log.Error("error :", e)
				}
				wg.Done()
			}()
			s.Run(ctx)
		}(v, &s.wg, s.ctx)
	}
	s.wg.Wait()
}

func NewDefaultGloryServer(ctx context.Context) *DefaultGloryServer {
	host, err := gostNet.GetLocalIP()
	if err != nil {
		panic("get server ip err: " + err.Error())
	}
	return &DefaultGloryServer{
		Services:     make(map[string]service.Service),
		ServerConfig: config.NewServerConfig(),
		GloryService: make(map[string]interface{}),
		GloryClient:  make(map[string]interface{}),
		ctx:          ctx,
		wg:           sync.WaitGroup{},
		localIp:      host,
	}
}

func GetDefaultServer() *DefaultGloryServer {
	return defaultGloryServer
}

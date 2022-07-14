package ginhttp

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/sync/errgroup"
)

const (
	GinHTTPServiceName = "ginhttp"
)

type ginHttpService struct {
	routers map[string]*gin.Engine
	configs map[string]*ginHttpServiceConfig
	g       errgroup.Group
}

var (
	srv  *ginHttpService
	once sync.Once
)

func GetGinHttpService() *ginHttpService {
	once.Do(func() {
		srv = &ginHttpService{
			routers: make(map[string]*gin.Engine),
			configs: make(map[string]*ginHttpServiceConfig),
		}
	})

	return srv
}

func (s *ginHttpService) GetEngine(name string) *gin.Engine {
	return s.routers[name]
}

func (s *ginHttpService) Name() string { return GinHTTPServiceName }

func (s *ginHttpService) Init(config map[string]interface{}) error {
	for k, v := range config {
		s.routers[k] = gin.Default()
		// 解析配置
		conf := &ginHttpServiceConfig{}
		if err := mapstructure.Decode(v, conf); err != nil {
			return err
		}
		s.configs[k] = conf
	}
	return nil
}

func (s *ginHttpService) Run() error {
	for k := range s.routers {
		name := k
		s.g.Go(func() error {
			server := http.Server{
				Addr:         s.configs[name].Addr,
				Handler:      s.routers[name],
				ReadTimeout:  time.Duration(s.configs[name].ReadTimeout) * time.Second,
				WriteTimeout: time.Duration(s.configs[name].WriteTimeout) * time.Second,
			}
			return server.ListenAndServe()
		})
	}

	return s.g.Wait()
}

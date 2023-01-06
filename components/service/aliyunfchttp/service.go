package aliyunfchttp

import (
	"context"
	"net/http"
	"sync"

	"github.com/aliyun/fc-runtime-go-sdk/fc"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
)

const (
	AliyunFCHTTPServiceName = "aliyun_fc_http"
)

type aliyunFCHttpService struct {
	router *gin.Engine
	config *aliyunFCHttpServiceConfig
}

var (
	srv  *aliyunFCHttpService
	once sync.Once
)

func getAliyunFCHttpService() *aliyunFCHttpService {
	once.Do(func() {
		srv = &aliyunFCHttpService{}
	})

	return srv
}

func GetEngine() *gin.Engine {
	return getAliyunFCHttpService().router
}

func (s *aliyunFCHttpService) Name() string { return AliyunFCHTTPServiceName }

func (s *aliyunFCHttpService) Init(config map[string]interface{}) error {
	s.router = gin.Default()
	s.router.ContextWithFallback = true
	// 解析配置
	conf := &aliyunFCHttpServiceConfig{}
	if err := mapstructure.Decode(config, conf); err != nil {
		return err
	}
	s.config = conf
	return nil
}

func (s *aliyunFCHttpService) Run() error {
	fc.StartHttp(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		s.router.ServeHTTP(w, r.Clone(ctx))
		return nil
	})
	return nil
}

package http

import (
	"github.com/gin-gonic/gin"
	"github.com/glory-go/glory/config"
	"github.com/glory-go/glory/service"
	"github.com/glory-go/glory/service/v2/http/middleware"
)

type HttpService struct {
	service.DefaultServiceBase

	engine *gin.Engine
}

func NewHttpService(name string) *HttpService {
	httpService := &HttpService{}
	httpService.Name = name
	httpService.LoadConfig(config.GlobalServerConf.ServiceConfigs[name])
	httpService.setup()
	return httpService
}

func (hs *HttpService) setup() {
	engine := gin.Default()
	engine.Use(
		middleware.NewMTraceMW().HandlerFunc,
		middleware.GLoggerMW{}.HandlerFunc,
	)

	hs.engine = engine
}

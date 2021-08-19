package service

import (
	"context"
	"net/http"
	"strconv"

	"github.com/glory-go/glory/config"
	ghttp "github.com/glory-go/glory/http"
	"github.com/glory-go/glory/service/middleware/jaeger"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

type HttpService struct {
	serviceBase
	router *mux.Router

	mws []negroni.Handler
}

func NewHttpService(name string) *HttpService {
	httpService := &HttpService{}
	httpService.name = name
	httpService.loadConfig(config.GlobalServerConf.ServiceConfigs[name])
	httpService.setup()
	return httpService
}

func (hs *HttpService) setup() {
	hs.router = mux.NewRouter()
	hs.UseMW(&jaeger.AliyunJaegerMW{})
}

func (hs *HttpService) UseMW(filters ...negroni.Handler) {
	hs.mws = append(hs.mws, filters...)
}

func (hs *HttpService) Run(ctx context.Context) {
	// handler := cors.Default().Handler(hs.router)
	s := negroni.Classic()
	for _, handler := range hs.mws {
		s.Use(handler)
	}
	s.UseHandler(hs.router)
	s.Run(":" + strconv.Itoa(hs.conf.addr.Port))
}

func (hs *HttpService) RegisterRouterWithRawHttpHandler(path string, handler func(w http.ResponseWriter, r *http.Request), method string) {
	hs.router.HandleFunc(path, handler).Methods(method)
}

// RegisterRouter 对用户暴露的接口
func (hs *HttpService) RegisterRouter(path string, handler func(*ghttp.GRegisterController) error, req interface{},
	rsp interface{}, method string, filters ...ghttp.Filter) {
	ghttp.RegisterRouter(path, hs.router, handler, req, rsp, method, filters)
}

// RegisterWSRouter 对用户暴露的接口
func (hs *HttpService) RegisterWSRouter(path string, handler func(*ghttp.GRegisterWSController)) {
	ghttp.RegisterWSRouter(path, hs.router, handler)
}

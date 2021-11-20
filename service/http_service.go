package service

import (
	"context"
	"net/http"
	"net/http/pprof"
	"strconv"

	"github.com/glory-go/glory/config"
	ghttp "github.com/glory-go/glory/http"
	"github.com/glory-go/glory/service/middleware/jaeger"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/negroni"
)

type HttpService struct {
	serviceBase
	router *mux.Router

	mws      []negroni.Handler
	gloryMWs []ghttp.Filter
}

func NewHttpService(name string, enableMetrics bool) *HttpService {
	httpService := &HttpService{}
	httpService.name = name
	httpService.loadConfig(config.GlobalServerConf.ServiceConfigs[name])
	httpService.setup()

	if enableMetrics {
		httpService.RegisterRouterWithRawHttpHandler("/metrics", promhttp.Handler().ServeHTTP, http.MethodGet)
		httpService.RegisterRouterWithRawHttpHandler("/debug/pprof/{action}", pprof.Index)
		httpService.RegisterRouterWithRawHttpHandler("/debug/pprof/", pprof.Index)
		httpService.RegisterRouterWithRawHttpHandler("/debug/pprof/symbol", pprof.Symbol)
	}
	return httpService
}

func (hs *HttpService) setup() {
	hs.router = mux.NewRouter()

	aliyunJaegerInstance := &jaeger.AliyunJaegerMW{}
	hs.UseGloryMW(aliyunJaegerInstance.GloryMW)
}

func (hs *HttpService) UseMW(filters ...negroni.Handler) {
	hs.mws = append(hs.mws, filters...)
}

func (hs *HttpService) UseGloryMW(filters ...ghttp.Filter) {
	hs.gloryMWs = append(hs.gloryMWs, filters...)
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

func (hs *HttpService) RegisterRouterWithRawHttpHandler(path string, handler func(w http.ResponseWriter, r *http.Request), method ...string) {
	r := hs.router.HandleFunc(path, handler)
	// 未传method则match all
	if len(method) != 0 {
		r.Methods(method...)
	}
}

// RegisterRouter 对用户暴露的接口
func (hs *HttpService) RegisterRouter(path string, handler func(*ghttp.GRegisterController) error, req interface{},
	rsp interface{}, method string, filters ...ghttp.Filter) {
	filters = append(hs.gloryMWs, filters...)
	ghttp.RegisterRouter(path, hs.router, handler, req, rsp, method, filters)
}

// RegisterWSRouter 对用户暴露的接口
func (hs *HttpService) RegisterWSRouter(path string, handler func(*ghttp.GRegisterWSController)) {
	ghttp.RegisterWSRouter(path, hs.router, handler)
}

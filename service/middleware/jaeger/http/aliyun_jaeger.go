package http

import (
	"net/http"

	"github.com/glory-go/glory/config"
	"github.com/glory-go/glory/log"
	"github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/transport"
)

type aliyunJaegerConfig struct {
	Endpoint string
}

type AliyunJaegerMW struct{}

func (m *AliyunJaegerMW) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	funcName := "AliyunJaegerMW"
	viper := config.GetViperConfig()
	jaegerConfig := &aliyunJaegerConfig{}
	if err := viper.Unmarshal(jaegerConfig); err != nil {
		log.CtxWarnf(r.Context(), "[%v] fail to parse config with err %v", funcName, err)
		next(rw, r)
		return
	}
	// 初始化jaeger实例
	sender := transport.NewHTTPTransport(
		jaegerConfig.Endpoint,
	)
	tracer, closer := jaeger.NewTracer(config.GlobalServerConf.GetAppKey(),
		jaeger.NewConstSampler(true),
		jaeger.NewRemoteReporter(sender),
	)
	// 生成span
	span := opentracing.SpanFromContext(r.Context())
	if span == nil {
		log.CtxWarnf(r.Context(), "[%v] fail to get span from request", funcName)
	}
	tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
	closer.Close()
	
	next(rw, r)
}

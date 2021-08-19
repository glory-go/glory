package jaeger

import (
	"net/http"

	"github.com/glory-go/glory/log"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
)

type aliyunJaegerConfig struct {
	Endpoint string
}

type AliyunJaegerMW struct{}

func (m *AliyunJaegerMW) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	nethttp.MiddlewareFunc(tracer, next)(rw, r)
	// 将trace-id写到header
	span := opentracing.SpanFromContext(r.Context())
	if span == nil {
		return
	}
	sc, ok := span.Context().(jaeger.SpanContext)
	if ok {
		return
	}
	rw.Header().Set(log.GetTraceIDKey(), sc.TraceID().String())
}

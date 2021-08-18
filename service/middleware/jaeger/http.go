package jaeger

import (
	"net/http"

	"github.com/opentracing-contrib/go-stdlib/nethttp"
)

type aliyunJaegerConfig struct {
	Endpoint string
}

type AliyunJaegerMW struct{}

func (m *AliyunJaegerMW) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	nethttp.MiddlewareFunc(tracer, next)(rw, r)
}

package jaeger

import (
	"context"
	"net/http"

	ghttp "github.com/glory-go/glory/http"
	"github.com/glory-go/glory/log"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
)

type AliyunJaegerMW struct{}

func addTraceID2ResHeader(ctx context.Context, rw http.ResponseWriter) {
	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		return
	}
	sc, ok := span.Context().(jaeger.SpanContext)
	if !ok {
		return
	}
	rw.Header().Set(log.GetTraceIDKey(), sc.TraceID().String())
}

func (m *AliyunJaegerMW) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	nethttp.MiddlewareFunc(tracer, next)(rw, r)
	// 将trace-id写到header
	addTraceID2ResHeader(r.Context(), rw)
}

func (m *AliyunJaegerMW) GloryMW(c *ghttp.GRegisterController, f ghttp.HandleFunc) (err error) {
	nethttp.MiddlewareFunc(tracer, func(rw http.ResponseWriter, r *http.Request) {
		c.R = r
		c.W = rw
		err = f(c)
	})(c.W, c.R)
	addTraceID2ResHeader(c.R.Context(), c.W)
	return
}

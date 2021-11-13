package jaeger

import (
	"github.com/glory-go/glory/log"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"google.golang.org/grpc"
)

func UnaryServerMW() grpc.UnaryServerInterceptor {
	return grpc_opentracing.UnaryServerInterceptor(
		grpc_opentracing.WithTracer(tracer),
	)
}

func UnaryClientMW() grpc.UnaryClientInterceptor {
	return grpc_opentracing.UnaryClientInterceptor(
		grpc_opentracing.WithTracer(tracer),
		grpc_opentracing.WithTraceHeaderName(log.GetTraceIDKey()),
	)
}

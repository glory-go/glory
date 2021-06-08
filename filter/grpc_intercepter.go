package filter

import (
	"context"

	"google.golang.org/grpc"
)

// Intercepter is the grpc invocation api, and is the entrance of glory filter
type Intercepter interface {
	ServerIntercepterHandle(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error)

	ClientIntercepterHandle(ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption) error
}

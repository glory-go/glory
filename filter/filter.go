package filter

import (
	"context"
)

import (
	"google.golang.org/grpc"
)

// GRPCFilter is the normal grpc filter interface
type GRPCFilter interface {
	SetNext(filter GRPCFilter)

	ServerHandle(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo) (resp interface{}, err error)

	ClientHandle(ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption) error
}

package intercepter_impl

import (
	"context"
)

import (
	"google.golang.org/grpc"
)

import (
	"github.com/glory-go/glory/config"
	"github.com/glory-go/glory/filter"
	"github.com/glory-go/glory/log"
	"github.com/glory-go/glory/plugin"
)

func NewDefaultGRPCIntercepter(filterKeys []string) (filter.Intercepter, error) {
	chainFilter, err := plugin.GetFilter("chain", &config.FilterConfig{
		SubFiltersKey: filterKeys,
	})
	if err != nil {
		log.Errorf("new default grpc intercepter error = %v\n", err)
		return nil, err
	}
	return &defaultIntercepter{
		chainFilter: chainFilter,
	}, nil
}

// defaultIntercepter is default grpc intercepter, it calls chain_filter and it self is the final filter of grpc invocation filter procedure
type defaultIntercepter struct {
	chainFilter filter.GRPCFilter
	handler     grpc.UnaryHandler
}

// ServerIntercepterHandle is the calling entrance of
func (j *defaultIntercepter) ServerIntercepterHandle(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp interface{}, err error) {

	j.handler = handler
	j.chainFilter.SetNext(j)
	return j.chainFilter.ServerHandle(ctx, req, info)
}

func (j *defaultIntercepter) ServerHandle(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo) (resp interface{}, err error) {
	return j.handler(ctx, req)
}

func (j *defaultIntercepter) ClientHandle(ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption) error {

	err := invoker(ctx, method, req, reply, cc, opts...)
	return err
}

// SetNext is useless as defaultIntercepter is the end of filter chain, so it doesn't need next filter
func (j *defaultIntercepter) SetNext(filter filter.GRPCFilter) {

}

func (j *defaultIntercepter) ClientIntercepterHandle(ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption) error {

	j.chainFilter.SetNext(j)
	return j.chainFilter.ClientHandle(ctx, method, req, reply, cc, invoker, opts...)
}

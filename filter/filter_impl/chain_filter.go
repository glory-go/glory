package filter_impl

import (
	"context"

	"github.com/glory-go/glory/config"
	"github.com/glory-go/glory/log"

	"github.com/glory-go/glory/plugin"

	"github.com/glory-go/glory/filter"
	"google.golang.org/grpc"
)

// ChainFilter is a filter_impl of GRPC filter
// 通过工厂New出来的filter是没有next的，因为next无法在构造参数中传递，需要在最终调用前，通过SetNext接口写入next参数。否则会报错
type ChainFilter struct {
	chain []filter.GRPCFilter
	next  filter.GRPCFilter
}

func init() {
	plugin.SetFilterFactory("chain", NewChainFilter)
}

func (cf *ChainFilter) ClientHandle(ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption) error {
	if len(cf.chain) == 0 {
		if cf.next == nil {
			log.Errorf("err = %v", filter.ErrNotSetNextFilter)
			return filter.ErrNotSetNextFilter
		}
		return cf.next.ClientHandle(ctx, method, req, reply, cc, invoker, opts...)
	}
	return cf.chain[0].ClientHandle(ctx, method, req, reply, cc, invoker, opts...)
}

func (cf *ChainFilter) ServerHandle(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo) (resp interface{}, err error) {
	if len(cf.chain) == 0 {
		if cf.next == nil {
			log.Errorf("err = %v", filter.ErrNotSetNextFilter)
			return nil, filter.ErrNotSetNextFilter
		}
		return cf.next.ServerHandle(ctx, req, info)
	}
	return cf.chain[0].ServerHandle(ctx, req, info)
}

func (cf *ChainFilter) SetNext(grpcFilter filter.GRPCFilter) {
	cf.next = grpcFilter
	if len(cf.chain) == 0 {
		return
	}
	cf.chain[len(cf.chain)-1].SetNext(grpcFilter)
}

// NewChainFilter never setup failed. If it's subfilter setup with error, chain filter will jump it and setup next.
// The worst result is there is no sub filter successfully setup, and this chain is useless.
func NewChainFilter(filterConfig *config.FilterConfig) (filter.GRPCFilter, error) {
	filterChain := make([]filter.GRPCFilter, 0)
	for _, k := range filterConfig.SubFiltersKey {
		filterConfig, ok := config.GlobalServerConf.FilterConfigMap[k]
		if !ok {
			log.Warnf("filter key %s not defined in config file's filters block, this filter not loaded!", k)
			continue
		}
		tempFilter, err := plugin.GetFilter(filterConfig.FilterName, filterConfig)
		if err != nil {
			log.Warnf("filter key %s setup error = %v", err)
			continue
		}
		filterChain = append(filterChain, tempFilter)
	}

	for i, f := range filterChain {
		if i == len(filterChain)-1 {
			break
		}
		f.SetNext(filterChain[i+1])
	}

	return &ChainFilter{
		chain: filterChain,
	}, nil
}

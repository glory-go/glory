package resolver

import (
	"google.golang.org/grpc/resolver"
)

import (
	"github.com/glory-go/glory/common"
	"github.com/glory-go/glory/log"
	"github.com/glory-go/glory/plugin"
)

func init() {
	plugin.SetGRPCResolverFactory("nacos", NewNacosResolver)
}

func NewNacosResolver(ch chan common.RegistryChangeEvent, cc resolver.ClientConn) resolver.Resolver {
	newResolver := &NacosResolver{
		ch: ch,
		cc: cc,
	}
	go newResolver.watcher()
	return newResolver
}

type NacosResolver struct {
	ch          chan common.RegistryChangeEvent
	cc          resolver.ClientConn
	addressList []resolver.Address
}

func (r *NacosResolver) Close() {
	panic("implement me")
}

func (r *NacosResolver) watcher() {
	// get target from glory_registry
	for {
		e := <-r.ch
		switch e.Opt {
		case common.RegistryUpdateToSerivcesListEvent:
			log.Debugf("nacos add event with addr = %v", e.Serivces)
			r.addressList = make([]resolver.Address, 0)
			for _, v := range e.Serivces {
				r.addressList = append(r.addressList, resolver.Address{Addr: v.GetUrl()})
			}
		}
		if err := r.cc.UpdateState(resolver.State{Addresses: r.addressList}); err != nil {
			log.Errorf("NacosResolver update state failed with error = %s", err)
		}
	}
}
func (r *NacosResolver) ResolveNow(clientName resolver.ResolveNowOptions) {

}

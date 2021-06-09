package resolver

import (
	"github.com/glory-go/glory/common"
	"github.com/glory-go/glory/log"
	"github.com/glory-go/glory/plugin"
	"google.golang.org/grpc/resolver"
)

func init() {
	plugin.SetGRPCResolverFactory("k8s", NewK8SResolver)
}

func NewK8SResolver(ch chan common.RegistryChangeEvent, cc resolver.ClientConn) resolver.Resolver {
	newResolver := &K8SResolver{
		ch: ch,
		cc: cc,
	}
	go newResolver.watcher()
	return newResolver
}

type K8SResolver struct {
	ch          chan common.RegistryChangeEvent
	cc          resolver.ClientConn
	addressList []resolver.Address
	existMap    map[string]bool
}

func (r *K8SResolver) Close() {
	panic("implement me")
}

func (r *K8SResolver) watcher() {
	// get target from glory_registry
	for {
		e := <-r.ch
		switch e.Opt {
		case common.RegistryAddEvent, common.RegistryUpdateEvent:
			log.Debugf("add event with addr = ", e.Addr.GetUrl())
			if _, ok := r.existMap[e.Addr.GetUrl()]; ok {
				continue
			}
			r.addressList = append(r.addressList, resolver.Address{Addr: e.Addr.GetUrl()})
			r.existMap[e.Addr.GetUrl()] = true
		case common.RegistryDeleteEvent:
			log.Debugf("delete event with addr = ", e.Addr.GetUrl())
			if _, ok := r.existMap[e.Addr.GetUrl()]; !ok {
				continue
			}
			newList := make([]resolver.Address, 0)
			for _, v := range r.addressList {
				if v.Addr != e.Addr.GetUrl() {
					newList = append(newList, v)
				}
			}
			delete(r.existMap, e.Addr.GetUrl())
			r.addressList = newList
		}
		r.cc.UpdateState(resolver.State{Addresses: r.addressList})
	}
}
func (r *K8SResolver) ResolveNow(clientName resolver.ResolveNowOptions) {

}

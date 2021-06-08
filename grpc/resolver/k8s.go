package resolver

import (
	"github.com/glory-go/glory/common"
	"github.com/glory-go/glory/config"
	"github.com/glory-go/glory/log"
	"github.com/glory-go/glory/plugin"
	perrors "github.com/pkg/errors"
	"google.golang.org/grpc/resolver"
)

func NewK8SResolverBuilder() resolver.Builder {
	return &K8SResolverBuilder{}
}

type K8SResolverBuilder struct {
}

func (r *K8SResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	// get target from glory_registry
	log.Debugf("in builder.Build(), target = %+v", target)
	clientConf := config.GlobalServerConf.ClientConfig[target.Endpoint]
	registryConfig, ok := config.GlobalServerConf.RegistryConfig[clientConf.RegistryKey]
	if !ok {
		log.Error("not found registry key = ", clientConf.RegistryKey, " in config file")
		return nil, perrors.Errorf("not found registry key = %s in config file", clientConf.RegistryKey)
	}
	reg := plugin.GetRegistry(registryConfig)
	if reg == nil {
		log.Error("get k8s registry failed")
		return nil, perrors.Errorf("get k8s registry failed")
	}
	ch, err := reg.Subscribe(clientConf.ServiceID)
	if err != nil {
		return nil, err
	}

	d := &K8SResolver{
		ch:          ch,
		cc:          cc,
		addressList: make([]resolver.Address, 0),
		existMap:    make(map[string]bool),
	}

	go d.watcher()
	d.ResolveNow(resolver.ResolveNowOptions{})
	return d, nil
}

// Scheme returns the naming scheme of this resolver builder, which is "k8s".
func (b *K8SResolverBuilder) Scheme() string {
	return "k8s"
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

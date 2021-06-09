package grpc

import (
	"github.com/glory-go/glory/config"
	"github.com/glory-go/glory/log"
	"github.com/glory-go/glory/plugin"
	perrors "github.com/pkg/errors"
	"google.golang.org/grpc/resolver"
)

func NewResolverBuilder(typ string) resolver.Builder {
	return &ResolverBuilder{typ: typ}
}

type ResolverBuilder struct {
	typ string
}

func (r *ResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
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

	rsvr := plugin.GetGRPCResolver(r.typ, ch, cc)

	rsvr.ResolveNow(resolver.ResolveNowOptions{})
	return rsvr, nil
}

// Scheme returns the naming scheme of this resolver builder, which is "k8s".
func (b *ResolverBuilder) Scheme() string {
	return b.typ
}

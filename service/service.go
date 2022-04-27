package service

import (
	"context"
)

import (
	"github.com/glory-go/glory/common"
	"github.com/glory-go/glory/config"
	_ "github.com/glory-go/glory/filter/filter_impl"
)

type Service interface {
	GetName() string
	Run(ctx context.Context)
	GetPort() int
	GetRegistryKey() string
	GetServiceID() string
	SetListeningAddr(addr common.Address)

	// @config.ServiceConfig 的 protocol 字段目前没用，是用户手动生成的对应协议的service，因为目前多协议只停留在应用层
	// 以后希望能将多协议下沉到协议层，从而可通过配置文件直接改动协议的选择，无需手动生成service
	loadConfig(conf *config.ServiceConfig)
}

type serviceConf struct {
	addr        common.Address
	RegistryKey string
	ServiceID   string
	filtersKey  []string
	protocol    string
}

type serviceBase struct {
	name string
	conf serviceConf
	//ctx  context.Context
}

func (sb *serviceBase) GetName() string {
	return sb.name
}
func (sb *serviceBase) GetPort() int {
	return sb.conf.addr.Port
}
func (sb *serviceBase) GetRegistryKey() string {
	return sb.conf.RegistryKey
}

func (sb *serviceBase) GetServiceID() string {
	return sb.conf.ServiceID
}

// SetListeningAddr called by server to set server network address
func (sb *serviceBase) SetListeningAddr(addr common.Address) {
	sb.conf.addr = addr
}
func (sb *serviceBase) loadConfig(conf *config.ServiceConfig) {
	sb.conf.ServiceID = conf.ServiceID
	sb.conf.RegistryKey = conf.RegistryKey
	// port是配置文件里开发者手动配置的
	sb.conf.addr.Port = conf.Port // 这里只有port确定了，host还未确定，要等server调用上面的SetListeningAddr才能确定绑定的host
	sb.conf.protocol = conf.Protocol
	sb.conf.filtersKey = conf.FiltersKey
}

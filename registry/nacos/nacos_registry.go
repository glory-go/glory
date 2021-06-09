package nacos

import (
	"fmt"
	"github.com/glory-go/glory/common"
	"github.com/glory-go/glory/config"
	"github.com/glory-go/glory/log"
	"github.com/glory-go/glory/plugin"
	"github.com/glory-go/glory/registry"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"strconv"
)

const (
	NacosTimeoutMs        = 5000
	NacosLogDir           = "."
	NacosCacheDir         = "."
	NacosRotateTime       = "1h"
	GloryNacosNamespaceID = "glory"
	NacosLogLevel         = "debug"
	NacosMaxAge           = 3
)

type nacosRegistry struct {
	client naming_client.INamingClient
}

func init() {
	plugin.SetRegistryFactory("nacos", newNacosRegistry)
}

func newNacosRegistry(registryConfig *config.RegistryConfig) registry.Registry {
	addr := common.NewAddress(registryConfig.Address)
	sc := []constant.ServerConfig{
		{
			IpAddr: addr.Host,
			Port:   uint64(addr.Port),
		},
	}

	cc := constant.ClientConfig{
		NamespaceId:         GloryNacosNamespaceID, //namespace id
		TimeoutMs:           NacosTimeoutMs,
		NotLoadCacheAtStart: true,
		LogDir:              NacosLogDir,
		CacheDir:            NacosCacheDir,
		RotateTime:          NacosRotateTime,
		MaxAge:              NacosMaxAge,
		LogLevel:            NacosLogLevel,
	}

	client, err := clients.CreateNamingClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	if err != nil {
		panic(err)
	}
	return &nacosRegistry{
		client: client,
	}
}

// Subscribe is undefined
func (nr *nacosRegistry) Subscribe(key string) (chan common.RegistryChangeEvent, error) {
	ch := make(chan common.RegistryChangeEvent)
	if err := nr.client.Subscribe(&vo.SubscribeParam{
		ServiceName: key,
		SubscribeCallback: func(services []model.SubscribeService, err error) {
			go func() {
				fmt.Println("subscribe get services = ", services)
				addrList := make([]common.Address, 0)
				for _, v := range services {
					addrList = append(addrList, *common.NewAddress(v.Ip + ":" + strconv.Itoa(int(v.Port))))
				}
				ch <- *common.NewReigstryUpdateToServiceEvent(addrList)
			}()
		},
	}); err != nil {
		log.Error("nacos registry subscribe with key = ", key, " error = ", err)
		return ch, err
	}
	return ch, nil
}

// Unsubscribe is undefined
func (nr *nacosRegistry) Unsubscribe(key string) error {
	// todo
	return nil
}

func (nr *nacosRegistry) Register(serviceID string, localAddress common.Address) {
	md := make(map[string]string)
	if ok, err := nr.client.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          localAddress.Host,
		Port:        uint64(localAddress.Port),
		ServiceName: serviceID,
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    md,
	}); err != nil || !ok {
		log.Error("nacos register with serviceID = ", serviceID, " local address = ", localAddress, " register error = ", err)
		return
	}
	log.Debug("nacos register success, with serviceID = ", serviceID, " localAddr = ", localAddress)
}

func (nr *nacosRegistry) UnRegister(serviceID string, localAddress common.Address) {
	if ok, err := nr.client.DeregisterInstance(vo.DeregisterInstanceParam{
		ServiceName: serviceID,
		Ip:          localAddress.Host,
		Port:        uint64(localAddress.Port),
		Ephemeral:   true,
	}); !ok || err != nil {
		log.Error("nacos unRegister with serviceID = ", serviceID, " local address = ", localAddress, "register error = ", err)
		return
	}
	log.Debug("nacos Unregister success, with serviceID = ", serviceID, "local address = ", localAddress)
}

func (nr *nacosRegistry) Refer(key string) []common.Address {
	service, err := nr.client.GetService(vo.GetServiceParam{
		ServiceName: key,
	})
	if err != nil {
		log.Error("nacos Refer with key = ", key, " error = ", err)
		return []common.Address{}
	}
	result := make([]common.Address, 0, 8)
	for _, instant := range service.Hosts {
		result = append(result, common.Address{
			Port: int(instant.Port),
			Host: instant.Ip,
		})
	}
	log.Debugf("nacos refer success, with addr list = %+v\n", result)
	return result
}

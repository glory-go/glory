/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package plugin

import (
	"google.golang.org/grpc/resolver"
)

import (
	"github.com/glory-go/glory/common"
	"github.com/glory-go/glory/config"
	"github.com/glory-go/glory/filter"
	"github.com/glory-go/glory/load_balance"
	"github.com/glory-go/glory/log"
	"github.com/glory-go/glory/protocol"
	"github.com/glory-go/glory/registry"
)

// protocol
type protocolFactory func(networkConfig *config.NetworkConfig, service ...interface{}) protocol.Protocol

var protocolPlugins = make(map[string]protocolFactory)

func SetProtocolFactory(protocolKey string, f protocolFactory) {
	protocolPlugins[protocolKey] = f
}

func GetProtocol(protocolKey string, networkConfig *config.NetworkConfig, service ...interface{}) protocol.Protocol {
	if f, ok := protocolPlugins[protocolKey]; ok {
		return f(networkConfig, service...)
	}
	log.Error("protocol key = plugin", protocolKey, " factory not registered!")
	return nil
}

// glory_registry
type registryFactory func(*config.RegistryConfig) registry.Registry

var registryPlugins = make(map[string]registryFactory)

func SetRegistryFactory(registryKey string, f registryFactory) {
	registryPlugins[registryKey] = f
}

func GetRegistry(registryConfig *config.RegistryConfig) registry.Registry {
	if f, ok := registryPlugins[registryConfig.Service]; ok {
		return f(registryConfig)
	}
	log.Error("glory_registry Service :", registryConfig.Service, " factory not registered!")
	return nil
}

// loadBalance
type loadBalanceFactory func() load_balance.LoadBalance

var loadBalancePlugins = make(map[string]loadBalanceFactory)

func SetLoadBalanceFactory(loadBalanceType string, f loadBalanceFactory) {
	loadBalancePlugins[loadBalanceType] = f
}

func GetLoadBalance(loadBalanceType string) load_balance.LoadBalance {
	return loadBalancePlugins[loadBalanceType]()
}

// filter
type filterFactory func(filterConfig *config.FilterConfig) (filter.GRPCFilter, error)

var filterFactoryPlugins = make(map[string]filterFactory)

func SetFilterFactory(filterKey string, f filterFactory) {
	filterFactoryPlugins[filterKey] = f
}

func GetFilter(filterKey string, filterConfig *config.FilterConfig) (filter.GRPCFilter, error) {
	return filterFactoryPlugins[filterKey](filterConfig)
}

// gRPCResolver
type gRPCResolverFactory func(ch chan common.RegistryChangeEvent, cc resolver.ClientConn) resolver.Resolver

var gRPCResolverPlugins = make(map[string]gRPCResolverFactory)

func SetGRPCResolverFactory(gRPCResolverType string, f gRPCResolverFactory) {
	gRPCResolverPlugins[gRPCResolverType] = f
}

func GetGRPCResolver(gRPCResolverType string, ch chan common.RegistryChangeEvent, cc resolver.ClientConn) resolver.Resolver {
	return gRPCResolverPlugins[gRPCResolverType](ch, cc)
}

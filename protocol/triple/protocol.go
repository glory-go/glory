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

package triple

import (
	"context"
	"reflect"

	"github.com/glory-go/glory/config"
	dubboCommon "github.com/apache/dubbo-go/common"
	dubbo3 "github.com/dubbogo/triple/pkg/triple"

	"github.com/glory-go/glory/plugin"

	"github.com/glory-go/glory/common"
	"github.com/glory-go/glory/protocol"
	dubboProto "github.com/apache/dubbo-go/protocol"
)

func init() {
	plugin.SetProtocolFactory("triple", NewTripleProtocol)
}

type Protocol struct {
	invoker         common.Invoker
	netWorkConfig   *config.NetworkConfig
	providerService dubbo3.Dubbo3GrpcService
	consumerService interface{}
}

func NewTripleProtocol(network *config.NetworkConfig, service ...interface{}) protocol.Protocol {
	var proService interface{}
	var conService interface{}
	var tripleGrpcService dubbo3.Dubbo3GrpcService
	if len(service) > 0 {
		proService = service[0]
		if len(service) > 1 {
			conService = service[1]
		}
	}
	if proService != nil {
		tripleGrpcService = proService.(dubbo3.Dubbo3GrpcService)
	}

	return &Protocol{
		netWorkConfig:   network,
		providerService: tripleGrpcService,
		consumerService: conService,
	}
}

type DubboProxyInvoker struct {
	invoker common.Invoker
}

func (di *DubboProxyInvoker) GetUrl() *dubboCommon.URL {
	return nil
}

func (di *DubboProxyInvoker) IsAvailable() bool {
	return true
}

func (di *DubboProxyInvoker) Destroy() {
}

func (di *DubboProxyInvoker) Invoke(ctx context.Context, inv dubboProto.Invocation) dubboProto.Result {
	ins := make([]interface{}, 0)
	for _, v := range inv.Arguments() {
		ins = append(ins, v)
	}
	params := &common.Params{
		MethodName: inv.MethodName(),
		Ins:        ins,
	}
	di.invoker.Invoke(ctx, params)
	return &dubboProto.RPCResult{
		Rest: params.Out,
		Err:  params.Error,
	}
}

type DubboProxyProvider struct {
}

// todo dubbo proxy invoker ctmd 太他妈脏了卧槽
// Export @invoker is glory common invoker, not dubbogo comon invoker, there needs change from Params to invocation!
func (g *Protocol) Export(ctx context.Context, invoker common.Invoker, addr common.Address) error {
	url, _ := dubboCommon.NewURL("dubbo3://127.0.0.1?param=1")
	url.Location = addr.GetUrl()
	url.Protocol = "dubbo3"
	m, ok := reflect.TypeOf(g.providerService).MethodByName("SetProxyImpl")
	if !ok {
		panic("method SetProxyImpl is necessary for triple service")
	}
	in := []reflect.Value{reflect.ValueOf(g.providerService)}
	dubboInvoker := &DubboProxyInvoker{
		invoker: invoker,
	}
	in = append(in, reflect.ValueOf(dubboInvoker))
	m.Func.Call(in)
	srv := dubbo3.NewTripleServer(url, g.providerService)
	srv.Start()
	select {}
	return nil
}

func (g *Protocol) Refer(addr common.Address) common.Invoker {
	invoker, _ := newTripleInvoker(addr.GetUrl(), g.netWorkConfig, g.consumerService)
	return invoker
}

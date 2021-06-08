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

package http

import (
	"context"

	"github.com/glory-go/glory/common"
	"github.com/glory-go/glory/config"
	"github.com/glory-go/glory/plugin"
	"github.com/glory-go/glory/protocol"
)

func init() {
	plugin.SetProtocolFactory("http", NewHTTPProtocol)
}

type Protocol struct {
	invoker       common.Invoker
	netWorkConfig *config.NetworkConfig
	//providerService dubbo3.Dubbo3GrpcService
	//consumerService interface{}
}

func NewHTTPProtocol(network *config.NetworkConfig, service ...interface{}) protocol.Protocol {
	return &Protocol{
		netWorkConfig: network,
	}
}

// Export @invoker is glory common invoker
func (g *Protocol) Export(ctx context.Context, invoker common.Invoker, addr common.Address) error {
	//httpService := service.NewHttpService()
	//httpService.RegisterRouter("/testwithfilter/{hello}/{hello2}", testHandler, &gloryHttpReq{}, &gloryHttpRsp{}, "POST", myFilter1, myFilter2)
	//httpService.Run(ctx)

	//url, _ := dubboCommon.NewURL("dubbo3://127.0.0.1?param=1")
	//url.Location = addr.GetUrl()
	//url.Protocol = "dubbo3"
	//m, ok := reflect.TypeOf(g.providerService).MethodByName("SetProxyImpl")
	//if !ok {
	//	panic("method SetProxyImpl is necessary for triple service")
	//}
	//in := []reflect.Value{reflect.ValueOf(g.providerService)}
	//dubboInvoker := &DubboProxyInvoker{
	//	invoker: invoker,
	//}
	//in = append(in, reflect.ValueOf(dubboInvoker))
	//m.Func.Call(in)
	//srv := dubbo3.NewTripleServer(url, g.providerService)
	//srv.Start()
	//select {}
	return nil
}

func (g *Protocol) Refer(addr common.Address) common.Invoker {
	//invoker, _ := newTripleInvoker(addr.GetUrl(), g.netWorkConfig, g.consumerService)
	//return invoker
	return nil
}

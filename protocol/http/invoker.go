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
	"reflect"

	"github.com/glory-go/glory/common"
	"github.com/glory-go/glory/config"
)

type Invoker struct {
	//targetAddress string
	addr *common.Address
	//tripleClient  *dubbo3.TripleClient
	//pendingMap    sync.Map
	//timeout       int
}

func newHTTPInvoker(targetAddress string, config *config.NetworkConfig, consumerService interface{}) (*Invoker, error) {
	//url, _ := dubboCommon.NewURL("dubbo3://127.0.0.1?param=1")
	//url.Location = targetAddress
	//url.Protocol = "dubbo3"
	//tripleConnClient, err := dubbo3.NewTripleClient(url, consumerService)
	//if err != nil {
	//	return nil, err
	//}
	//timeout := 3000
	//if config != nil {
	//	timeout = config.Timeout
	//}
	//newInvoker := &Invoker{
	//	addr:          common.NewAddress(targetAddress),
	//	targetAddress: targetAddress,
	//	tripleClient:  tripleConnClient,
	//	pendingMap:    sync.Map{},
	//	timeout:       timeout,
	//}
	return nil, nil
}

// StreamRecv the naming method makes me sick!
// in condition of invoker, I define this function "receive", in invoker level it's a receive event
// but in condition of global, I define this function "send", as common.StreamSendPkg bellow, in global level, it's send.
func (gi *Invoker) StreamRecv(param *common.Params) error {
	return nil
}

// StreamInvoke invoker server and start a streaming invocation
func (gi *Invoker) StreamInvoke(ctx context.Context, param *common.Params, rspChanType reflect.Type) (reflect.Value, *common.Address, uint64, error) {

	// todo: now StreamInvoke never return error from server
	return reflect.Value{}, gi.addr, param.Seq, nil
}

func (gi *Invoker) Invoke(ctx context.Context, param *common.Params) error {
	//in := make([]reflect.Value, 0, 16)
	//in = append(in, reflect.ValueOf(ctx))
	//if len(param.Ins) > 0 {
	//	for _, v := range param.Ins {
	//		in = append(in, reflect.ValueOf(v))
	//	}
	//}
	//
	//method := gi.tripleClient.Invoker.MethodByName(param.MethodName)
	//res := method.Call(in)
	//
	//if len(res) == 0 {
	//} else {
	//	if !res[1].IsNil() {
	//		param.Error = res[1].Interface().(error)
	//	}
	//	hessian2.ReflectResponse(res[0], param.Out)
	//}
	return nil
}

func (gi *Invoker) GetAddr() *common.Address {
	return gi.addr
}

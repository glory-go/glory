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

package round_robin

import (
	"context"
	"reflect"
	"sync"
	"time"

	err "github.com/glory-go/glory/error"

	"github.com/glory-go/glory/tools"
	"go.uber.org/atomic"

	"github.com/glory-go/glory/log"

	"github.com/glory-go/glory/protocol"
	"github.com/glory-go/glory/registry"

	"github.com/glory-go/glory/common"
)

type RoundRobinInvoker struct {
	invokerList     []common.Invoker
	reg             registry.Registry
	existAddrMap    map[string]bool
	serviceID       string
	protoc          protocol.Protocol
	seq             *atomic.Int64
	invokerListlock sync.RWMutex
}

func newRoundRobinInvoker(invokerList []common.Invoker, registry registry.Registry, serviceID string, protoc protocol.Protocol) *RoundRobinInvoker {
	existMap := make(map[string]bool)
	for _, v := range invokerList {
		existMap[v.GetAddr().GetUrl()] = true
	}
	newRRInvoker := &RoundRobinInvoker{
		invokerList:     invokerList,
		reg:             registry,
		serviceID:       serviceID,
		protoc:          protoc,
		existAddrMap:    existMap,
		seq:             atomic.NewInt64(0),
		invokerListlock: sync.RWMutex{},
	}
	log.Debug("RRInvoker have list = ", len(newRRInvoker.invokerList))
	if registry != nil {
		go tools.SetTimeClickFunction(time.Second*5, newRRInvoker.refresh)
	}
	return newRRInvoker
}

func (rri *RoundRobinInvoker) Invoke(ctx context.Context, in *common.Params) error {
	in.Seq = uint64(rri.seq.Inc())
	// todo divide zero error
	rri.invokerListlock.RLock()
	defer rri.invokerListlock.RUnlock()
	if len(rri.invokerList) == 0 {
		return err.GloryFrameworkErrorTargetInvokerNotFound
	}
	return rri.invokerList[in.Seq%uint64(len(rri.invokerList))].Invoke(ctx, in)
}

func (rri *RoundRobinInvoker) StreamInvoke(ctx context.Context, in *common.Params, rspType reflect.Type) (reflect.Value, *common.Address, uint64, error) {
	in.Seq = uint64(rri.seq.Inc())
	rri.invokerListlock.RLock()
	defer rri.invokerListlock.RUnlock()
	if len(rri.invokerList) == 0 {
		return reflect.Value{}, nil, 0, err.GloryFrameworkErrorTargetInvokerNotFound
	}
	return rri.invokerList[in.Seq%uint64(len(rri.invokerList))].StreamInvoke(ctx, in, rspType)
}

func (rri *RoundRobinInvoker) StreamRecv(in *common.Params) error {
	log.Debug("SteramRecv get in.Param = ", in.Addr.GetUrl())
	for i, _ := range rri.invokerList {
		if rri.invokerList[i].GetAddr().Equal(in.Addr) {
			log.Debugf("streamRecv called ", i, "th invoker")
			rri.invokerList[i].StreamRecv(in)
			return nil
		}
	}
	log.Error("StreamRecv: target invoker no found, maybe connection closed by remote")
	return err.GloryFrameworkErrorTargetInvokerNotFound
}

func (rri *RoundRobinInvoker) refresh() {
	log.Debugf("refresh was called!\n")
	addrs := rri.reg.Refer(rri.serviceID)
	newInvokerList := make([]common.Invoker, 0)
	stillExistMap := make(map[string]bool)
	rri.invokerListlock.Lock()
	defer rri.invokerListlock.Unlock()
	oldLen := len(rri.invokerList)
	for _, v := range addrs {
		if rri.existAddrMap[v.GetUrl()] {
			stillExistMap[v.GetUrl()] = true
		} else {
			newInvokerList = append(newInvokerList, rri.protoc.Refer(v))
		}
	}
	for i := oldLen - 1; i >= 0; i-- {
		if _, ok := stillExistMap[rri.invokerList[i].GetAddr().GetUrl()]; ok {
			newInvokerList = append(newInvokerList, rri.invokerList[i])
		}
	}
	rri.invokerList = newInvokerList
	log.Debugf("new invokerList = %+v\n", rri.invokerList)
}
func (rri *RoundRobinInvoker) GetAddr() *common.Address {
	return nil
}

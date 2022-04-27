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

package glory

import (
	"context"
	"reflect"
	"sync"
	"time"
)

import (
	"github.com/glory-go/glory/common"
	"github.com/glory-go/glory/config"
	"github.com/glory-go/glory/log"
)

type GloryInvoker struct {
	handler         *GloryPkgHandler
	targetAddress   string
	addr            *common.Address
	gloryConnClient *gloryConnClient
	pendingMap      sync.Map
	timeout         int
}

func newGloryInvoker(targetAddress string, config *config.NetworkConfig) (*GloryInvoker, error) {
	gloryConnClient, err := newGloryConnClient(targetAddress)
	if err != nil {
		return nil, err
	}
	timeout := 3000
	if config != nil {
		timeout = config.Timeout
	}
	newInvoker := &GloryInvoker{
		addr:            common.NewAddress(targetAddress),
		targetAddress:   targetAddress,
		gloryConnClient: gloryConnClient,
		handler:         NewGloryPkgHandler(),
		pendingMap:      sync.Map{},
		timeout:         timeout,
	}
	go newInvoker.startRspListening()
	return newInvoker, nil
}

func (gi *GloryInvoker) startRspListening() {
	for {
		buf, err := gi.gloryConnClient.ReadFrame()
		if err != nil {
			log.Error("glory conn client read frame error")
			break
		}
		rspPkg, err := gi.handler.Unmarshal(buf)
		if err != nil {
			log.Error("get rsp Pkg error:", err)
			continue
		}
		rsp, ok := rspPkg.(*ResponsePackage)
		if !ok {
			log.Error("rsp package is not ResponsePackage type")
			continue
		}
		log.Debug("rspPkg.Seq = ", rsp.Header.Seq)
		val, ok := gi.pendingMap.Load(rsp.Header.Seq)
		if !ok {
			log.Error("gi.pendingMap.Load with seq = ", rsp.Header.Seq, "err, key not exist!")
			continue
		}
		rspChannel := val.(chan interface{})
		rspChannel <- rspPkg
	}
}

// StreamRecv the naming method makes me sick!
// in condition of invoker, I define this function "receive", in invoker level it's a receive event
// but in condition of global, I define this function "send", as common.StreamSendPkg bellow, in global level, it's send.
func (gi *GloryInvoker) StreamRecv(param *common.Params) error {
	gloryPkg := newGloryRequestPackage("", param.MethodName, uint64(common.StreamSendPkg), param.Seq)
	gloryPkg.Params = append(gloryPkg.Params, param.Value)
	gloryPkg.Header.ChanOffset = param.ChanOffset
	gloryPkg.Header.Seq = param.Seq
	if err := gloryPkg.sendToConn(gi.gloryConnClient, gi.handler); err != nil {
		log.Error("StreamRecv: gloryPkg.sendToConn(gi.conn, gi.handler) err =", err)
		return GloryErrorConnErr
	}
	return nil
}

// StreamInvoke invoker server and start a streaming invocation
func (gi *GloryInvoker) StreamInvoke(ctx context.Context, param *common.Params, rspChanType reflect.Type) (reflect.Value, *common.Address, uint64, error) {
	gloryPkg := newGloryRequestPackage("", param.MethodName, uint64(common.StreamRequestPkg), param.Seq)
	gloryPkg.Params = param.Ins
	// only one rspChannel for once invoke
	rspChannel := make(chan interface{})
	gi.pendingMap.Store(param.Seq, rspChannel)
	rspChan := reflect.MakeChan(rspChanType, 0)
	if err := gloryPkg.sendToConn(gi.gloryConnClient, gi.handler); err != nil {
		log.Error("StreamInvoke:gloryPkg.sendToConn(gi.conn, gi.handler) err =", err)
		return rspChan, nil, 0, GloryErrorConnErr
	}
	timeoutCaller := time.After(time.Millisecond * time.Duration(gi.timeout))
	for {
		var rspRawPkg interface{}
		select {
		case <-timeoutCaller:
			log.Error("stream invoke timeout")
			close(rspChannel)
			gi.pendingMap.Delete(param.Seq)
			return rspChan, nil, 0, GloryErrorConnErr
		case rspRawPkg = <-rspChannel: // wait until receive StreamReady Pkg:
		}
		rspPkg, ok := rspRawPkg.(*ResponsePackage)
		if !ok {
			log.Error("StreamInvoke:rspRawPkg assert not *ResponsePackage err")
			return rspChan, nil, 0, GloryErrorProtocol
		}
		if rspPkg.Error.Code != GloryErrorNoErr.Code { // stream rpc invoke not success
			gi.pendingMap.Delete(param.Seq)
			close(rspChannel)
			return rspChan, nil, 0, rspPkg.Error
		}
		if common.PkgType(rspPkg.Header.PkgType) == common.StreamReadyPkg {
			break
		}
	}

	go func() {
		for {
			rspPkg := (<-rspChannel).(*ResponsePackage)
			if common.PkgType(rspPkg.Header.PkgType) == common.StreamRecvPkg {
				rspChan.Send(reflect.ValueOf(rspPkg.Result[0]).Elem())
			}
		}
	}()
	// todo: now StreamInvoke never return error from server
	return rspChan, gi.addr, param.Seq, nil
}

func (gi *GloryInvoker) Invoke(ctx context.Context, param *common.Params) error {
	gloryPkg := newGloryRequestPackage("", param.MethodName, uint64(common.NormalRequestPkg), param.Seq)
	gloryPkg.Params = param.Ins
	rspChannel := make(chan interface{})
	gi.pendingMap.Store(param.Seq, rspChannel)
	defer func() {
		close(rspChannel)
		gi.pendingMap.Delete(param.Seq)
	}()
	if err := gloryPkg.sendToConn(gi.gloryConnClient, gi.handler); err != nil {
		log.Error("Invoke: gloryPkg.sendToConn(gi.conn, gi.handler) err = ", err)
		return GloryErrorConnErr
	}
	timeoutcaller := time.After(time.Duration(gi.timeout) * time.Millisecond)

	var recvPkg interface{}
	select {
	case <-timeoutcaller:
		log.Error(`"invoke timeout"`)
		return GloryErrorTimeoutErr
	case recvPkg = <-rspChannel:
	}

	recv, ok := recvPkg.(*ResponsePackage) // check package
	if !ok {
		log.Error("Invoke: recvPkg assert *ResponsePackage err")
		return GloryErrorProtocol
	}

	if len(recv.Result) == 0 {
	} else {
		param.Out = recv.Result[0]
	}
	param.Error = recv.Error

	return nil
}

func (gi *GloryInvoker) GetAddr() *common.Address {
	return gi.addr
}

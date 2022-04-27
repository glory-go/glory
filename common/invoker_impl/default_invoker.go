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

package invoker_impl

import (
	"context"
	"reflect"
)

import (
	"github.com/glory-go/glory/common"
	err "github.com/glory-go/glory/error"
	"github.com/glory-go/glory/log"
)

type DefaultInvoker struct {
	fValMap  map[string]reflect.Value
	fTypeMap map[string]reflect.Type

	// streamMethod store methodName which is stream
	streamMethod map[string]bool

	// chanValMap store alive channel of stream function, first key by param.Key, second by channel offset
	chanValMap map[string][]reflect.Value
}

func newDefaultInvoker() *DefaultInvoker {
	return &DefaultInvoker{
		fValMap:      make(map[string]reflect.Value),
		fTypeMap:     make(map[string]reflect.Type),
		streamMethod: make(map[string]bool),
		chanValMap:   make(map[string][]reflect.Value),
	}
}

// StreamRecv send stream param @v, with target stream param @in.ChanOffset
func (d *DefaultInvoker) StreamRecv(in *common.Params) error {
	log.Debug("StreamRecv", in.MethodName, "offset =", in.ChanOffset, " seq = ", in.Seq, "len = ", len(d.chanValMap[in.GetUniqueCallKey()]))
	targetRecvChans, ok := d.chanValMap[in.GetUniqueCallKey()]
	if !ok {
		log.Error("stream receive with uniqueCallKey = ", in.GetUniqueCallKey())
		return err.GloryFrameworkErrorServerStreamReceiveUniqueKeyNotFound
	}
	if len(targetRecvChans) < int(in.ChanOffset)+1 {
		log.Error("stream receive with uniqueCallKey = ", in.GetUniqueCallKey(), "offset = ", in.ChanOffset, " channel not exist")
		return err.GloryFrameworkErrorServerStreamReceiveChanOffsetNotFound
	}
	targetRecvChans[in.ChanOffset].Send(reflect.ValueOf(in.Value).Elem())
	return nil
}

// StreamInvoke start call a stream function and  return @reflect.Value means the rpc rsp channel
func (d *DefaultInvoker) StreamInvoke(ctx context.Context, in *common.Params, rspType reflect.Type) (reflect.Value, *common.Address, uint64, error) {
	fVal, okVal := d.fValMap[in.MethodName]
	if !okVal {
		log.Error("server not found stream method: ", in.MethodName)
		return reflect.Value{}, nil, in.Seq, err.GloryFrameworkErrorServerStreamMethodNotFound
	}
	fType, okType := d.fTypeMap[in.MethodName]
	if !okType {
		return reflect.Value{}, nil, in.Seq, err.GloryFrameworkErrorServerStreamMethodNotFound
	}
	paramList := in.Ins

	// store full duplex channel of two or more
	valueChanList := []reflect.Value{
		reflect.ValueOf(ctx),
	}
	for _, v := range paramList {
		valueChanList = append(valueChanList, reflect.ValueOf(v))
	}

	numIN := fType.NumIn()
	for j := 0; j < numIN; j++ {
		if fType.In(j).Kind() == reflect.Chan {
			typ := fType.In(j)
			valueChanList = append(valueChanList, reflect.MakeChan(typ, 0))
		}
	}
	log.Debugf("call value list = %+v\n", valueChanList)
	go func() {
		defer func() {
			if e := recover(); e != nil {
				log.Error("stream default invoke call err = ", e)
			}
		}()
		fVal.Call(valueChanList)
	}()

	d.chanValMap[in.GetUniqueCallKey()] = valueChanList[1+len(in.Ins):]
	log.Debug("set unique call key = ", in.GetUniqueCallKey())

	//chainListValue[0].Send(
	//	reflect.ValueOf(ReqStruct{
	//		Input: "hi",
	//	}),
	//)
	//rsp, _ := chainListValue[1].Recv()
	return valueChanList[len(valueChanList)-1], nil, in.Seq, nil
}

func (d *DefaultInvoker) Invoke(ctx context.Context, in *common.Params) error {
	f, ok := d.fValMap[in.MethodName]
	if !ok {
		// not found target method
		log.Error("server not found unary method: ", in.MethodName)
		return err.GloryFrameworkErrorServerUnaryMethodNotFound
	}
	paramList := in.Ins

	valueList := []reflect.Value{
		reflect.ValueOf(ctx),
	}
	for _, v := range paramList {
		valueList = append(valueList, reflect.ValueOf(v))
	}
	outValue := f.Call(valueList)

	in.Out = outValue[0].Interface()

	if outValue[len(outValue)-1].Interface() == nil {
		return nil
	}
	return outValue[len(outValue)-1].Interface().(error)
}

func (d *DefaultInvoker) setFunc(methodName string, f reflect.Value, t reflect.Type) {
	d.fValMap[methodName] = f
	d.fTypeMap[methodName] = t
}

//NewInvokerFromProvider create an invoker from all provider's function
func NewInvokerFromProvider(provider interface{}) common.Invoker {
	invoker := newDefaultInvoker()
	val := reflect.ValueOf(provider)
	typ := reflect.TypeOf(provider)
	nums := val.NumMethod()
	for i := 0; i < nums; i++ {
		invoker.setFunc(typ.Method(i).Name, val.Method(i), typ.Method(i).Type)
		//log.Debug("type = ", typ.Method(i).Name, " val = ", val.Method(i))
	}

	return invoker
}

func (d *DefaultInvoker) GetAddr() *common.Address {
	return nil
}

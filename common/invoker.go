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

package common

import (
	"context"
	"reflect"
	"strconv"
)

type Invoker interface {
	Invoke(ctx context.Context, in *Params) error
	StreamRecv(in *Params) error
	StreamInvoke(ctx context.Context, in *Params, rspType reflect.Type) (reflect.Value, *Address, uint64, error)
	GetAddr() *Address
}

type Params struct {
	// normal
	MethodName string
	Ins        []interface{}
	Out        interface{}
	Seq        uint64
	Error      error

	// stream
	ChanOffset uint8       // send channel number offset
	Value      interface{} // used in StreamRecv, get value from client
	Addr       *Address    // used in StreamRecv, send recv pkg to old connection address's invoker
}

func (p *Params) GetUniqueCallKey() string {
	return p.MethodName + strconv.Itoa(int(p.Seq))
}

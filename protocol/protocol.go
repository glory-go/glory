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

package protocol

import (
	"context"

	"github.com/glory-go/glory/common"
)

type Protocol interface {
	Export(ctx context.Context, invoker common.Invoker, address common.Address) error // 将一个服务通过该协议暴露
	Refer(address common.Address) common.Invoker                                      // 获取一个使用该协议来调用服务端的实体：invoker_impl
}

type MiddleProtocol interface {
	Export(invoker common.Invoker, port int) (common.Invoker, error) // 将一个服务通过该协议暴露
}

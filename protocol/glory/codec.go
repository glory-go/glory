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
	hessian "github.com/apache/dubbo-go-hessian2"
)

import (
	"github.com/glory-go/glory/log"
)

func NewGloryPkgHandler() *GloryPkgHandler {
	return &GloryPkgHandler{}
}

type GloryPkgHandler struct {
}

func (g *GloryPkgHandler) Marshal(pkg interface{}) ([]byte, error) {
	encoder := hessian.NewEncoder()
	if err := encoder.Encode(pkg); err != nil {
		return nil, err
	}
	return encoder.Buffer(), nil
}

func (g *GloryPkgHandler) Unmarshal(data []byte) (interface{}, error) {
	obj, err := hessian.NewDecoder(data).Decode()
	if err != nil {
		log.Error("glory pkg handler unmarshal err:", err)
		return nil, err
	}
	return obj, nil
}

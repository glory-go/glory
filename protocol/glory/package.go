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

const (
	defaultVersion = "1.0.0"
)

// Header is glory network header field,
type Header struct {
	Version    string // Version is glory version
	TraceID    string // TraceID is used to get all link way trace usage
	PkgType    uint64 // PkgType define the package type
	Seq        uint64 // Seq is the request seq number
	ChanOffset uint8  // ChanOffset is used when stream channel send data, define which target Chan sends
}

func init() {
	hessian.RegisterPOJO(&Header{})
	hessian.RegisterPOJO(&RequestPackage{})
	hessian.RegisterPOJO(&ResponsePackage{})
}

func (Header) JavaClassName() string {
	return "GLORY_GloryHeader"
}

// RequestPackage is sent from client to server
type RequestPackage struct {
	Header     *Header       // Header is glory header
	MethodName string        // MethodName is method to invoke
	Params     []interface{} // Params is param list
}

func newGloryRequestPackage(traceID, methodName string, pkgType, seq uint64) *RequestPackage {
	return &RequestPackage{
		Header: &Header{
			TraceID: traceID,
			PkgType: pkgType,
			Seq:     seq,
			Version: defaultVersion,
		},
		MethodName: methodName,
		Params:     make([]interface{}, 0),
	}
}

func (RequestPackage) JavaClassName() string {
	return "GLORY_RequestPackage"
}

func (g *RequestPackage) sendToConn(client *gloryConnClient, handler *GloryPkgHandler) error {
	log.Debugf("SendToConn: Pkg = %+v", *g)
	data, err := handler.Marshal(g)
	if err != nil {
		log.Error("GloryPackage handler.Marshal(g) err = ", err)
		return err
	}
	if _, err := client.WriteFrame(data); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

// ResponsePackage is sent from server to client
type ResponsePackage struct {
	Header *Header       // Header is glory header
	Result []interface{} // Result is response result list
	Error  *Error        // Error defined the response status
}

func NewResponsePackage(traceID string, pkgType, seq uint64, err *Error) *ResponsePackage {
	return &ResponsePackage{
		Header: &Header{
			TraceID: traceID,
			PkgType: pkgType,
			Seq:     seq,
			Version: defaultVersion,
		},
		Result: make([]interface{}, 0),
		Error:  err,
	}
}

func (ResponsePackage) JavaClassName() string {
	return "GLORY_ResponsePackage"
}

func (g *ResponsePackage) sendToConn(client *gloryConnClient, handler *GloryPkgHandler) error {
	log.Debugf("SendToConn: Pkg = %+v", *g)
	data, err := handler.Marshal(g)
	if err != nil {
		log.Error("GloryPackage handler.Marshal(g) err = ", err)
		return err
	}
	if _, err := client.WriteFrame(data); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

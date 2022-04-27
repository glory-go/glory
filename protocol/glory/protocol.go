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
	"net"
	"reflect"
)

import (
	"github.com/glory-go/glory/common"
	"github.com/glory-go/glory/config"
	"github.com/glory-go/glory/log"
	"github.com/glory-go/glory/plugin"
	"github.com/glory-go/glory/protocol"
)

func init() {
	plugin.SetProtocolFactory("glory", NewGloryProtocol)
}

type GloryProtocol struct {
	gloryPkgHandler *GloryPkgHandler
	invoker         common.Invoker
	netWorkConfig   *config.NetworkConfig
}

// NewGloryProtocol create new glory protocol, @opt is for pb extension
func NewGloryProtocol(network *config.NetworkConfig, opt ...interface{}) protocol.Protocol {
	return &GloryProtocol{
		gloryPkgHandler: NewGloryPkgHandler(),
		netWorkConfig:   network,
	}
}

func (g *GloryProtocol) Export(ctx context.Context, invoker common.Invoker, addr common.Address) error {
	g.invoker = invoker
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr.GetUrl())
	if err != nil {
		log.Error("net.ResolveTCP Addr err = ", err)
		return err
	}

	lst, err := net.ListenTCP("tcp", tcpAddr)
	log.Debug("listening addr = ", tcpAddr)
	if err != nil {
		log.Error("net.ListenTCP error = ", err)
		return err
	}

	for {
		conn, err := lst.AcceptTCP()
		log.Debug("accept")
		if err != nil {
			log.Error("lst.AcceptTCP error = ", err)
			continue
		}
		chanByteStr := make(chan string)
		grStopChain := make(chan interface{}, 1)
		connCloseStopChain := make(chan interface{}, 1)
		gloryClient := newGloryConnClientFromConn(conn)
		// chanValue -> chan interface{}
		go func() {
			for {
				select {
				case <-grStopChain:
					log.Debug("stopChain")
					return
				default:
					buf, err := gloryClient.ReadFrame()
					if err != nil {
						connCloseStopChain <- struct{}{}
						return
					}
					chanByteStr <- string(buf)
				}
			}
		}()

		go func() {
			for {
				select {
				case <-ctx.Done():
					grStopChain <- struct{}{}
					return
				case <-connCloseStopChain:
					return

				case dataStr := <-chanByteStr:
					data := []byte(dataStr)
					pkg, err := g.gloryPkgHandler.Unmarshal(data)
					if err != nil {
						log.Error("g.gloryPkgHandler.Unmarshal Error", err)
						continue
					}
					req, ok := pkg.(*RequestPackage)
					if !ok {
						log.Error("g.gloryPkgHandler get pkg is not RequestPackage")
						continue
					}
					// checkVersion
					if req.Header.Version != defaultVersion {
						log.Error("recv unsupportex version = ", req.Header.Version, " support version = ", defaultVersion)
						g.sendErrorResponse(req, gloryClient, GloryErrorVersion)
						continue
					}
					log.Debugf("get pkg = %+v", *req)
					switch common.PkgType(req.Header.PkgType) {
					case common.NormalRequestPkg:
						go g.handleNormalRequest(gloryClient, req)
					case common.StreamRequestPkg:
						go g.handleStreamRequest(ctx, gloryClient, req)
					case common.StreamSendPkg:
						go g.handleStreamSendPkg(req)
					default:
						// unexpected Type error
						g.sendErrorResponse(req, gloryClient, GloryErrorPkgTypeError)
						log.Error("recv unsupported pkg.Header.PkgType = ", req.Header.PkgType)
					}
				}

			}
		}()

	}
}

func (g *GloryProtocol) sendErrorResponse(req *RequestPackage, client *gloryConnClient, err *Error) {
	rspPkg := NewResponsePackage(req.Header.TraceID, uint64(common.ErrorRspPkg), req.Header.Seq, err)
	_ = rspPkg.sendToConn(client, g.gloryPkgHandler)
}

func (g *GloryProtocol) handleStreamSendPkg(pkg *RequestPackage) {
	callStreamSendPkg(g.invoker, pkg)
}

func (g *GloryProtocol) handleStreamRequest(ctx context.Context, client *gloryConnClient, pkg *RequestPackage) {

	rspChan, err := callStreamInvoker(ctx, g.invoker, pkg)
	if err != nil {
		log.Error("callStreamInvoker err = ", err)
		g.sendErrorResponse(pkg, client, err)
		return
	}
	rspPkg := NewResponsePackage(pkg.Header.TraceID, uint64(common.StreamReadyPkg), pkg.Header.Seq, GloryErrorNoErr)
	if err := rspPkg.sendToConn(client, g.gloryPkgHandler); err != nil {
		log.Error("handleStreamRequest:rspPkg.sendToConn(conn, g.gloryPkgHandler) err =", err)
		return
	}
	closeChan := make(chan interface{})
	valueChan := make(chan reflect.Value)
	go func() {
		for {
			select {
			case <-closeChan:
				return
			default:
				rsp, ok := rspChan.Recv()
				if !ok {
					log.Error("rspChan.Recv() not ok ")
					return
				}
				valueChan <- rsp
			}
		}
	}()
	for {
		select {
		case <-ctx.Done():
			closeChan <- struct{}{}
			return
		case rsp := <-valueChan:
			rspPkg := NewResponsePackage(pkg.Header.TraceID, uint64(common.StreamRecvPkg), pkg.Header.Seq, GloryErrorNoErr)
			rspPkg.Result = append(rspPkg.Result, rsp.Interface())
			if err := rspPkg.sendToConn(client, g.gloryPkgHandler); err != nil {
				log.Error("handleStreamRequest:rspPkg.sendToConn(conn, g.gloryPkgHandler) err =", err)
				return
			}
		}

	}
}

func (g *GloryProtocol) handleNormalRequest(client *gloryConnClient, pkg *RequestPackage) {
	rsp, err := callNormalInvoker(g.invoker, pkg)
	// err can only be user error or no error, as next level is user level(framework level), not protocol level.
	rspPkg := NewResponsePackage(pkg.Header.TraceID, uint64(common.NormalResponsePkg), pkg.Header.Seq, err)

	rspList, ok := rsp.([]interface{})
	if !ok || len(rspList) > 0 { // call successfully
		rspPkg.Result = append(rspPkg.Result, rsp)
	}

	if err := rspPkg.sendToConn(client, g.gloryPkgHandler); err != nil {
		log.Error("handleNormalRequest: rspPkg.sendToConn(conn, g.gloryPkgHandler) err:", err)
	}
}

func callNormalInvoker(invoker common.Invoker, pkg *RequestPackage) (interface{}, *Error) {
	params := &common.Params{
		MethodName: pkg.MethodName,
		Ins:        pkg.Params,
		Seq:        pkg.Header.Seq,
		Out:        make([]interface{}, 0),
	}
	err := invoker.Invoke(context.Background(), params) // framework(user) level error! not protocol level error
	if err != nil {
		return params.Out, NewGloryErrorUserError(err)
	}
	return params.Out, GloryErrorNoErr
}

func callStreamInvoker(ctx context.Context, invoker common.Invoker, pkg *RequestPackage) (reflect.Value, *Error) {
	params := &common.Params{
		MethodName: pkg.MethodName,
		Seq:        pkg.Header.Seq,
		Ins:        pkg.Params,
	}
	rspChan, _, _, err := invoker.StreamInvoke(ctx, params, nil)
	if err != nil {
		log.Error("invoker.StreamInvoke err = ", err)
		return reflect.Value{}, NewGloryErrorUserError(err)
	}
	return rspChan, nil
}

func callStreamSendPkg(invoker common.Invoker, pkg *RequestPackage) {
	params := &common.Params{
		MethodName: pkg.MethodName,
		ChanOffset: pkg.Header.ChanOffset,
		Value:      pkg.Params[0],
		Seq:        pkg.Header.Seq,
	}
	_ = invoker.StreamRecv(params)
}

func (g *GloryProtocol) Refer(addr common.Address) common.Invoker {
	invoker, _ := newGloryInvoker(addr.GetUrl(), g.netWorkConfig)
	return invoker
}

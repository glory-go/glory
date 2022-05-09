package interceptor

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"
)

import (
	"google.golang.org/protobuf/types/known/emptypb"
)

import (
	"github.com/glory-go/glory/autowire/util"
	"github.com/glory-go/glory/debug/api/glory/boot"
	"github.com/glory-go/glory/debug/common"
)

var sendRecvChWatchEditMap = make(map[string]sendRecvCh)

type DebugServerImpl struct {
	editInterceptor         *EditInterceptor
	watchInterceptor        *WatchInterceptor
	allInterfaceMetadataMap map[string]*common.DebugMetadata
	boot.UnimplementedDebugServiceServer
}

func (d *DebugServerImpl) ListServices(ctx context.Context, empty *emptypb.Empty) (*boot.ListServiceResponse, error) {
	structsMetadatas := make(MetadataSorter, 0)
	for key, v := range d.allInterfaceMetadataMap {
		methods := make([]string, 0)
		for key := range v.GuardMap {
			methods = append(methods, key)
		}

		pair := strings.Split(key, "-")
		interfaceName := pair[0]
		implName := pair[1]
		structsMetadatas = append(structsMetadatas, &boot.ServiceMetadata{
			Methods:            methods,
			InterfaceName:      interfaceName,
			ImplementationName: implName,
		})
	}
	sort.Sort(structsMetadatas)

	return &boot.ListServiceResponse{
		ServiceMetadata: structsMetadatas,
	}, nil
}

func (d *DebugServerImpl) Watch(req *boot.WatchRequest, watchSever boot.DebugService_WatchServer) error {
	interfaceImplId := util.GetIdByNamePair(req.GetInterfaceName(), req.GetImplementationName())
	method := req.GetMethod()
	input := req.GetInput()
	output := req.GetOutput()
	sendCh := make(chan *boot.WatchResponse)
	fmt.Printf("interceptor server recv watch %+v\n", req)
	fmt.Println(interfaceImplId)
	fmt.Println(method)
	fmt.Println(input)
	fmt.Println(output)
	var fieldMatcher *FieldMatcher
	for _, matcher := range req.GetMatchers() {
		// todo multi match support
		fieldMatcher = &FieldMatcher{
			FieldIndex: int(matcher.Index),
			MatchRule:  matcher.GetMatchPath() + "=" + matcher.GetMatchValue(),
		}
	}
	if input {
		d.watchInterceptor.Watch(interfaceImplId, method, true, &WatchContext{
			Ch:           sendCh,
			FieldMatcher: fieldMatcher,
		})
	}

	if output {
		d.watchInterceptor.Watch(interfaceImplId, method, false, &WatchContext{
			Ch:           sendCh,
			FieldMatcher: fieldMatcher,
		})
	}

	done := watchSever.Context().Done()
	for {
		select {
		case <-done:
			// watch stop
			if input {
				d.watchInterceptor.UnWatch(interfaceImplId, method, true)
			}
			if output {
				d.watchInterceptor.UnWatch(interfaceImplId, method, false)
			}

			return nil
		case watchRsp := <-sendCh:
			if err := watchSever.Send(watchRsp); err != nil {
				return err
			}
		}
	}
}

type sendRecvCh struct {
	sendCh chan *boot.WatchResponse
	recvCh chan *EditData
}

func (d *DebugServerImpl) WatchEdit(watchEditServerReq boot.DebugService_WatchEditServer) error {
	interfaceImplId := ""
	method := ""
	isParam := false
	for {
		req, err := watchEditServerReq.Recv()
		if err != nil {
			d.watchInterceptor.UnWatch(interfaceImplId, method, isParam)
			return err
		}
		interfaceImplId = util.GetIdByNamePair(req.GetInterfaceName(), req.GetImplementationName())
		method = req.GetMethod()
		isParam = req.GetIsParam()
		uniqueMethodKey := getMethodUniqueKey(interfaceImplId, method, isParam)
		if !req.IsEdit {
			// start new watch
			_, ok := sendRecvChWatchEditMap[uniqueMethodKey]
			if ok {
				// if already watch, unwatch
				d.editInterceptor.UnWatchEdit(interfaceImplId, method, isParam)
			}
			var fieldMatcher *FieldMatcher
			sendCh := make(chan *boot.WatchResponse)
			recvCh := make(chan *EditData)
			for _, matcher := range req.GetMatchers() {
				// todo multi match support
				fieldMatcher = &FieldMatcher{
					FieldIndex: int(matcher.Index),
					MatchRule:  matcher.GetMatchPath() + "=" + matcher.GetMatchValue(),
				}
			}
			d.editInterceptor.WatchEdit(
				interfaceImplId, method, isParam,
				&EditContext{
					RecvCh:       recvCh,
					SendCh:       sendCh,
					FieldMatcher: fieldMatcher,
				})
			// start send gr
			go func() {
				toShowData := <-sendCh
				if err := watchEditServerReq.Send(toShowData); err != nil {
					log.Printf("send error = %s\n", err)
					return
				}
			}()
			sendRecvChWatchEditMap[uniqueMethodKey] = sendRecvCh{
				sendCh: sendCh,
				recvCh: recvCh,
			}
		} else {
			// edit
			oldSendRecvCh, ok := sendRecvChWatchEditMap[uniqueMethodKey]
			if !ok {
				log.Printf("uniqueMethodKey = %s old subscription shou be exist.\n", uniqueMethodKey)
				continue
			}
			if len(req.EditRequests) == 0 {
				continue
			}
			// todo support multi edit
			oldSendRecvCh.recvCh <- &EditData{
				FieldIndex: int(req.EditRequests[0].Index),
				FieldPath:  req.EditRequests[0].Path,
				Value:      req.EditRequests[0].Value,
			}
		}
	}
}

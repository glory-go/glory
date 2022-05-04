package interceptor

import (
	"context"
	"fmt"
)

import (
	"google.golang.org/protobuf/types/known/emptypb"
)

import (
	"github.com/glory-go/glory/boot/api/glory/boot"
	"github.com/glory-go/glory/boot/common"
	"github.com/glory-go/glory/boot/util"
	"github.com/glory-go/glory/log"
)

var sendRecvChWatchEditMap = make(map[string]sendRecvCh)

type DebugServerImpl struct {
	editInterceptor         *EditInterceptor
	watchInterceptor        *WatchInterceptor
	allInterfaceMetadataMap map[string]common.RegisterServiceMetadata
	boot.UnimplementedDebugServiceServer
}

func (d *DebugServerImpl) ListServices(ctx context.Context, empty *emptypb.Empty) (*boot.ListServiceResponse, error) {
	serviceMetadata := make([]*boot.ServiceMetadata, 0)
	for _, v := range d.allInterfaceMetadataMap {
		methods := make([]string, 0)
		for key, _ := range v.GuardMap {
			methods = append(methods, key)
		}
		serviceMetadata = append(serviceMetadata, &boot.ServiceMetadata{
			Methods:            methods,
			InterfaceName:      util.GetName(v.InterfaceStruct),
			ImplementationName: util.GetName(v.SvcStructPtr),
		})
	}
	return &boot.ListServiceResponse{
		ServiceMetadata: serviceMetadata,
	}, nil
}

func (d *DebugServerImpl) Watch(req *boot.WatchRequest, watchSever boot.DebugService_WatchServer) error {
	interfaceImplId := util.GetInterfaceIdByNames(req.GetInterfaceName(), req.GetImplementationName())
	method := req.GetMethod()
	isParam := req.GetIsParam()
	sendCh := make(chan string)
	fmt.Printf("interceptor server recv watch %+v\n", req)
	fmt.Println(interfaceImplId)
	fmt.Println(method)
	fmt.Println(isParam)
	var fieldMatcher *FieldMatcher
	for _, matcher := range req.GetMatchers() {
		// todo multi match support
		fieldMatcher = &FieldMatcher{
			FieldIndex: int(matcher.Index),
			MatchRule:  matcher.GetMatchPath() + "=" + matcher.GetMatchValue(),
		}
	}
	d.watchInterceptor.Watch(interfaceImplId, method, isParam, &WatchContext{
		Ch:           sendCh,
		FieldMatcher: fieldMatcher,
	})

	done := watchSever.Context().Done()
	for {
		select {
		case <-done:
			// watch stop
			d.watchInterceptor.UnWatch(interfaceImplId, method, isParam)
			return nil
		case data := <-sendCh:
			if err := watchSever.Send(&boot.WatchResponse{
				Content: data,
			}); err != nil {
				return err
			}
		}
	}
}

type sendRecvCh struct {
	sendCh chan string
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
		interfaceImplId = util.GetInterfaceIdByNames(req.GetInterfaceName(), req.GetImplementationName())
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
			sendCh := make(chan string)
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
				if err := watchEditServerReq.Send(&boot.WatchResponse{
					Content: toShowData,
				}); err != nil {
					log.Errorf("send error = %s", err)
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
				log.Errorf("uniqueMethodKey = %s old subscription shou be exist.", uniqueMethodKey)
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

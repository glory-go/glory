package interceptor

import (
	"context"
	"reflect"
	"strings"
	"testing"
	"time"
)

import (
	"github.com/stretchr/testify/assert"
)

import (
	"github.com/glory-go/glory/debug/api/glory/boot"
)

func TestEditInterceptorWithCondition(t *testing.T) {
	editInterceptor := GetEditInterceptor()
	interfaceImplId := "Service-ServiceFoo"
	methodName := "Invoke"
	sendCh := make(chan *boot.WatchResponse, 10)
	recvCh := make(chan *EditData, 10)
	controlSendCh := make(chan *boot.WatchResponse, 10)
	controlRecvCh := make(chan *EditData, 10)
	go func() {
		for {
			info := <-sendCh
			controlSendCh <- info
		}
	}()
	go func() {
		for {
			info := <-controlRecvCh
			recvCh <- info
		}
	}()
	editInterceptor.WatchEdit(interfaceImplId, methodName, true, &EditContext{
		SendCh: sendCh,
		RecvCh: recvCh,
		FieldMatcher: &FieldMatcher{
			FieldIndex: 2,
			MatchRule:  "User.Name=lizhixin",
		},
	})

	service := &ServiceFoo{}
	ctx := context.Background()

	param := &RequestParam{
		User: &User{
			Name: "lizhixin",
		},
	}
	go func() {
		controlRecvCh <- &EditData{
			FieldIndex: 2,
			FieldPath:  "User.Name",
			Value:      "laurence",
		}
	}()

	editInterceptor.Invoke(interfaceImplId, methodName, true,
		[]reflect.Value{reflect.ValueOf(service), reflect.ValueOf(ctx), reflect.ValueOf(param)})

	rsp, err := service.Invoke(ctx, param)

	time.Sleep(time.Millisecond * 500)
	info := &boot.WatchResponse{}
	select {
	case info = <-controlSendCh:
	default:
	}
	assert.Equal(t, "Service", info.InterfaceName)
	assert.Equal(t, "ServiceFoo", info.ImplementationName)
	assert.Equal(t, "Invoke", info.MethodName)
	assert.Equal(t, true, info.IsParam)
	assert.True(t, strings.Contains(info.Params[1], "lizhixin"))

	assert.Nil(t, err)
	assert.Equal(t, "laurence", rsp.Name)
}

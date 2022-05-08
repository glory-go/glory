package interceptor

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"
)

import (
	"github.com/stretchr/testify/assert"
)

func TestEditInterceptorWithCondition(t *testing.T) {
	editInterceptor := GetEditInterceptor()
	interfaceImplId := "Service-ServiceFoo"
	methodName := "Invoke"
	sendCh := make(chan string)
	recvCh := make(chan *EditData)
	controlSendCh := make(chan string)
	controlRecvCh := make(chan *EditData)
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

	// match
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

	time.Sleep(time.Second * 3)
	info := ""
	select {
	case info = <-controlSendCh:
	default:
	}
	fmt.Println(info)
	assert.True(t, strings.Contains(info, "lizhixin"))
	assert.True(t, strings.HasPrefix(info, "Invoke Service-ServiceFoo.Invoke"))

	assert.Nil(t, err)
	fmt.Println(rsp)
	assert.Equal(t, "laurence", rsp.Name)
}

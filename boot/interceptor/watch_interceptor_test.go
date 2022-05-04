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

type User struct {
	Name string
}

type RequestParam struct {
	User *User
}

type Response struct {
	Name string
}

type ServiceFoo struct {
}

func (s *ServiceFoo) Invoke(ctx context.Context, param *RequestParam) (*Response, error) {
	return &Response{
		Name: param.User.Name,
	}, nil
}

func TestWatchInterceptor(t *testing.T) {
	watchInterceptor := GetWatchInterceptor()
	interfaceImplId := "Service-ServiceFoo"
	methodName := "Invoke"
	sendCh := make(chan string)
	controlCh := make(chan string)
	go func() {
		info := <-sendCh
		controlCh <- info
		info = <-sendCh
		controlCh <- info
	}()
	watchInterceptor.Watch(interfaceImplId, methodName, true, &WatchContext{
		Ch: sendCh,
	})
	watchInterceptor.Watch(interfaceImplId, methodName, false, &WatchContext{
		Ch: sendCh,
	})

	service := &ServiceFoo{}
	ctx := context.Background()
	param := &RequestParam{
		User: &User{
			Name: "laurence",
		},
	}

	watchInterceptor.Invoke(interfaceImplId, methodName, true,
		[]reflect.Value{reflect.ValueOf(service), reflect.ValueOf(ctx), reflect.ValueOf(param)})
	rsp, err := service.Invoke(ctx, param)
	info := <-controlCh
	fmt.Println(info)
	assert.True(t, strings.HasPrefix(info, "Invoke Service-ServiceFoo.Invoke"))
	watchInterceptor.Invoke(interfaceImplId, methodName, false,
		[]reflect.Value{reflect.ValueOf(service), reflect.ValueOf(rsp), reflect.ValueOf(err)})
	info = <-controlCh
	fmt.Println(info)
	assert.True(t, strings.HasPrefix(info, "After Invoke Service-ServiceFoo.Invoke"))
}

func TestWatchInterceptorWithCondition(t *testing.T) {
	watchInterceptor := GetWatchInterceptor()
	interfaceImplId := "Service-ServiceFoo"
	methodName := "Invoke"
	sendCh := make(chan string)
	controlCh := make(chan string)
	go func() {
		for {
			info := <-sendCh
			controlCh <- info
		}
	}()
	watchInterceptor.Watch(interfaceImplId, methodName, true, &WatchContext{
		Ch: sendCh,
		FieldMatcher: &FieldMatcher{
			FieldIndex: 2,
			MatchRule:  "User.Name=lizhixin",
		},
	})

	service := &ServiceFoo{}
	ctx := context.Background()

	// not match
	param := &RequestParam{
		User: &User{
			Name: "laurence",
		},
	}
	watchInterceptor.Invoke(interfaceImplId, methodName, true,
		[]reflect.Value{reflect.ValueOf(service), reflect.ValueOf(ctx), reflect.ValueOf(param)})
	rsp, err := service.Invoke(ctx, param)
	info := ""
	time.Sleep(time.Millisecond * 500)
	select {
	case info = <-controlCh:
	default:
	}
	assert.Equal(t, "", info)
	watchInterceptor.Invoke(interfaceImplId, methodName, false,
		[]reflect.Value{reflect.ValueOf(service), reflect.ValueOf(rsp), reflect.ValueOf(err)})
	time.Sleep(time.Millisecond * 500)
	select {
	case info = <-controlCh:
	default:
	}
	assert.Equal(t, "", info)

	// match
	param.User.Name = "lizhixin"
	watchInterceptor.Invoke(interfaceImplId, methodName, true,
		[]reflect.Value{reflect.ValueOf(service), reflect.ValueOf(ctx), reflect.ValueOf(param)})
	rsp, err = service.Invoke(ctx, param)
	time.Sleep(time.Millisecond * 500)
	select {
	case info = <-controlCh:
	default:
	}
	assert.True(t, strings.HasPrefix(info, "Invoke Service-ServiceFoo.Invoke"))
	watchInterceptor.Invoke(interfaceImplId, methodName, false,
		[]reflect.Value{reflect.ValueOf(service), reflect.ValueOf(rsp), reflect.ValueOf(err)})
	time.Sleep(time.Millisecond * 500)
	info = ""
	select {
	case info = <-controlCh:
	default:
	}
	assert.Equal(t, "", info)

	// not match
	param.User.Name = "lizhixin"
	watchInterceptor.UnWatch(interfaceImplId, methodName, true)
	watchInterceptor.Invoke(interfaceImplId, methodName, true,
		[]reflect.Value{reflect.ValueOf(service), reflect.ValueOf(ctx), reflect.ValueOf(param)})
	rsp, err = service.Invoke(ctx, param)
	time.Sleep(time.Millisecond * 500)
	info = ""
	select {
	case info = <-controlCh:
	default:
	}
	assert.Equal(t, "", info)
}

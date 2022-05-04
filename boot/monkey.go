package boot

import (
	"fmt"
	"reflect"
)

import (
	"github.com/glory-go/monkey"
)

import (
	"github.com/glory-go/glory/boot/interceptor"
)

var paramInterceptors = make([]interceptor.Interceptor, 0)
var responseInterceptors = make([]interceptor.Interceptor, 0)

var debugConfig map[string]string

func init() {
	paramInterceptors = append(paramInterceptors, interceptor.GetWatchInterceptor())
	paramInterceptors = append(paramInterceptors, interceptor.GetEditInterceptor())

	responseInterceptors = append(responseInterceptors, interceptor.GetWatchInterceptor())
	responseInterceptors = append(responseInterceptors, interceptor.GetEditInterceptor())
}

func implMonkey(servicePtr interface{}, tempInterfaceId string) {
	valueOf := reflect.ValueOf(servicePtr)
	typeOf := reflect.TypeOf(servicePtr)
	valueOfElem := valueOf.Elem()
	typeOfElem := valueOfElem.Type()
	if typeOfElem.Kind() != reflect.Struct {
		panic("invalid struct ptr")
	}

	numField := valueOf.NumMethod()
	for i := 0; i < numField; i++ {
		methodType := typeOf.Method(i)
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(servicePtr), methodType.Name,
			reflect.MakeFunc(methodType.Type, makeCallProxy(tempInterfaceId, methodType.Name)).Interface(),
		)
		registeredMap[tempInterfaceId].GuardMap[methodType.Name] = guard
		continue
	}
}

func makeCallProxy(tempInterfaceId, methodName string) func(in []reflect.Value) []reflect.Value {
	return func(in []reflect.Value) []reflect.Value {
		registeredMap[tempInterfaceId].GuardMap[methodName].Unpatch()
		defer func() {
			registeredMap[tempInterfaceId].GuardMap[methodName].Unpatch()
			implMonkey(in[0].Interface(), tempInterfaceId)
		}()
		// interceptor
		fmt.Println("call proxy", tempInterfaceId, methodName)
		for _, i := range paramInterceptors {
			in = i.Invoke(tempInterfaceId, methodName, true, in)
		}
		out := in[0].MethodByName(methodName).Call(in[1:])
		for _, i := range responseInterceptors {
			out = i.Invoke(tempInterfaceId, methodName, false, out)
		}
		return out
	}
}

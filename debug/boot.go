package debug

import (
	"fmt"
	"log"
	"reflect"
)

import (
	"github.com/glory-go/monkey"
)

import (
	"github.com/glory-go/glory/autowire"
	"github.com/glory-go/glory/config"
	"github.com/glory-go/glory/debug/common"
	"github.com/glory-go/glory/debug/interceptor"
)

var paramInterceptors = make([]interceptor.Interceptor, 0)
var responseInterceptors = make([]interceptor.Interceptor, 0)

func init() {
	paramInterceptors = append(paramInterceptors, interceptor.GetWatchInterceptor())
	paramInterceptors = append(paramInterceptors, interceptor.GetEditInterceptor())

	responseInterceptors = append(responseInterceptors, interceptor.GetWatchInterceptor())
	responseInterceptors = append(responseInterceptors, interceptor.GetEditInterceptor())

	autowire.RegisterMonkeyFunction(implMonkey)
}

var guardMap = make(map[string]*common.DebugMetadata)

func Load() error {
	// start debug port if enabled
	bootConfig := &Config{}
	if err := config.LoadConfigByPrefix("debug", bootConfig); err == nil && !bootConfig.Enable {
		return nil
	}
	if bootConfig.Port == "" {
		bootConfig.Port = "1999"
	}
	log.Println("glory boot debug port start at :" + bootConfig.Port)
	if err := interceptor.Start(bootConfig.Port, guardMap); err != nil {
		return err
	}
	return nil
}

func implMonkey(servicePtr interface{}, tempInterfaceId string) {
	if _, ok := guardMap[tempInterfaceId]; !ok {
		guardMap[tempInterfaceId] = &common.DebugMetadata{
			GuardMap: map[string]*common.GuardInfo{},
		}
	}
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
		if _, ok := guardMap[tempInterfaceId].GuardMap[methodType.Name]; !ok {
			guardMap[tempInterfaceId].GuardMap[methodType.Name] = &common.GuardInfo{}
		}
		if guardMap[tempInterfaceId].GuardMap[methodType.Name].Guard == nil {
			// each method of one type should only injected once
			guard := monkey.PatchInstanceMethod(reflect.TypeOf(servicePtr), methodType.Name,
				reflect.MakeFunc(methodType.Type, makeCallProxy(tempInterfaceId, methodType.Name)).Interface(),
			)
			guardMap[tempInterfaceId].GuardMap[methodType.Name].Guard = guard
		}
		continue
	}
}

func makeCallProxy(tempInterfaceId, methodName string) func(in []reflect.Value) []reflect.Value {
	return func(in []reflect.Value) []reflect.Value {
		guardMap[tempInterfaceId].GuardMap[methodName].Lock.Lock()
		guardMap[tempInterfaceId].GuardMap[methodName].Guard.Unpatch()
		defer func() {
			guardMap[tempInterfaceId].GuardMap[methodName].Guard.Restore()
			guardMap[tempInterfaceId].GuardMap[methodName].Lock.Unlock()
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

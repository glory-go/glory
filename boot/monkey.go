package boot

import (
	"fmt"
	"reflect"
)

import (
	"github.com/glory-go/monkey"
)

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
		registeredMap[tempInterfaceId].guardMap[methodType.Name] = guard
		continue
	}
}

func makeCallProxy(tempInterfaceId, methodName string) func(in []reflect.Value) []reflect.Value {
	return func(in []reflect.Value) []reflect.Value {
		registeredMap[tempInterfaceId].guardMap[methodName].Unpatch()
		defer registeredMap[tempInterfaceId].guardMap[methodName].Unpatch()
		// interceptor
		fmt.Printf("call inject func = %+v\n", in)
		for _, v := range in {
			fmt.Printf("param = %+v\n", v.Interface())
		}
		return in[0].MethodByName(methodName).Call(in[1:])
	}
}

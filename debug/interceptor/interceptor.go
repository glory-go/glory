package interceptor

import (
	"reflect"
)

type Interceptor interface {
	Invoke(interfaceImplId, methodName string, isParam bool, value []reflect.Value) []reflect.Value
}

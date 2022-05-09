package interceptor

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"sync"
)

import (
	"github.com/davecgh/go-spew/spew"
)

import (
	"github.com/glory-go/glory/debug/api/glory/boot"
)

type WatchInterceptor struct {
	watch sync.Map
}

func (w *WatchInterceptor) Invoke(interfaceImplId, methodName string, isParam bool, values []reflect.Value) []reflect.Value {
	methodUniqueKey := getMethodUniqueKey(interfaceImplId, methodName, isParam)
	watchCtxInterface, ok := w.watch.Load(methodUniqueKey)
	if !ok {
		return values
	}
	watchCtx := watchCtxInterface.(*WatchContext)
	if watchCtx.FieldMatcher != nil && !watchCtx.FieldMatcher.Match(values) {
		// doesn't match
		return values
	}
	sendValues(interfaceImplId, methodName, isParam, values, watchCtx.Ch)
	return values
}

func sendValues(interfaceImplId, methodName string, isParam bool, values []reflect.Value, sendCh chan *boot.WatchResponse) {
	splitedSDID := strings.Split(interfaceImplId, "-")
	invokeDetail := &boot.WatchResponse{
		IsParam:            isParam,
		InterfaceName:      splitedSDID[0],
		MethodName:         methodName,
		ImplementationName: splitedSDID[1],
	}
	i := 0
	if isParam {
		// param first value is struct ptr, should skip it.
		i = 1
	}
	for ; i < len(values); i++ {
		if !values[i].IsValid() {
			invokeDetail.Params = append(invokeDetail.Params, "nil")
			continue
		}
		invokeDetail.Params = append(invokeDetail.Params, spew.Sdump(values[i].Interface()))
	}
	select {
	case sendCh <- invokeDetail:
	default:
	}
}

type WatchContext struct {
	Ch           chan *boot.WatchResponse
	FieldMatcher *FieldMatcher
}

type FieldMatcher struct {
	FieldIndex int
	MatchRule  string // A.B.C=xxx
}

func (f *FieldMatcher) Match(values []reflect.Value) bool {
	if len(values) < f.FieldIndex {
		return false
	}
	targetVal := values[f.FieldIndex]
	data, err := json.Marshal(targetVal.Interface())
	if err != nil {
		return false
	}
	anyTypeMap := make(map[string]interface{})
	if err := json.Unmarshal(data, &anyTypeMap); err != nil {
		return false
	}
	rules := strings.Split(f.MatchRule, "=")
	paths := rules[0]
	expectedValue := rules[1]
	splitedPaths := strings.Split(paths, ".")
	for i, p := range splitedPaths {
		subInterface, ok := anyTypeMap[p]
		if !ok {
			return false
		}
		if i == len(splitedPaths)-1 {
			// final must be string
			realStr, ok := subInterface.(string)
			if !ok {
				return false
			}
			if realStr != expectedValue {
				return false
			}
		} else {
			// not final, subInterface should be map[string]interface{}
			anyTypeMap, ok = subInterface.(map[string]interface{})
			if !ok {
				return false
			}
		}
	}
	return true
}

func (w *WatchInterceptor) Watch(interfaceImplId, methodName string, isParam bool, watchCtx *WatchContext) {
	methodUniqueKey := getMethodUniqueKey(interfaceImplId, methodName, isParam)
	w.watch.Store(methodUniqueKey, watchCtx)
}

func (w *WatchInterceptor) UnWatch(interfaceImplId, methodName string, isParam bool) {
	methodUniqueKey := getMethodUniqueKey(interfaceImplId, methodName, isParam)
	w.watch.Delete(methodUniqueKey)
}

var watchInterceptorSingleton *WatchInterceptor

func GetWatchInterceptor() *WatchInterceptor {
	if watchInterceptorSingleton == nil {
		watchInterceptorSingleton = &WatchInterceptor{}
	}
	return watchInterceptorSingleton
}

func getMethodUniqueKey(interfaceImplId, methodName string, isParam bool) string {
	return strings.Join([]string{interfaceImplId, methodName, fmt.Sprintf("%t", isParam)}, "-")
}

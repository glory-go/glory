package util

import (
	"reflect"
	"strings"
)

func GetInterfaceIdByNames(interfaceName, structPtrName string) string {
	return strings.Join([]string{interfaceName, structPtrName}, "-")
}

func GetName(v interface{}) string {
	if v == nil {
		return ""
	}
	typeOfInterface := GetTypeFromInterface(v)
	return typeOfInterface.Name()
}

func GetTypeFromInterface(v interface{}) reflect.Type {
	valueOfInterface := reflect.ValueOf(v)
	valueOfElemInterface := valueOfInterface.Elem()
	return valueOfElemInterface.Type()
}

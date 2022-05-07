package util

import (
	"reflect"
	"strings"
)

func GetIdByInterfaceAndImplPtr(interfaceStruct, implStructPtr interface{}) string {
	interfaceName := GetStructName(interfaceStruct)
	structPtrName := GetStructName(implStructPtr)
	return GetIdByNamePair(interfaceName, structPtrName)
}

func GetIdByNamePair(interfaceName, structPtrName string) string {
	return strings.Join([]string{interfaceName, structPtrName}, "-")
}

func GetStructName(v interface{}) string {
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

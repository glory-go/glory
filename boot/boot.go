package boot

import (
	"reflect"
	"strings"
)

import (
	perrors "github.com/pkg/errors"
)

type RegisterServicePair struct {
	interfaceStruct    interface{}
	svcStructPtr       interface{}
	constructFunctions []func(interface{})
}

var registeredMap = make(map[string]RegisterServicePair)
var implCompletedMap = make(map[string]interface{})
var controllerMap = make(map[string]interface{})

func RegisterController(controller interface{}) {
	controllerMap[getName(controller)] = controller
}

func RegisterService(interfaceStruct, structPtr interface{}, constructFunction ...func(interface{})) {
	newPair := RegisterServicePair{
		interfaceStruct:    interfaceStruct,
		svcStructPtr:       structPtr,
		constructFunctions: constructFunction,
	}

	serviceId := getInterfaceId(newPair)
	registeredMap[serviceId] = newPair
}

func Load() {
	// set impl
	for serviceId, v := range registeredMap {
		if _, ok := implCompletedMap[serviceId]; ok {
			continue
		}
		impl(v)
	}

	// impl controller
	for _, v := range controllerMap {
		implController(v)
	}
}

func implController(c interface{}) interface{} {
	valueOf := reflect.ValueOf(c)
	valueOfElem := valueOf.Elem()
	typeOf := valueOfElem.Type()
	if typeOf.Kind() != reflect.Struct {
		panic("invalid struct ptr")
	}

	numField := valueOfElem.NumField()
	for i := 0; i < numField; i++ {
		t := typeOf.Field(i)
		svcImplStructName := t.Tag.Get("service")
		if svcImplStructName == "" {
			// not autowire field
			continue
		}

		// get impled sub service
		impledPtr := impl(registeredMap[getInterfaceIdByNames(t.Type.Name(), svcImplStructName)])

		// set field
		subService := valueOfElem.Field(i)
		if !(subService.Kind() == reflect.Interface && subService.IsValid() && subService.CanSet()) {
			err := perrors.Errorf("Failed to autowire controller %s 's. It's field %s with tag 'service:\"%s\"', please check if the field is exported",
				getName(c), t.Type.Name(), svcImplStructName)
			panic(err)
		}
		subService.Set(reflect.ValueOf(impledPtr))
	}
	return c
}

func impl(p RegisterServicePair) interface{} {
	tempInterfaceId := getInterfaceId(p)
	if impledPtr, ok := implCompletedMap[tempInterfaceId]; ok {
		// if already impleted, return
		return impledPtr
	}
	defer func() { // assure the impl procedure of one service run once
		implCompletedMap[tempInterfaceId] = p.svcStructPtr
	}()

	valueOf := reflect.ValueOf(p.svcStructPtr)
	valueOfElem := valueOf.Elem()
	typeOf := valueOfElem.Type()
	if typeOf.Kind() != reflect.Struct {
		panic("invalid struct ptr")
	}

	numField := valueOfElem.NumField()
	for i := 0; i < numField; i++ {
		t := typeOf.Field(i)
		svcImplStructName := t.Tag.Get("service")
		if svcImplStructName == "" {
			// not autowire field
			continue
		}

		// get impled sub service
		impledPtr := impl(registeredMap[getInterfaceIdByNames(t.Type.Name(), svcImplStructName)])

		// set field
		subService := valueOfElem.Field(i)
		if !(subService.Kind() == reflect.Interface && subService.IsValid() && subService.CanSet()) {
			err := perrors.Errorf("Failed to autowire interface %s 's impl %s service. It's field %s with tag 'service:\"%s\"', please check if the field is exported",
				getName(p.interfaceStruct), getName(p.svcStructPtr), t.Type.Name(), svcImplStructName)
			panic(err)
		}
		subService.Set(reflect.ValueOf(impledPtr))
	}
	for _, f := range p.constructFunctions {
		f(p.svcStructPtr)
	}
	return p.svcStructPtr
}

func getInterfaceId(p RegisterServicePair) string {
	interfaceName := getName(p.interfaceStruct)
	structPtrName := getName(p.svcStructPtr)
	return getInterfaceIdByNames(interfaceName, structPtrName)
}

func getInterfaceIdByNames(interfaceName, structPtrName string) string {
	return strings.Join([]string{interfaceName, structPtrName}, "-")
}

func getName(v interface{}) string {
	typeOfInterface := getTypeFromInterface(v)
	return typeOfInterface.Name()
}

func getTypeFromInterface(v interface{}) reflect.Type {
	valueOfInterface := reflect.ValueOf(v)
	valueOfElemInterface := valueOfInterface.Elem()
	return valueOfElemInterface.Type()
}

package boot

import (
	"os"
	"reflect"
	"strings"
)

import (
	"github.com/glory-go/monkey"

	perrors "github.com/pkg/errors"
)

import (
	"github.com/glory-go/glory/log"
)

type RegisterServicePair struct {
	interfaceStruct    interface{}
	svcStructPtr       interface{}
	constructFunctions []func(interface{})
	guardMap           map[string]*monkey.PatchGuard
	isController       bool
}

var registeredMap = make(map[string]RegisterServicePair)
var implCompletedMap = make(map[string]interface{})
var grpcImplCompletedMap = make(map[string]interface{})
var controllerMap = make(map[string]interface{})

func RegisterController(controller interface{}) {
	controllerMap[getName(controller)] = controller
}

func RegisterService(interfaceStruct, structPtr interface{}, constructFunction ...func(interface{})) {
	newPair := RegisterServicePair{
		interfaceStruct:    interfaceStruct,
		svcStructPtr:       structPtr,
		constructFunctions: constructFunction,
		guardMap:           make(map[string]*monkey.PatchGuard),
		isController:       false,
	}

	serviceId := getInterfaceId(newPair)
	registeredMap[serviceId] = newPair
}

func Load(config map[interface{}]interface{}) {
	userConfig = config
	// set impl
	for serviceId, v := range registeredMap {
		if _, ok := implCompletedMap[serviceId]; ok {
			continue
		}
		impl(v)
	}

	// impl controller
	for _, v := range controllerMap {
		impl(RegisterServicePair{
			svcStructPtr: v,
			isController: true,
		})
	}
}

func impl(p RegisterServicePair) interface{} {
	tempInterfaceId := getInterfaceId(p)
	if impledPtr, ok := implCompletedMap[tempInterfaceId]; ok {
		// if already impleted, return
		return impledPtr
	}
	defer func() { // assure the impl procedure of one service run once
		if r := recover(); r != nil {
			log.Errorf("recover panic = %s", r)
		}
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
		var impledPtr interface{}
		tagKey := ""
		tagValue := ""
		if svcImplStructName := t.Tag.Get("service"); svcImplStructName != "" {
			// get impled sub local service
			fieldTypeName := t.Type.Name()
			if fieldTypeName == "" { // autowire struct ptr
				fieldTypeName = svcImplStructName
			}
			impledPtr = impl(registeredMap[getInterfaceIdByNames(fieldTypeName, svcImplStructName)])
			tagKey = "service"
			tagValue = svcImplStructName
		} else if grpcClientName := t.Tag.Get("grpc"); grpcClientName != "" {
			// `service:"grpc"` means auto wire grpc client
			impledPtr = implGRPC(t.Type.Name(), grpcClientName, t.Tag.Get("interceptorsKey"))
			tagKey = "grpc"
			tagValue = grpcClientName
		} else if configName := t.Tag.Get("config"); configName != "" {
			// XXX string `config:"mysvc.config"` means auto wire config
			configStr, err := implConfig(configName)
			if err != nil {
				panic(err)
			}
			tagKey = "config"
			tagValue = configName

			subService := valueOfElem.Field(i)
			if !(subService.Kind() == reflect.String && subService.IsValid() && subService.CanSet()) {
				err := perrors.Errorf("Failed to autowire interface %s's confige. It's field %s with tag '%s:\"%s\"', please check if the field is exported",
					getName(p.interfaceStruct), t.Type.Name(), tagKey, tagValue)
				panic(err)
			}
			subService.Set(reflect.ValueOf(configStr))
			continue
		}
		if tagKey == "" && tagValue == "" {
			continue
		}
		// set field
		subService := valueOfElem.Field(i)
		if !(subService.IsValid() && subService.CanSet()) {
			err := perrors.Errorf("Failed to autowire interface %s's impl %s service. It's field %s with tag '%s:\"%s\"', please check if the field is exported",
				getName(p.interfaceStruct), getName(p.svcStructPtr), t.Type.Name(), tagKey, tagValue)
			panic(err)
		}
		subService.Set(reflect.ValueOf(impledPtr))
	}
	for _, f := range p.constructFunctions {
		f(p.svcStructPtr)
	}
	// todo control if using monkey
	if !p.isController && os.Getenv("GOARCH") == "amd64" {
		// only service, only amd64 mod can inject monkey function
		implMonkey(p.svcStructPtr, tempInterfaceId)
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
	if v == nil {
		return ""
	}
	typeOfInterface := getTypeFromInterface(v)
	return typeOfInterface.Name()
}

func getTypeFromInterface(v interface{}) reflect.Type {
	valueOfInterface := reflect.ValueOf(v)
	valueOfElemInterface := valueOfInterface.Elem()
	return valueOfElemInterface.Type()
}

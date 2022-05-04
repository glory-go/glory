package boot

import (
	"os"
	"reflect"
)

import (
	"github.com/glory-go/monkey"

	perrors "github.com/pkg/errors"
)

import (
	"github.com/glory-go/glory/boot/common"
	"github.com/glory-go/glory/boot/interceptor"
	"github.com/glory-go/glory/boot/util"
	"github.com/glory-go/glory/config"
	"github.com/glory-go/glory/log"
)

var registeredMap = make(map[string]common.RegisterServiceMetadata)
var implCompletedMap = make(map[string]interface{})
var grpcImplCompletedMap = make(map[string]interface{})
var controllerMap = make(map[string]interface{})

func RegisterController(controller interface{}) {
	controllerMap[util.GetName(controller)] = controller
}

func RegisterService(interfaceStruct, structPtr interface{}, constructFunction ...func(interface{})) {
	newPair := common.RegisterServiceMetadata{
		InterfaceStruct:    interfaceStruct,
		SvcStructPtr:       structPtr,
		ConstructFunctions: constructFunction,
		GuardMap:           make(map[string]*monkey.PatchGuard),
		IsController:       false,
	}

	serviceId := getInterfaceId(newPair)
	registeredMap[serviceId] = newPair
}

func Load(config *config.ServerConfig) {
	userConfig = config.UserConfig
	debugConfig = config.DebugConfig
	// set impl
	for serviceId, v := range registeredMap {
		if _, ok := implCompletedMap[serviceId]; ok {
			continue
		}
		impl(v)
	}

	// impl controller
	for _, v := range controllerMap {
		impl(common.RegisterServiceMetadata{
			SvcStructPtr: v,
			IsController: true,
		})
	}

	// start debug port if enabled
	if enable := debugConfig["enable"]; enable != "false" {
		debugPort := debugConfig["port"]
		if debugPort == "" {
			debugPort = "1999"
		}
		go func() {
			if err := interceptor.Run(debugPort, registeredMap); err != nil {
				panic(err)
			}
		}()
	}
}

func impl(p common.RegisterServiceMetadata) interface{} {
	tempInterfaceId := getInterfaceId(p)
	if impledPtr, ok := implCompletedMap[tempInterfaceId]; ok {
		// if already impleted, return
		return impledPtr
	}
	defer func() { // assure the impl procedure of one service run once
		if r := recover(); r != nil {
			log.Errorf("recover panic = %s", r)
		}
		implCompletedMap[tempInterfaceId] = p.SvcStructPtr
	}()

	valueOf := reflect.ValueOf(p.SvcStructPtr)
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
			impledPtr = impl(registeredMap[util.GetInterfaceIdByNames(fieldTypeName, svcImplStructName)])
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
					util.GetName(p.InterfaceStruct), t.Type.Name(), tagKey, tagValue)
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
				util.GetName(p.InterfaceStruct), util.GetName(p.SvcStructPtr), t.Type.Name(), tagKey, tagValue)
			panic(err)
		}
		subService.Set(reflect.ValueOf(impledPtr))
	}
	for _, f := range p.ConstructFunctions {
		f(p.SvcStructPtr)
	}
	// todo control if using monkey
	if !p.IsController && os.Getenv("GOARCH") == "amd64" {
		// only service, only amd64 mod can inject monkey function
		implMonkey(p.SvcStructPtr, tempInterfaceId)
	}
	return p.SvcStructPtr
}

func getInterfaceId(p common.RegisterServiceMetadata) string {
	interfaceName := util.GetName(p.InterfaceStruct)
	structPtrName := util.GetName(p.SvcStructPtr)
	return util.GetInterfaceIdByNames(interfaceName, structPtrName)
}

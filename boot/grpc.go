package boot

import (
	"reflect"
)

import (
	"google.golang.org/grpc"
)

import (
	gloryGRPC "github.com/glory-go/glory/grpc"
)

var gRPCServiceMap = make(map[string]RegisterGRPCServicePair)
var gRPCUnaryClientInterceptorMap = make(map[string][]grpc.UnaryClientInterceptor)
var defaultUnaryClientInterceptor = make([]grpc.UnaryClientInterceptor, 0)

type RegisterGRPCServicePair struct {
	interfaceName string
	implFunction  interface{}
}

func RegisterGRPCService(f interface{}) {
	typeOf := reflect.TypeOf(f)
	gRPCInterfaceType := typeOf.Out(0)
	name := gRPCInterfaceType.Name()
	gRPCServiceMap[name] = RegisterGRPCServicePair{
		interfaceName: name,
		implFunction:  f,
	}
}

func RegisterGRPCUnaryClientInterceptorLists(name string, interceptors []grpc.UnaryClientInterceptor) {
	gRPCUnaryClientInterceptorMap[name] = interceptors
	if name == "default" {
		defaultUnaryClientInterceptor = interceptors
	}
}

func implGRPC(interfaceName, clientName, interceptorsKey string) interface{} {
	if impledPtr, ok := grpcImplCompletedMap[interfaceName]; ok {
		// if already impleted, return
		return impledPtr
	}

	f := gRPCServiceMap[interfaceName]
	valueOf := reflect.ValueOf(f.implFunction)
	// get conn
	interceptors := make([]grpc.UnaryClientInterceptor, 0)
	for _, interceptor := range defaultUnaryClientInterceptor {
		interceptors = append(interceptors, interceptor)
	}
	if interceptors, ok := gRPCUnaryClientInterceptorMap[interceptorsKey]; interceptorsKey != "" && ok {
		for _, interceptor := range interceptors {
			interceptors = append(interceptors, interceptor)
		}
	}
	grpcClient := gloryGRPC.NewGrpcClient(clientName, interceptors...)
	conn := grpcClient.GetConn()

	// call impl
	rsp := valueOf.Call([]reflect.Value{reflect.ValueOf(conn)})
	impledPtr := rsp[0].Interface()
	grpcImplCompletedMap[interfaceName] = impledPtr
	return impledPtr
}

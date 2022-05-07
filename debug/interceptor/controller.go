package interceptor

import (
	"log"
	"net"
)

import (
	"google.golang.org/grpc"
)

import (
	"github.com/glory-go/glory/debug/api/glory/boot"
	"github.com/glory-go/glory/debug/common"
)

func Start(port string, allInterfaceMetadataMap map[string]*common.DebugMetadata) error {
	grpcServer := grpc.NewServer()
	grpcServer.RegisterService(&boot.DebugService_ServiceDesc, &DebugServerImpl{
		editInterceptor:         GetEditInterceptor(),
		watchInterceptor:        GetWatchInterceptor(),
		allInterfaceMetadataMap: allInterfaceMetadataMap,
	})
	lst, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}
	go func() {
		if err := grpcServer.Serve(lst); err != nil {
			log.Println("debug server start failed with error = ", err)
			return
		}
	}()
	return nil
}

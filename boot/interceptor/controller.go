package interceptor

import (
	"net"
)

import (
	"google.golang.org/grpc"
)

import (
	"github.com/glory-go/glory/boot/api/glory/boot"
	"github.com/glory-go/glory/boot/common"
)

func Run(port string, allInterfaceMetadataMap map[string]common.RegisterServiceMetadata) error {
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
	if err := grpcServer.Serve(lst); err != nil {
		return err
	}
	return nil
}

package interceptor

import (
	"net"
)

import (
	"github.com/fatih/color"

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
		color.Red("[Debug] Debug server listening port :%s failed with error = %s", port, err)
		return err
	}
	go func() {
		color.Blue("[Debug] Debug server listening at :%s", port)
		if err := grpcServer.Serve(lst); err != nil {
			color.Red("[Debug] Debug server run with error = ", err)
			return
		}
	}()
	return nil
}

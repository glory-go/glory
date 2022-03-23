package service

import (
	"context"
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"

	"github.com/glory-go/glory/log"

	"github.com/glory-go/glory/config"
)

type jsonRPCService struct {
	DefaultServiceBase
}

func NewJsonRPCService(name string) *jsonRPCService {
	newJsonRPCService := &jsonRPCService{}
	newJsonRPCService.Name = name
	newJsonRPCService.LoadConfig(config.GlobalServerConf.ServiceConfigs[name])
	return newJsonRPCService
}

func (jrpcs *jsonRPCService) Run(ctx context.Context) {
	// start jsonRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", jrpcs.conf.addr.Port))
	if err != nil {
		log.Errorf("failed to listen grpc: %v", err)
	}
	fmt.Println("jsonrpc start listening on", fmt.Sprintf(":%v", jrpcs.conf.addr.Port))
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Errorf("jsonrpc listen err with %s", err.Error())
			continue
		}
		jsonrpc.ServeConn(conn)
	}
}

// Register 对用户暴露的接口
func (jrpcs *jsonRPCService) Register(rcvr interface{}) {
	// todo param check
	// 确保结构体所有函数可导出，所有函数的最后一个参数为指针
	if err := rpc.Register(rcvr); err != nil {
		log.Errorf("jsonrpc register struct: err = %s", err.Error())
	}
}

package jsonrpc

import (
	"net/rpc"
	"net/rpc/jsonrpc"

	"github.com/glory-go/glory/config"
	"github.com/glory-go/glory/log"
)

type jsonRPCClient struct {
	client        *rpc.Client
	targetAddress string
}

func (jrpc *jsonRPCClient) loadConfig(conf *config.JsonRPCClientConfig) {
	//if conf.ConfigSource == "env" {
	//	if err := tools.ReadAllConfFromEnv(conf); err != nil {
	//		log.Error("grpc-client: get conf struct err = ", err)
	//		return
	//	}
	//}
	jrpc.targetAddress = conf.ServerAddress
	log.Info(jrpc.targetAddress)
}

func (jrpc *jsonRPCClient) setup() {
	var err error
	jrpc.client, err = jsonrpc.Dial("tcp", jrpc.targetAddress)
	if err != nil {
		log.Error(err)
	}
}

func (jrpc *jsonRPCClient) Call(targetMethod string, params interface{}, rsp interface{}) error {
	// todo check 检查rsp是否为指针，params是否都为导出
	return jrpc.client.Call(targetMethod, params, rsp)
}

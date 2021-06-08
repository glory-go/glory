package jsonrpc

import "github.com/glory-go/glory/config"

func NewJsonRPCClient(jsonRPCClientName string) *jsonRPCClient {
	jsonRPCClient := &jsonRPCClient{}
	jsonRPCClient.loadConfig(config.GlobalServerConf.JsonRPCClientConfig[jsonRPCClientName])
	jsonRPCClient.setup()
	return jsonRPCClient
}

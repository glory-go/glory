package jsonrpc

import (
	"github.com/glory-go/glory/config"
)

func NewJsonRPCClient(jsonRPCClientName string) *jsonRPCClient {
	jsonRPCClient := &jsonRPCClient{}
	jsonRPCClient.loadConfig(config.GlobalServerConf.ClientConfig[jsonRPCClientName])
	jsonRPCClient.setup()
	return jsonRPCClient
}

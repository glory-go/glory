package glory

import (
	"github.com/glory-go/glory/grmanager"
	"github.com/glory-go/glory/server"
)

// NewServer 新建一个Glory框架服务
func NewServer() server.GloryServer {
	ctx := grmanager.NewCtx()
	server := server.NewDefaultGloryServer(ctx)
	server.LoadConfig()
	return server
}

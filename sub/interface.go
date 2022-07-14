package sub

import (
	"context"

	"github.com/glory-go/glory/v2/service"
	"github.com/lileio/pubsub"
)

//go:generate mockgen -source interface.go -destination mock/interface.go

type SubProvider interface {
	service.Service
	// Subscribe 实现了sub的订阅逻辑
	Subscribe(ctx context.Context, topic string, h pubsub.MsgHandler)
}

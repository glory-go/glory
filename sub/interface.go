package sub

import (
	"github.com/glory-go/glory/v2/service"
)

//go:generate mockgen -source interface.go -destination mock/interface.go

type SubProvider interface {
	service.Service
}

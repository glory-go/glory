package sub

import (
	"sync"

	"github.com/glory-go/glory/v2/service"
)

var (
	registerOnce sync.Once
)

func register() {
	registerOnce.Do(func() {
		service.GetService().RegisterService(GetSub())
	})
}

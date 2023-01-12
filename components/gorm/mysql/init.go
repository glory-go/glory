package mysql

import (
	"sync"

	"github.com/glory-go/glory/v2/config"
)

var (
	registerOnce sync.Once
)

func register() {
	registerOnce.Do(func() {
		config.RegisterComponent(GetMysqlGormComponent())
	})
}

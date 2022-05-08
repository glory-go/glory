package glory

import (
	"github.com/glory-go/glory/autowire"
	"github.com/glory-go/glory/config"
	"github.com/glory-go/glory/debug"
)

func Load() error {
	printLogo()

	// 1. load config
	if err := config.Load(); err != nil {
		return err
	}

	// 2. load debug
	if err := debug.Load(); err != nil {
		return err
	}

	// 3. load autowire
	return autowire.Load()
}

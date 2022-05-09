package glory

import (
	"github.com/fatih/color"
)

import (
	"github.com/glory-go/glory/autowire"
	"github.com/glory-go/glory/config"
	"github.com/glory-go/glory/debug"
)

func Load() error {
	printLogo()
	color.Cyan("Welcome to use glory-go!")

	// 1. load config
	color.Blue("[Boot] Start to load glory config")
	if err := config.Load(); err != nil {
		return err
	}

	// 2. load debug
	color.Blue("[Boot] Start to load debug")
	if err := debug.Load(); err != nil {
		return err
	}

	// 3. load autowire
	color.Blue("[Boot] Start to load autowire")
	return autowire.Load()
}

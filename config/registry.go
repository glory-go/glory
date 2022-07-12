package config

import "sync"

var (
	configCenterRegistry sync.Map
)

func RegisterConfigCenter(center ConfigCenter) {}

func LoadConfigCenter(name string) (ConfigCenter, bool) {}

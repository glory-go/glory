package config_center

import "github.com/glory-go/glory/config"

type ConfigCenter interface {
	Conn(config *config.ConfigCenterConfig) error
	LoadConfig() (*config.ServerConfig, error)
}

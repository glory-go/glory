package sqlite

import "github.com/glory-go/glory/v2/components/gorm"

type sqliteConfig struct {
	gorm.GormConfig
	Path string `mapstructure:"path"`
}

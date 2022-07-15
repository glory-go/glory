package mysql

import "github.com/glory-go/glory/v2/components/gorm"

type mysqlConfig struct {
	gorm.GormConfig
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Addr     string `mapstructure:"addr"`
	DB       string `mapstructure:"db"`
	Params   string `mapstructure:"params"`
}

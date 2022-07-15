package postgres

import "github.com/glory-go/glory/v2/components/gorm"

type postgresConfig struct {
	gorm.GormConfig
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	DB       string `mapstructure:"db"`
	TimeZone string `mapstructure:"time_zone"`
}

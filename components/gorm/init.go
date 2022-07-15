package gorm

import (
	"time"

	"gorm.io/gorm"
)

func Init(db *gorm.DB, conf *GormConfig) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	if conf.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(conf.MaxIdleConns)
	}
	if conf.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(conf.MaxOpenConns)
	}
	if conf.ConnMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(time.Duration(conf.ConnMaxLifetime) * time.Second)
	}

	return nil
}

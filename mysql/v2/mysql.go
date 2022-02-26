package mysql

import (
	"github.com/glory-go/glory/log"
	"github.com/glory-go/glory/service/middleware/jaeger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type mysqlService struct {
	conf *MysqlConfig
	db   *gorm.DB
}

func (srv *mysqlService) getMysqlLinkStr() string {
	return srv.conf.Username + ":" + srv.conf.Password + "@tcp(" + srv.conf.Host + ":" + srv.conf.Port + ")/" + srv.conf.DBName +
		"?charset=utf8&parseTime=True&loc=Local"
}

func (srv *mysqlService) conn() error {
	db, err := gorm.Open(mysql.Open(srv.getMysqlLinkStr()), &gorm.Config{
		Logger: NewGormLogger(),
	})
	if err != nil {
		log.Error("open db error ", err, "with db config = ", srv.conf)
		return err
	}
	if err := jaeger.GormUseTrace(srv.db); err != nil {
		log.Error("register tracer meets error", err)
		return nil
	}
	srv.db = db
	return nil
}

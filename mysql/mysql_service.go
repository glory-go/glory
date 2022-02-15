package mysql

import (
	"github.com/glory-go/glory/log"
	"github.com/glory-go/glory/service/middleware/jaeger"
	"gorm.io/gorm"

	"github.com/glory-go/glory/config"
	"gorm.io/driver/mysql"
)

// MysqlService 保存单个mysql链接
type MysqlService struct {
	DB   *gorm.DB
	conf config.MysqlConfig
}

func newMysqlService() *MysqlService {
	return &MysqlService{}
}

func getMysqlLinkStr(conf config.MysqlConfig) string {
	return conf.Username + ":" + conf.Password + "@tcp(" + conf.Host + ":" + conf.Port + ")/" + conf.DBName +
		"?charset=utf8&parseTime=True&loc=Local"
}

func (ms *MysqlService) loadConfig(conf config.MysqlConfig) error {
	ms.conf = conf
	return nil
}

func (ms *MysqlService) openDB(conf config.MysqlConfig) error {
	var err error
	if err := ms.loadConfig(conf); err != nil {
		log.Error("opendb error with err = ", err)
		return err
	}
	ms.DB, err = gorm.Open(mysql.Open(getMysqlLinkStr(ms.conf)), &gorm.Config{
		Logger: NewGormLogger(),
	})
	if err != nil {
		log.Error("open db error ", err, "with db config = ", ms.conf)
		return err
	}
	if err := jaeger.GormUseTrace(ms.DB); err != nil {
		log.Error("register tracer meets error", err)
		return nil
	}
	return nil
}

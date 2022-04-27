package mysql

import (
	"errors"
)

import (
	"gorm.io/driver/mysql"

	"gorm.io/gorm"
)

import (
	"github.com/glory-go/glory/config"
	"github.com/glory-go/glory/log"
	"github.com/glory-go/glory/service/middleware/jaeger"
)

// MysqlService 保存单个mysql链接
type MysqlService struct {
	DB     *gorm.DB
	tables map[string]*MysqlTable
	conf   config.MysqlConfig
}

func newMysqlService() *MysqlService {
	return &MysqlService{
		tables: make(map[string]*MysqlTable),
	}
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

func (ms *MysqlService) registerModel(model UserDefinedModel) (*MysqlTable, error) {
	table := newMysqlTable(ms.DB)
	if err := table.registerModel(model); err != nil {
		log.Error("mysql service register model err")
		return nil, err
	}

	return table, nil
}

func (ms *MysqlService) GetTable(tableName string) (*MysqlTable, error) {
	table, ok := ms.tables[tableName]
	if !ok {
		log.Error("table name = ", tableName, " not registered")
		return nil, errors.New("table name = " + tableName + " not registered")
	}
	return table, nil
}

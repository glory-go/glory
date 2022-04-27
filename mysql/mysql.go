package mysql

import (
	"errors"
)

import (
	"github.com/glory-go/glory/config"
	"github.com/glory-go/glory/log"
)

func init() {
	defaultMysqlHandler = newMysqlHandler()
	defaultMysqlHandler.setup(config.GlobalServerConf.MysqlConfigs)
}

type MysqlHandler struct {
	mysqlServices map[string]*MysqlService
}

var defaultMysqlHandler *MysqlHandler

func (mh *MysqlHandler) setup(conf map[string]*config.MysqlConfig) {
	for k, v := range conf {
		tempService := newMysqlService()
		if err := tempService.openDB(*v); err != nil {
			log.Errorf("opendb with key = %s, err = %s", k, err)
			continue
		}
		mh.mysqlServices[k] = tempService
	}
}

func newMysqlHandler() *MysqlHandler {
	return &MysqlHandler{
		mysqlServices: make(map[string]*MysqlService),
	}
}

func RegisterModel(mysqlServiceName string, model UserDefinedModel) (*MysqlTable, error) {
	service, ok := defaultMysqlHandler.mysqlServices[mysqlServiceName]
	if !ok {
		log.Error("mysql service name = ", mysqlServiceName, " not setup successful")
		return nil, errors.New("mysql service name = " + mysqlServiceName + " not setup successful")
	}
	return service.registerModel(model)
}

func GetService(mysqlServiceName string) (*MysqlService, bool) {
	s, ok := defaultMysqlHandler.mysqlServices[mysqlServiceName]
	return s, ok
}

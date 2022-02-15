package mysql

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

func GetService(mysqlServiceName string) (*MysqlService, bool) {
	s, ok := defaultMysqlHandler.mysqlServices[mysqlServiceName]
	return s, ok
}

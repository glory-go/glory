package mysql

import (
	"fmt"
	"sync"

	config "github.com/glory-go/glory/config/v2"
	"github.com/glory-go/glory/tools"
	"gorm.io/gorm"
)

const (
	ConfigName = "mysql"
)

func init() {
	config.RegisterConfig(ConfigName, mysqlBuilder)
}

var (
	mysqlServiceRegistry sync.Map
)

func mysqlBuilder(name string, rawConf map[string]interface{}) error {
	_, ok := mysqlServiceRegistry.Load(name)
	if ok {
		return fmt.Errorf("mysql with %s already registered", name)
	}
	conf := &MysqlConfig{}
	if err := tools.ConvertInto(rawConf, conf); err != nil {
		return err
	}

	srv := &mysqlService{
		conf: conf,
	}
	if err := srv.conn(); err != nil {
		return err
	}
	mysqlServiceRegistry.Store(name, srv)
	return nil
}

func GetDB(name string) (*gorm.DB, error) {
	srvInterface, ok := mysqlServiceRegistry.Load(name)
	if !ok {
		return nil, fmt.Errorf("mysql with name %s not registered", name)
	}
	srv := srvInterface.(*mysqlService)
	return srv.db, nil
}

package mysql

import (
	"fmt"
	"sync"

	gormcomponent "github.com/glory-go/glory/v2/components/gorm"
	"github.com/mitchellh/mapstructure"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	MysqlGormComponentName = "mysql_gorm"
)

type mysqlGormComponent struct {
	config map[string]*mysqlConfig
	dbs    map[string]*gorm.DB
}

var (
	component *mysqlGormComponent
	once      sync.Once
)

func GetMysqlGormComponent() *mysqlGormComponent {
	once.Do(func() {
		component = &mysqlGormComponent{
			config: make(map[string]*mysqlConfig),
			dbs:    make(map[string]*gorm.DB),
		}
	})
	return component
}

func (g *mysqlGormComponent) GetDB(name string) *gorm.DB {
	return g.dbs[name]
}

func (*mysqlGormComponent) Name() string { return MysqlGormComponentName }

func (g *mysqlGormComponent) Init(config map[string]any) error {
	for name := range config {
		raw := config[name]
		conf := &mysqlConfig{}
		if err := mapstructure.Decode(raw, conf); err != nil {
			return err
		}
		// 初始化数据库连接
		db, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s)/%s?%s", conf.Username, conf.Password, conf.Addr, conf.DB, conf.Params)), &gorm.Config{})
		if err != nil {
			return err
		}
		if err := gormcomponent.Init(db, &conf.GormConfig); err != nil {
			return err
		}
		g.dbs[name] = db
		g.config[name] = conf
	}

	return nil
}

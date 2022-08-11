package sqlite

import (
	"sync"

	gormcomponent "github.com/glory-go/glory/v2/components/gorm"
	"github.com/mitchellh/mapstructure"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	SqliteGormComponentName = "sqlite_gorm"
)

type sqliteGormComponent struct {
	config map[string]*sqliteConfig
	dbs    map[string]*gorm.DB
}

var (
	component *sqliteGormComponent
	once      sync.Once
)

func GetSqliteGormComponent() *sqliteGormComponent {
	once.Do(func() {
		component = &sqliteGormComponent{
			config: make(map[string]*sqliteConfig),
			dbs:    make(map[string]*gorm.DB),
		}
	})
	return component
}

func (g *sqliteGormComponent) GetDB(name string) *gorm.DB {
	return g.dbs[name]
}

func (*sqliteGormComponent) Name() string { return SqliteGormComponentName }

func (g *sqliteGormComponent) Init(config map[string]any) error {
	for name := range config {
		raw := config[name]
		conf := &sqliteConfig{}
		if err := mapstructure.Decode(raw, conf); err != nil {
			return err
		}
		// 初始化数据库连接
		db, err := gorm.Open(sqlite.Open(conf.Path), &gorm.Config{})
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

package postgres

import (
	"fmt"
	"sync"

	gormcomponent "github.com/glory-go/glory/v2/components/gorm"
	"github.com/mitchellh/mapstructure"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	PostgresGormComponentName = "postgres_gorm"
)

type postgresGormComponent struct {
	config map[string]*postgresConfig
	dbs    map[string]*gorm.DB
}

var (
	component *postgresGormComponent
	once      sync.Once
)

func GetPostgresGormComponent() *postgresGormComponent {
	once.Do(func() {
		component = &postgresGormComponent{
			config: make(map[string]*postgresConfig),
			dbs:    make(map[string]*gorm.DB),
		}
	})
	return component
}

func (g *postgresGormComponent) GetDB(name string) *gorm.DB {
	return g.dbs[name]
}

func (*postgresGormComponent) Name() string { return PostgresGormComponentName }

func (g *postgresGormComponent) Init(config map[string]any) error {
	for name := range config {
		raw := config[name]
		conf := &postgresConfig{}
		if err := mapstructure.Decode(raw, conf); err != nil {
			return err
		}
		// 初始化数据库连接
		db, err := gorm.Open(postgres.Open(fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable TimeZone=%s",
			conf.Username, conf.Password, conf.Host, conf.Port, conf.DB, conf.TimeZone)), &gorm.Config{})
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

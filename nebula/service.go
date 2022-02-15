package nebula

import (
	"strconv"

	"github.com/glory-go/glory/config"
	"github.com/glory-go/glory/log"
	nebula "github.com/vesoft-inc/nebula-go/v2"
)

// MysqlService 保存单个mysql链接
type NebulaService struct {
	pool *nebula.ConnectionPool
	conf config.NebulaConfig
}

func newNebulaService() *NebulaService {
	return &NebulaService{}
}

func (ms *NebulaService) loadConfig(conf config.NebulaConfig) error {
	ms.conf = conf
	return nil
}

func (ms *NebulaService) openDB(conf config.NebulaConfig) error {
	var err error
	if err := ms.loadConfig(conf); err != nil {
		log.Error("opendb error with err = ", err)
		return err
	}
	if conf.Port == "" {
		conf.Port = "9669"
	}
	port, err := strconv.Atoi(conf.Port)
	if err != nil {
		log.Errorf("nebula fail to open db when parse port in conf, err: %v", err)
		return err
	}
	ms.pool, err = nebula.NewConnectionPool([]nebula.HostAddress{
		{
			Host: conf.Host,
			Port: port,
		},
	}, nebula.GetDefaultConf(), logger)
	if err != nil {
		log.Error("open db error ", err, "with db config = ", ms.conf)
		return err
	}
	return nil
}

func (ms *NebulaService) GetSession() (*nebula.Session, error) {
	session, err := ms.pool.GetSession(ms.conf.Username, ms.conf.Password)
	// TODO: 使用jager打点
	return session, err
}

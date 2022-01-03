package redis

import (
	"sync"

	"github.com/glory-go/glory/log"
	"github.com/go-redis/redis/v8"

	"github.com/glory-go/glory/config"
	_ "github.com/go-sql-driver/mysql"
)

// RedisService 保存多个redis的库
type RedisService struct {
	db     map[int]*redis.Client
	conf   config.RedisConfig
	isInit bool
	sync.RWMutex
}

func newRedisService() *RedisService {
	return &RedisService{
		db: make(map[int]*redis.Client),
	}
}

func (ms *RedisService) loadConfig(conf config.RedisConfig) error {
	ms.conf = conf
	return nil
}

func (ms *RedisService) openDB(conf config.RedisConfig) error {
	if err := ms.loadConfig(conf); err != nil {
		log.Error("opendb error with err = ", err)
		return err
	}
	return nil
}

func (ms *RedisService) registerModel(db int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     ms.conf.Host + ":" + ms.conf.Port,
		Password: ms.conf.Password,
		DB:       db,
	})
}

package redis

import (
	"errors"

	"github.com/glory-go/glory/config"
	"github.com/glory-go/glory/log"
	"github.com/go-redis/redis/v8"
)

func init() {
	defaultRedisHandler = newRedisHandler()
}

type RedisHandler struct {
	redisServices map[string]*RedisService
}

var defaultRedisHandler *RedisHandler

func (mh *RedisHandler) setup(conf map[string]*config.RedisConfig) {
	for k, v := range conf {
		tempService := newRedisService()
		tempService.conf = *v
		mh.redisServices[k] = tempService
	}
}

func newRedisHandler() *RedisHandler {
	return &RedisHandler{
		redisServices: make(map[string]*RedisService),
	}
}

func NewRedisClient(redisServiceName string, db int) (*redis.Client, error) {
	defaultRedisHandler.setup(config.GlobalServerConf.RedisConfig)
	service, ok := defaultRedisHandler.redisServices[redisServiceName]
	if !ok {
		log.Error("redis service name = ", redisServiceName, " not registered in config")
		return nil, errors.New("mysql service name = " + redisServiceName + " not registered in config")
	}
	// 检查是否初始化，未初始化则进行初始化，这里是为了保证使用时才去初始化
	service.RLock()
	if !service.isInit {
		service.RUnlock()
		service.Lock()
		defer service.Unlock()
		if err := service.openDB(service.conf); err != nil {
			log.Error("connect redis with key = ", redisServiceName, "err")
			service.isInit = true
			return nil, err
		}
	}
	return service.registerModel(db), nil
}

func GetService(mysqlServiceName string) (*RedisService, bool) {
	s, ok := defaultRedisHandler.redisServices[mysqlServiceName]
	return s, ok
}

package redis

import (
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/mitchellh/mapstructure"
)

const (
	RedisComponentName = "redis_6"
)

type redisComponent struct {
	config  map[string]*redisConfig
	clients map[string]*redis.Client
}

var (
	component *redisComponent
	once      sync.Once
)

func getRedisComponent() *redisComponent {
	once.Do(func() {
		component = &redisComponent{
			config:  map[string]*redisConfig{},
			clients: map[string]*redis.Client{},
		}
	})
	return component
}

func GetRedisClient(name string) *redis.Client {
	return getRedisComponent().clients[name]
}

func (c *redisComponent) Name() string { return RedisComponentName }

func (c *redisComponent) Init(config map[string]any) error {
	for name, raw := range config {
		conf := &redisConfig{}
		if err := mapstructure.Decode(raw, conf); err != nil {
			return err
		}
		c.config[name] = conf

		redisclient := redis.NewClient(&redis.Options{
			Addr:     conf.Addr,
			Username: conf.Username,
			Password: conf.Password,
			DB:       conf.DB,
		})
		c.clients[name] = redisclient
	}

	return nil
}

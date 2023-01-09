package tests

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	gloryredis "github.com/glory-go/glory/v2/components/redis/v8"
	"github.com/glory-go/glory/v2/config"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func Test_Redis(t *testing.T) {
	setup()
	s1 := gloryredis.GetRedisClient("mock_server_1")
	s2 := gloryredis.GetRedisClient("mock_server_2")
	// 确保s1和s2都可以正常调用，且数据不互通
	ctx := context.Background()
	key1, key2 := "test_1", "test_2"
	s1.Set(ctx, key1, 1, time.Hour)
	s2.Set(ctx, key2, 1, time.Hour)

	v1, err := s1.Get(ctx, key1).Int()
	assert.Nil(t, err)
	assert.Equal(t, 1, v1)
	_, err = s1.Get(ctx, key2).Int()
	assert.Equal(t, redis.Nil, err)

	_, err = s2.Get(ctx, key1).Int()
	assert.Equal(t, redis.Nil, err)
	v2, err := s2.Get(ctx, key2).Int()
	assert.Nil(t, err)
	assert.Equal(t, 1, v2)
}

func setup() {
	// 初始化redis连接
	redis1 := miniredis.NewMiniRedis()
	redis1.Start()
	redis2 := miniredis.NewMiniRedis()
	redis2.Start()
	// 初始化配置
	content := fmt.Sprintf(`
redis_6:
  mock_server_1:
    addr: %s
  mock_server_2:
    addr: %s
`, redis1.Addr(), redis2.Addr())
	path := fmt.Sprintf("/tmp/test_redis_%d.yaml", time.Now().Unix())
	config.ChangeDefaultConfigPath(path)
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	_, err = file.WriteString(content)
	if err != nil {
		panic(err)
	}
	config.Init()
}

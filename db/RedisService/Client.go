package RedisService

import (
	"context"
	"github.com/StephenChristianW/go-movies-open/config"
	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()
var RedisClient *redis.Client

func init() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: config.RedisPassword, // 没有密码就 ""
		DB:       config.RedisDB,       // 默认 0
	})
	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		panic("Redis 连接失败: " + err.Error())
	}
}

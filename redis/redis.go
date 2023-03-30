package redis

import (
	"github.com/redis/go-redis/v9"
	gc "xiuianserver/config"
)

var RedisClient *redis.Client

func NewRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     gc.GlobalConfig.RedisConfig.RemoteAddr,
		Password: gc.GlobalConfig.RedisConfig.Pass, // 没有密码，默认值
		DB:       gc.GlobalConfig.RedisConfig.DB,   // 默认DB 0
	})
}

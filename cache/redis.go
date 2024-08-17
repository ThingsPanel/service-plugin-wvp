package cache

import (
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

var REDIS *redis.Client

func RedisInit() {
	REDIS = redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.addr"),     // Redis 服务器地址
		Password: viper.GetString("redis.password"), // 没有密码时保持为空
		DB:       viper.GetInt("redis.db"),          // 使用默认的 DB
	})
}

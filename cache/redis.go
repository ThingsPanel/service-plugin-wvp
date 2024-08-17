package cache

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"plugin_wvp/model"
)

var REDIS *redis.Client

func RedisInit() {
	REDIS = redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.addr"),     // Redis 服务器地址
		Password: viper.GetString("redis.password"), // 没有密码时保持为空
		DB:       viper.GetInt("redis.db"),          // 使用默认的 DB
	})
}

func SetWvpConfig(ctx context.Context, config model.WvpForm) error {
	redisKey := "wvpConfigRedisCacheKey"
	configRedisKey := fmt.Sprintf("%s:%d", config.Server, config.Port)
	REDIS.SAdd(ctx, redisKey, configRedisKey)

	_, err := REDIS.Set(ctx, configRedisKey, config, 0).Result()
	if err != nil {
		return err
	}
	return nil
}

func GetWvpConfigKey(ctx context.Context) ([]string, error) {
	redisKey := "wvpConfigRedisCacheKey"
	return REDIS.SMembers(ctx, redisKey).Result()
}

func GetWvpConfig(ctx context.Context, key string) (model.WvpForm, error) {
	var result model.WvpForm
	return result, REDIS.Get(ctx, key).Scan(&result)
}

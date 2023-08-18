package redisclient

import (
	"github.com/dathuynh1108/clean-arch-base/pkg/config"
	"github.com/dathuynh1108/clean-arch-base/pkg/singleton"
	"github.com/redis/go-redis/v9"
)

var (
	redisClientSingleton *singleton.Singleton[redis.UniversalClient]
)

func InitRedis() error {
	config := config.GetConfig()
	redisClientSingleton = singleton.NewSingleton[redis.UniversalClient](func() redis.UniversalClient {
		return redis.NewUniversalClient(&redis.UniversalOptions{
			Addrs:    config.RedisConfig.Addrs,
			Username: config.RedisConfig.Username,
			Password: config.RedisConfig.Password,
			DB:       config.RedisConfig.DB,
		})
	})
	return nil
}

func GetRedis() redis.UniversalClient {
	return redisClientSingleton.Get()
}

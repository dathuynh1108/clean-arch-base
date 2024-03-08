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
	redisClientSingleton = singleton.NewSingleton(func() redis.UniversalClient {
		dsn := config.RedisConfig.DSN
		if config.RedisConfig.IsCluster {
			clusterOpts, err := redis.ParseClusterURL(dsn)
			if err != nil {
				panic(err)
			}
			return redis.NewClusterClient(clusterOpts)
		}

		opts, err := redis.ParseURL(dsn)
		if err != nil {
			panic(err)
		}

		return redis.NewClient(opts)
	}, true)
	return nil
}

func GetRedis() redis.UniversalClient {
	return redisClientSingleton.Get()
}

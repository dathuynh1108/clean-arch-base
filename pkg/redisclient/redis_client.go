package redisclient

import (
	"context"

	"github.com/dathuynh1108/clean-arch-base/pkg/config"
	"github.com/dathuynh1108/clean-arch-base/pkg/logger"
	"github.com/dathuynh1108/clean-arch-base/pkg/singleton"

	"github.com/redis/go-redis/v9"
)

var (
	redisClientSingleton *singleton.Singleton[redis.UniversalClient]
)

func InitRedis() error {
	config := config.GetConfig()
	dsn := config.RedisConfig.DSN
	var client redis.UniversalClient
	if config.RedisConfig.IsCluster {
		logger.GetLogger().Infof("Init redis cluster")
		clusterOpts, err := redis.ParseClusterURL(dsn)
		if err != nil {
			return err
		}
		client = redis.NewClusterClient(clusterOpts)
	} else {
		logger.GetLogger().Infof("Init redis client")
		opts, err := redis.ParseURL(dsn)
		if err != nil {
			return err
		}
		client = redis.NewClient(opts)
	}
	redisClientSingleton = singleton.NewSingletonInstance(client)
	return redisClientSingleton.Get().Ping(context.Background()).Err()
}

func GetRedis() redis.UniversalClient {
	return redisClientSingleton.Get()
}

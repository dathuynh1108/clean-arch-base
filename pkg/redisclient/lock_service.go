package redisclient

import (
	"github.com/dathuynh1108/clean-arch-base/pkg/singleton"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
)

var (
	lockServiceSingleton *singleton.Singleton[*redsync.Redsync]
)

func InitLockService() error {
	lockServiceSingleton = singleton.NewSingleton(func() *redsync.Redsync {
		return redsync.New(goredis.NewPool(GetRedis()))
	}, true)
	return nil
}

package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/dathuynh1108/clean-arch-base/pkg/comerr"
	"github.com/dathuynh1108/clean-arch-base/pkg/msgpack"
	"github.com/dathuynh1108/clean-arch-base/pkg/redisclient"
)

var (
	GetCache = redisclient.GetRedis
)

func GetFast[T any](
	ctx context.Context,
	cache redis.UniversalClient,
	key string,
	ttl time.Duration,
	creator func(ctx context.Context) (T, error),
) (value T, err error) {
	key = decorator(key)

	cmd := cache.Get(ctx, key)
	if cmd.Err() != nil {
		if errors.Is(cmd.Err(), redis.Nil) {
			value, err := creator(ctx)
			if err == nil {
				text, err := msgpack.MsgPackMarshal(value, false)
				if err == nil {
					_ = cache.Set(ctx, key, text, ttl).Err()
				}
			}
			return value, err
		}
		err = comerr.WrapMessage(cmd.Err(), "failed to get the value from the cache")
		return
	}

	err = msgpack.MsgPackUnmarshal(
		[]byte(cmd.Val()),
		&value,
		false,
	)
	if err != nil {
		err = comerr.WrapMessage(err, "failed to unmarshal the value from the cache")
		return
	}
	return value, err
}

func Set[T any](
	ctx context.Context,
	cache redis.UniversalClient,
	key string,
	ttl time.Duration,
	value T,
) (err error) {
	key = decorator(key)

	text, err := msgpack.MsgPackMarshal(value, false)
	if err != nil {
		return comerr.WrapMessage(err, "failed to marshal the value for the cache")
	}

	return cache.Set(ctx, key, text, ttl).Err()
}

func SetBatch[T any](
	ctx context.Context,
	cache redis.UniversalClient,
	ttl time.Duration,
	values map[string]T,
) (err error) {
	_, err = cache.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for k, v := range values {
			text, err := msgpack.MsgPackMarshal(v, false)
			if err != nil {
				return comerr.WrapMessage(err, "failed to marshal the value for the cache")
			}
			pipe.Set(ctx, decorator(k), text, ttl)
		}
		return nil
	})
	if err != nil {
		return comerr.WrapMessage(err, "failed to set data pipeline")
	}
	return
}

func decorator(key string) string {
	return fmt.Sprintf("%s:%s", "", key)
}

package redis

import (
	"context"
	"errors"
	"time"

	pkgerrors "github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

func setSingleValue[T string | int64 | float64](
	ctx context.Context,
	rdb redis.Cmdable,
	key string,
	value T,
	expiration time.Duration,
	mode setMode,
) error {
	// 1. Prepare SET command arguments
	args := redis.SetArgs{
		Mode:    mode.String(),
		KeepTTL: true,
	}
	if expiration > 0 {
		args.KeepTTL = false
		args.TTL = expiration
	}

	// 2. Set value to redis
	status, err := rdb.SetArgs(ctx, key, value, args).Result()
	if err != nil {
		return pkgerrors.WithStack(err)
	}

	if mode != setModeNone && status != statusOK {
		return ErrFailToSetValue
	}

	// 3. Return result
	return nil
}

func getSingleValue[T string | int64 | float64](ctx context.Context, rdb redis.Cmdable, key string) (T, error) {
	var rs T
	if err := rdb.Get(ctx, key).Scan(&rs); err != nil {
		if errors.Is(err, redis.Nil) {
			return *new(T), nil
		}

		return rs, pkgerrors.WithStack(err)
	}

	return rs, nil
}

func (client redisClient) SetString(ctx context.Context, key string, value string, expiration time.Duration) error {
	return setSingleValue(ctx, client.rdb, key, value, expiration, setModeNone)
}

func (client redisClient) SetStringIfNotExist(ctx context.Context, key string, value string, expiration time.Duration) error {
	return setSingleValue(ctx, client.rdb, key, value, expiration, setModeNX)
}

func (client redisClient) SetStringIfExist(ctx context.Context, key string, value string, expiration time.Duration) error {
	return setSingleValue(ctx, client.rdb, key, value, expiration, setModeXX)
}

func (client redisClient) GetString(ctx context.Context, key string) (string, error) {
	return getSingleValue[string](ctx, client.rdb, key)
}

func (client redisClient) SetInt(ctx context.Context, key string, value int64, expiration time.Duration) error {
	return setSingleValue(ctx, client.rdb, key, value, expiration, setModeNone)
}

func (client redisClient) SetIntIfNotExist(ctx context.Context, key string, value int64, expiration time.Duration) error {
	return setSingleValue(ctx, client.rdb, key, value, expiration, setModeNX)
}

func (client redisClient) SetIntIfExist(ctx context.Context, key string, value int64, expiration time.Duration) error {
	return setSingleValue(ctx, client.rdb, key, value, expiration, setModeXX)
}

func (client redisClient) GetInt(ctx context.Context, key string) (int64, error) {
	return getSingleValue[int64](ctx, client.rdb, key)
}

func (client redisClient) IncrementBy(ctx context.Context, key string, value int64) (int64, error) {
	val, err := client.rdb.IncrBy(ctx, key, value).Result()
	if err != nil {
		return 0, pkgerrors.WithStack(err)
	}

	return val, nil
}

func (client redisClient) DecrementBy(ctx context.Context, key string, value int64) (int64, error) {
	val, err := client.rdb.DecrBy(ctx, key, value).Result()
	if err != nil {
		return 0, pkgerrors.WithStack(err)
	}

	return val, nil
}

func (client redisClient) SetFloat(ctx context.Context, key string, value float64, expiration time.Duration) error {
	return setSingleValue(ctx, client.rdb, key, value, expiration, setModeNone)
}

func (client redisClient) SetFloatIfNotExist(ctx context.Context, key string, value float64, expiration time.Duration) error {
	return setSingleValue(ctx, client.rdb, key, value, expiration, setModeNX)
}

func (client redisClient) SetFloatIfExist(ctx context.Context, key string, value float64, expiration time.Duration) error {
	return setSingleValue(ctx, client.rdb, key, value, expiration, setModeXX)
}

func (client redisClient) IncrementFloatBy(ctx context.Context, key string, value float64) (float64, error) {
	val, err := client.rdb.IncrByFloat(ctx, key, value).Result()
	if err != nil {
		return 0, pkgerrors.WithStack(err)
	}

	return val, nil
}

func (client redisClient) GetFloat(ctx context.Context, key string) (float64, error) {
	return getSingleValue[float64](ctx, client.rdb, key)
}

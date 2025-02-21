package redis

import (
	"context"
	"time"

	pkgerrors "github.com/pkg/errors"
)

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

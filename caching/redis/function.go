package redis

import (
	"context"
	"errors"
	"time"

	pkgerrors "github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

func setSingleValue[T Type](
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

func getSingleValue[T Type](ctx context.Context, rdb redis.Cmdable, key string) (T, error) {
	var rs T
	if err := rdb.Get(ctx, key).Scan(&rs); err != nil {
		if errors.Is(err, redis.Nil) {
			return *new(T), nil
		}

		return rs, pkgerrors.WithStack(err)
	}

	return rs, nil
}

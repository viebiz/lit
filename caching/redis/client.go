package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"

	pkgerrors "github.com/pkg/errors"
)

type redisClient struct {
	rdb redis.UniversalClient
}

func (client redisClient) Ping(ctx context.Context) error {
	if _, err := client.rdb.Ping(ctx).Result(); err != nil {
		return pkgerrors.WithStack(err)
	}

	return nil
}

func (client redisClient) Do(ctx context.Context, cmd string, args ...interface{}) (interface{}, error) {
	params := make([]interface{}, len(args)+1)
	params[0] = cmd
	copy(params[1:], args)

	val, err := client.rdb.Do(ctx, params...).Result()
	if err != nil {
		return nil, pkgerrors.WithStack(err)
	}

	return val, nil
}

// DoInBatch executes multiple Redis commands in a single call using a pipeline.
// Note: All commands are executed at once, and no results are returned until the function completes.
func (client redisClient) DoInBatch(ctx context.Context, fn func(cmder Commander) error) error {
	var pipelineFn = func(pl redis.Pipeliner) error {
		return fn(commander{pipeliner: pl})
	}

	if _, err := client.rdb.Pipelined(ctx, pipelineFn); err != nil {
		return pkgerrors.WithStack(err)
	}

	return nil
}

func (client redisClient) Close() error {
	return client.rdb.Close()
}

func (client redisClient) Delete(ctx context.Context, key string) (int64, error) {
	rs, err := client.rdb.Del(ctx, key).Result()
	if err != nil {
		return 0, pkgerrors.WithStack(err)
	}

	return rs, nil
}

func (client redisClient) Expire(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	rs, err := client.rdb.Expire(ctx, key, expiration).Result()
	if err != nil {
		return false, pkgerrors.WithStack(err)
	}

	return rs, nil
}

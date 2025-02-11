package redis

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/redis/go-redis/v9"

	pkgerrors "github.com/pkg/errors"

	"github.com/viebiz/lit/monitoring"
)

type Client interface {
	Ping(ctx context.Context) error

	Do(ctx context.Context, cmd string, args ...interface{}) (interface{}, error)

	DoInBatch(ctx context.Context, fn func(cmder Commander) error) error

	Close() error

	Delete(ctx context.Context, key string) (int64, error)

	Expire(ctx context.Context, key string, expiration time.Duration) (bool, error)

	SetString(ctx context.Context, key string, value string, expiration time.Duration) error

	SetStringIfNotExist(ctx context.Context, key string, value string, expiration time.Duration) error

	SetStringIfExist(ctx context.Context, key string, value string, expiration time.Duration) error

	SetInt(ctx context.Context, key string, value int64, expiration time.Duration) error

	SetIntIfNotExist(ctx context.Context, key string, value int64, expiration time.Duration) error

	SetIntIfExist(ctx context.Context, key string, value int64, expiration time.Duration) error

	IncrementBy(ctx context.Context, key string, value int64) (int64, error)

	DecrementBy(ctx context.Context, key string, value int64) (int64, error)

	SetFloat(ctx context.Context, key string, value float64, expiration time.Duration) error

	SetFloatIfNotExist(ctx context.Context, key string, value float64, expiration time.Duration) error

	SetFloatIfExist(ctx context.Context, key string, value float64, expiration time.Duration) error

	IncrementFloatBy(ctx context.Context, key string, value float64) (float64, error)

	GetString(ctx context.Context, key string) (string, error)

	GetInt(ctx context.Context, key string) (int64, error)

	GetFloat(ctx context.Context, key string) (float64, error)

	HashSet(ctx context.Context, key string, value interface{}) error

	HashGetAll(ctx context.Context, key string, out interface{}) error

	HashGetField(ctx context.Context, key string, field string, out interface{}) error
}

type redisClient struct {
	rdb redis.UniversalClient
}

// NewClient creates new redis client
//
// Examples:
//
//	client, err := NewClient("redis://user:password@localhost:6379/0?protocol=3")
func NewClient(url string) (Client, error) {
	return NewClientWithTLS(url, nil)
}

// NewClientWithTLS creates new redis client with TLS configuration
func NewClientWithTLS(url string, tlsConfig *tls.Config) (Client, error) {
	// 1. Prepare Redis client configurations
	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, pkgerrors.WithStack(err)
	}

	// 1.2. Add TLS config to Redis client options
	opts.TLSConfig = tlsConfig

	// 2. Create new Redis client
	rdb := redis.NewClient(opts)

	// 3. Setup Redis tracing hook
	info := monitoring.NewExternalServiceInfo(url)
	rdb.AddHook(newTracingHook(info))

	return redisClient{
		rdb: rdb,
	}, nil
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
	pipelineFunc := func(pl redis.Pipeliner) error {
		return fn(commander{pipeliner: pl})
	}

	if _, err := client.rdb.Pipelined(ctx, pipelineFunc); err != nil {
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

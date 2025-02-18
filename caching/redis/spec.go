package redis

import (
	"context"
	"time"
)

// Client represents redis client interface
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

// Commander represents redis pipeline supported commands
type Commander interface {
	Discard()

	Execute(ctx context.Context) error

	Delete(ctx context.Context, key string) (int64, error)

	Expire(ctx context.Context, key string, expiration time.Duration) (bool, error)

	SetString(ctx context.Context, key string, value string, expiration time.Duration) error

	SetStringIfNotExist(ctx context.Context, key string, value string, expiration time.Duration) error

	SetStringIfExist(ctx context.Context, key string, value string, expiration time.Duration) error

	GetString(ctx context.Context, key string) (string, error)

	SetInt(ctx context.Context, key string, value int64, expiration time.Duration) error

	SetIntIfNotExist(ctx context.Context, key string, value int64, expiration time.Duration) error

	SetIntIfExist(ctx context.Context, key string, value int64, expiration time.Duration) error

	IncrementBy(ctx context.Context, key string, value int64) (int64, error)

	DecrementBy(ctx context.Context, key string, value int64) (int64, error)

	GetInt(ctx context.Context, key string) (int64, error)

	SetFloat(ctx context.Context, key string, value float64, expiration time.Duration) error

	SetFloatIfNotExist(ctx context.Context, key string, value float64, expiration time.Duration) error

	SetFloatIfExist(ctx context.Context, key string, value float64, expiration time.Duration) error

	IncrementFloatBy(ctx context.Context, key string, value float64) (float64, error)

	GetFloat(ctx context.Context, key string) (float64, error)

	HashSet(ctx context.Context, key string, value interface{}) error

	HashGetAll(ctx context.Context, key string, out interface{}) error

	HashGetField(ctx context.Context, key string, field string, out interface{}) error
}

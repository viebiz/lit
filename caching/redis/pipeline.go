package redis

import (
	"context"
	"reflect"
	"time"

	pkgerrors "github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

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

type commander struct {
	pipeliner redis.Pipeliner
}

func (e commander) Discard() {
	e.pipeliner.Discard()
}

func (e commander) Execute(ctx context.Context) error {
	if _, err := e.pipeliner.Exec(ctx); err != nil {
		return pkgerrors.WithStack(err)
	}

	return nil
}

func (e commander) Delete(ctx context.Context, key string) (int64, error) {
	rs, err := e.pipeliner.Del(ctx, key).Result()
	if err != nil {
		return 0, pkgerrors.WithStack(err)
	}

	return rs, nil
}

func (e commander) Expire(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	rs, err := e.pipeliner.Expire(ctx, key, expiration).Result()
	if err != nil {
		return false, pkgerrors.WithStack(err)
	}

	return rs, nil
}

func (e commander) SetString(ctx context.Context, key string, value string, expiration time.Duration) error {
	return setSingleValue(ctx, e.pipeliner, key, value, expiration, setModeNone)
}

func (e commander) SetStringIfNotExist(ctx context.Context, key string, value string, expiration time.Duration) error {
	return setSingleValue(ctx, e.pipeliner, key, value, expiration, setModeNX)
}

func (e commander) SetStringIfExist(ctx context.Context, key string, value string, expiration time.Duration) error {
	return setSingleValue(ctx, e.pipeliner, key, value, expiration, setModeXX)
}

func (e commander) GetString(ctx context.Context, key string) (string, error) {
	return getSingleValue[string](ctx, e.pipeliner, key)
}

func (e commander) SetInt(ctx context.Context, key string, value int64, expiration time.Duration) error {
	return setSingleValue(ctx, e.pipeliner, key, value, expiration, setModeNone)
}

func (e commander) SetIntIfNotExist(ctx context.Context, key string, value int64, expiration time.Duration) error {
	return setSingleValue(ctx, e.pipeliner, key, value, expiration, setModeNX)
}

func (e commander) SetIntIfExist(ctx context.Context, key string, value int64, expiration time.Duration) error {
	return setSingleValue(ctx, e.pipeliner, key, value, expiration, setModeXX)
}

func (e commander) IncrementBy(ctx context.Context, key string, value int64) (int64, error) {
	rs, err := e.pipeliner.IncrBy(ctx, key, value).Result()
	if err != nil {
		return 0, pkgerrors.WithStack(err)
	}

	return rs, nil
}

func (e commander) DecrementBy(ctx context.Context, key string, value int64) (int64, error) {
	rs, err := e.pipeliner.DecrBy(ctx, key, value).Result()
	if err != nil {
		return 0, pkgerrors.WithStack(err)
	}

	return rs, nil
}

func (e commander) GetInt(ctx context.Context, key string) (int64, error) {
	return getSingleValue[int64](ctx, e.pipeliner, key)
}

func (e commander) SetFloat(ctx context.Context, key string, value float64, expiration time.Duration) error {
	return setSingleValue(ctx, e.pipeliner, key, value, expiration, setModeNone)
}

func (e commander) SetFloatIfNotExist(ctx context.Context, key string, value float64, expiration time.Duration) error {
	return setSingleValue(ctx, e.pipeliner, key, value, expiration, setModeNX)
}

func (e commander) SetFloatIfExist(ctx context.Context, key string, value float64, expiration time.Duration) error {
	return setSingleValue(ctx, e.pipeliner, key, value, expiration, setModeXX)
}

func (e commander) IncrementFloatBy(ctx context.Context, key string, value float64) (float64, error) {
	rs, err := e.pipeliner.IncrByFloat(ctx, key, value).Result()
	if err != nil {
		return 0, pkgerrors.WithStack(err)
	}

	return rs, nil
}

func (e commander) GetFloat(ctx context.Context, key string) (float64, error) {
	return getSingleValue[float64](ctx, e.pipeliner, key)
}

func (e commander) HashSet(ctx context.Context, key string, value interface{}) error {
	// 1. Return error if given value input is not in Struct or map
	if reflect.TypeOf(value).Kind() != reflect.Struct && reflect.TypeOf(value).Kind() != reflect.Map {
		return ErrUnsupportedInputType
	}

	// 2. Set values to redis
	v, err := e.pipeliner.HSet(ctx, key, value).Result()
	if err != nil {
		return pkgerrors.WithStack(err)
	}

	// 3. Return error if no any field set
	if v == 0 {
		return ErrFailToSetValue
	}

	return nil
}

func (e commander) HashGetAll(ctx context.Context, key string, out interface{}) error {
	// 1. Return error if given input is not pointer
	if reflect.TypeOf(out).Kind() != reflect.Ptr {
		return ErrUnsupportedInputType
	}

	// 2. Get all values from redis by key
	if err := e.pipeliner.HGetAll(ctx, key).Scan(out); err != nil {
		return pkgerrors.WithStack(err)
	}

	return nil
}

func (e commander) HashGetField(ctx context.Context, key string, field string, out interface{}) error {
	// 1. Return error if given input is not pointer
	if reflect.TypeOf(out).Kind() != reflect.Ptr {
		return ErrUnsupportedInputType
	}

	// 2. Get a field from redis by key and field name
	if err := e.pipeliner.HGet(ctx, key, field).Scan(out); err != nil {
		return pkgerrors.WithStack(err)
	}

	return nil
}

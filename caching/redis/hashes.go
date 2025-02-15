package redis

import (
	"context"
	"reflect"

	pkgerrors "github.com/pkg/errors"
)

// HashSet accepts values in following formats:
//
//   - HSet("myhash", map[string]interface{}{"key1": "value1", "key2": "value2"})
//
//     Playing struct With "redis" tag.
//     type MyHash struct { Key1 string `redis:"key1"`; Key2 int `redis:"key2"` }
//
//   - HSet("myhash", MyHash{"value1", "value2"}) Warn: redis-server >= 4.0
//     For struct, can be a structure pointer type, we only parse the field whose tag is redis.
//     if you don't want the field to be read, you can use the `redis:"-"` flag to ignore it,
//     or you don't need to set the redis tag.
//     For the type of structure field, we only support simple data types:
//     string, int/uint(8,16,32,64), float(32,64), time.Time(to RFC3339Nano), time.Duration(to Nanoseconds ),
//     if you are other more complex or custom data types, please implement the encoding.BinaryMarshaler interface.
func (client redisClient) HashSet(ctx context.Context, key string, value interface{}) error {
	// 1. Return error if given value input is not in Struct or map
	if reflect.TypeOf(value).Kind() != reflect.Struct && reflect.TypeOf(value).Kind() != reflect.Map &&
		reflect.TypeOf(value).Kind() != reflect.Slice {
		return ErrUnsupportedInputType
	}

	// 2. Set values to redis
	v, err := client.rdb.HSet(ctx, key, value).Result()
	if err != nil {
		return pkgerrors.WithStack(err)
	}

	// 3. Return error if no any field set
	if v == 0 {
		return ErrFailToSetValue
	}

	return nil
}

// HashGetAll retrieves all fields and values of a Redis hash by the specified key and populates the provided pointer.
// Returns an error if the provided output is not a pointer or if any operation in Redis fails.
func (client redisClient) HashGetAll(ctx context.Context, key string, out interface{}) error {
	// 1. Return error if given input is not pointer
	if reflect.TypeOf(out).Kind() != reflect.Ptr {
		return ErrUnsupportedInputType
	}

	// 2. Get all values from redis by key
	if err := client.rdb.HGetAll(ctx, key).Scan(out); err != nil {
		return pkgerrors.WithStack(err)
	}

	return nil
}

// HashGetField retrieves the value of a specific field in a Redis hash by key and field name and populates the output.
// Returns an error if the output is not a pointer or if any Redis operation fails.
func (client redisClient) HashGetField(ctx context.Context, key string, field string, out interface{}) error {
	// 1. Return error if given input is not pointer
	if reflect.TypeOf(out).Kind() != reflect.Ptr {
		return ErrUnsupportedInputType
	}

	// 2. Get a field from redis by key and field name
	if err := client.rdb.HGet(ctx, key, field).Scan(out); err != nil {
		return pkgerrors.WithStack(err)
	}

	return nil
}

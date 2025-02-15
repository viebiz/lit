package redis

import (
	"context"
	"errors"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/viebiz/lit/mocks/mockredis"
)

func Test_redisClient_HashSet(t *testing.T) {
	type redisClientArg struct {
		givenContext context.Context
		givenKey     string
		givenValue   interface{}
		expCmd       *redis.IntCmd
	}
	type arg struct {
		givenMockRedisClientArg func() redisClientArg
		givenContext            context.Context
		givenKey                string
		givenValue              interface{}
		givenRedisClientArg     redisClientArg
		expErr                  error
	}
	tcs := map[string]arg{
		"error: ErrUnsupportedInputType": {
			givenMockRedisClientArg: func() redisClientArg {
				var cmd redis.IntCmd
				cmd.SetErr(ErrUnsupportedInputType)
				return redisClientArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   "value",
					expCmd:       &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			givenValue:   "value",
			expErr:       ErrUnsupportedInputType,
		},
		"error: set value": {
			givenMockRedisClientArg: func() redisClientArg {
				var cmd redis.IntCmd
				cmd.SetErr(errors.New("redis: hset error"))
				return redisClientArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   map[string]string{"key1": "value1"},
					expCmd:       &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			givenValue:   map[string]string{"key1": "value1"},
			expErr:       errors.New("redis: hset error"),
		},
		"error: no field": {
			givenMockRedisClientArg: func() redisClientArg {
				var cmd redis.IntCmd
				cmd.SetVal(0)
				return redisClientArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   map[string]string{"key1": "value1"},
					expCmd:       &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			givenValue:   map[string]string{"key1": "value1"},
			expErr:       ErrFailToSetValue,
		},
		"success": {
			givenMockRedisClientArg: func() redisClientArg {
				var cmd redis.IntCmd
				cmd.SetVal(1)
				return redisClientArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   map[string]string{"key1": "value1"},
					expCmd:       &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			givenValue:   map[string]string{"key1": "value1"},
		},
	}
	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()
			// Mocks
			mockRedisClient := new(mockredis.MockUniversalClient)

			// Given
			tc.givenRedisClientArg = tc.givenMockRedisClientArg()
			mockRedisClient.ExpectedCalls = []*mock.Call{
				mockRedisClient.On("HSet",
					tc.givenRedisClientArg.givenContext,
					tc.givenRedisClientArg.givenKey,
					tc.givenRedisClientArg.givenValue,
				).Return(tc.givenRedisClientArg.expCmd),
			}

			// When
			instance := redisClient{
				rdb: mockRedisClient,
			}
			err := instance.HashSet(tc.givenContext, tc.givenKey, tc.givenValue)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_redisClient_HashGetAll(t *testing.T) {
	type redisClientArg struct {
		givenContext context.Context
		givenKey     string
		expCmd       *redis.MapStringStringCmd
	}
	type arg struct {
		givenMockRedisClientArg func() redisClientArg
		givenContext            context.Context
		givenKey                string
		givenValue              interface{}
		givenRedisClientArg     redisClientArg
		expErr                  error
	}
	tcs := map[string]arg{
		"error: ErrUnsupportedInputType": {
			givenMockRedisClientArg: func() redisClientArg {
				var cmd redis.MapStringStringCmd
				cmd.SetErr(ErrUnsupportedInputType)
				return redisClientArg{
					givenContext: context.Background(),
					givenKey:     "key",
					expCmd:       &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			givenValue:   "value",
			expErr:       ErrUnsupportedInputType,
		},
		"error": {
			givenMockRedisClientArg: func() redisClientArg {
				var cmd redis.MapStringStringCmd
				cmd.SetErr(errors.New("redis: hash get all error"))
				return redisClientArg{
					givenContext: context.Background(),
					givenKey:     "key",
					expCmd:       &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			givenValue:   &map[string]string{},
			expErr:       errors.New("redis: hash get all error"),
		},
		"success": {
			givenMockRedisClientArg: func() redisClientArg {
				var cmd redis.MapStringStringCmd
				//cmd.SetVal(map[string]string{"key1": "value1"})
				return redisClientArg{
					givenContext: context.Background(),
					givenKey:     "key",
					expCmd:       &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			givenValue:   &arg{},
		},
	}
	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()
			// Mocks
			mockRedisClient := new(mockredis.MockUniversalClient)

			// Given
			tc.givenRedisClientArg = tc.givenMockRedisClientArg()
			mockRedisClient.ExpectedCalls = []*mock.Call{
				mockRedisClient.On("HGetAll",
					tc.givenRedisClientArg.givenContext,
					tc.givenRedisClientArg.givenKey,
				).Return(tc.givenRedisClientArg.expCmd),
			}

			// When
			instance := redisClient{
				rdb: mockRedisClient,
			}
			err := instance.HashGetAll(tc.givenContext, tc.givenKey, tc.givenValue)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_redisClient_HashGetField(t *testing.T) {
	type redisClientArg struct {
		givenContext context.Context
		givenKey     string
		givenField   string
		expCmd       *redis.StringCmd
	}
	type arg struct {
		givenMockRedisClientArg func() redisClientArg
		givenContext            context.Context
		givenKey                string
		givenField              string
		givenOut                interface{}
		givenRedisClientArg     redisClientArg
		expErr                  error
	}
	tcs := map[string]arg{
		"error: ErrUnsupportedInputType": {
			givenMockRedisClientArg: func() redisClientArg {
				var cmd redis.StringCmd
				cmd.SetErr(ErrUnsupportedInputType)
				return redisClientArg{
					givenContext: context.Background(),
					givenKey:     "key",
					expCmd:       &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			givenField:   "givenField",
			givenOut:     "",
			expErr:       ErrUnsupportedInputType,
		},
		"error": {
			givenMockRedisClientArg: func() redisClientArg {
				var cmd redis.StringCmd
				cmd.SetErr(errors.New("redis: hash get error"))
				return redisClientArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenField:   "givenField",
					expCmd:       &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			givenField:   "givenField",
			givenOut:     StringPtr(""),
			expErr:       errors.New("redis: hash get error"),
		},
		"success": {
			givenMockRedisClientArg: func() redisClientArg {
				var cmd redis.StringCmd
				cmd.SetVal("value")
				return redisClientArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenField:   "givenField",
					expCmd:       &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			givenField:   "givenField",
			givenOut:     StringPtr(""),
		},
	}
	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()
			// Mocks
			mockRedisClient := new(mockredis.MockUniversalClient)

			// Given
			tc.givenRedisClientArg = tc.givenMockRedisClientArg()
			mockRedisClient.ExpectedCalls = []*mock.Call{
				mockRedisClient.On("HGet",
					tc.givenRedisClientArg.givenContext,
					tc.givenRedisClientArg.givenKey,
					tc.givenRedisClientArg.givenField,
				).Return(tc.givenRedisClientArg.expCmd),
			}

			// When
			instance := redisClient{
				rdb: mockRedisClient,
			}
			err := instance.HashGetField(tc.givenContext, tc.givenKey, tc.givenField, tc.givenOut)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

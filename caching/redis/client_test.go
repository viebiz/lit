package redis

import (
	"context"
	"crypto/tls"
	"errors"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/viebiz/lit/mocks/mockredis"
)

func TestNewClient(t *testing.T) {
	type arg struct {
		givenURL string
		expErr   error
	}
	tcs := map[string]arg{
		"error": {
			givenURL: "",
			expErr:   errors.New("redis: invalid URL scheme: "),
		},
		"success": {
			givenURL: "redis://localhost:6379/1",
		},
	}
	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()
			// Given

			// When
			instance, err := NewClient(tc.givenURL)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
				require.NotNil(t, instance)
			}
		})
	}
}

func TestNewClientWithTLS(t *testing.T) {
	type arg struct {
		givenURL    string
		givenConfig *tls.Config
		expErr      error
	}
	tcs := map[string]arg{
		"error": {
			givenURL: "",
			expErr:   errors.New("redis: invalid URL scheme: "),
		},
		"success": {
			givenURL: "redis://localhost:6379/1",
		},
	}
	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()
			// Given

			// When
			instance, err := NewClientWithTLS(tc.givenURL, tc.givenConfig)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
				require.NotNil(t, instance)
			}
		})
	}
}

func Test_redisClient_Ping(t *testing.T) {
	type redisClientArg struct {
		givenContext context.Context
		expCmd       *redis.StatusCmd
	}
	type arg struct {
		givenMockRedisClientArg func() redisClientArg
		givenContext            context.Context
		givenRedisClientArg     redisClientArg
		expErr                  error
	}
	tcs := map[string]arg{
		"error": {
			givenMockRedisClientArg: func() redisClientArg {
				var cmd redis.StatusCmd
				cmd.SetErr(errors.New("redis: ping error"))
				return redisClientArg{
					givenContext: context.Background(),
					expCmd:       &cmd,
				}
			},
			givenContext: context.Background(),
			expErr:       errors.New("redis: ping error"),
		},
		"success": {
			givenMockRedisClientArg: func() redisClientArg {
				return redisClientArg{
					givenContext: context.Background(),
					expCmd:       &redis.StatusCmd{},
				}
			},
			givenContext: context.Background(),
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
				mockRedisClient.On("Ping", tc.givenRedisClientArg.givenContext).Return(tc.givenRedisClientArg.expCmd),
			}

			// When
			instance := redisClient{
				rdb: mockRedisClient,
			}
			err := instance.Ping(tc.givenContext)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_redisClient_Do(t *testing.T) {
	type redisClientArg struct {
		givenContext context.Context
		givenCmd     string
		expCmd       *redis.Cmd
	}
	type arg struct {
		givenMockRedisClientArg func() redisClientArg
		givenContext            context.Context
		givenCmd                string
		givenRedisClientArg     redisClientArg
		expResult               interface{}
		expErr                  error
	}
	tcs := map[string]arg{
		"error": {
			givenMockRedisClientArg: func() redisClientArg {
				var cmd redis.Cmd
				cmd.SetErr(errors.New("redis: do cmd error"))
				return redisClientArg{
					givenContext: context.Background(),
					expCmd:       &cmd,
				}
			},
			givenContext: context.Background(),
			expErr:       errors.New("redis: do cmd error"),
		},
		"success": {
			givenMockRedisClientArg: func() redisClientArg {
				cmd := redis.Cmd{}
				cmd.SetVal("ok")
				return redisClientArg{
					givenContext: context.Background(),
					expCmd:       &cmd,
				}
			},
			givenContext: context.Background(),
			expResult:    "ok",
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
				mockRedisClient.On("Do", tc.givenRedisClientArg.givenContext, tc.givenRedisClientArg.givenCmd).Return(tc.givenRedisClientArg.expCmd),
			}

			// When
			instance := redisClient{
				rdb: mockRedisClient,
			}
			result, err := instance.Do(tc.givenContext, tc.givenCmd)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expResult, result)
			}
		})
	}
}

func Test_redisClient_DoInBatch(t *testing.T) {
	type redisClientArg struct {
		givenContext context.Context
		expCmd       []redis.Cmder
		expErr       error
	}
	type arg struct {
		givenMockRedisClientArg func() redisClientArg
		givenContext            context.Context
		givenCmdFn              func(cmder Commander) error
		givenRedisClientArg     redisClientArg
		expResult               interface{}
		expErr                  error
	}
	tcs := map[string]arg{
		"error": {
			givenMockRedisClientArg: func() redisClientArg {
				var cmd []redis.Cmder
				return redisClientArg{
					givenContext: context.Background(),
					expCmd:       cmd,
					expErr:       errors.New("redis: do in batch error"),
				}
			},
			givenContext: context.Background(),
			givenCmdFn: func(cmder Commander) error {
				return errors.New("redis: do in batch error")
			},
			expErr: errors.New("redis: do in batch error"),
		},
		"success": {
			givenMockRedisClientArg: func() redisClientArg {
				var cmd []redis.Cmder
				return redisClientArg{
					givenContext: context.Background(),
					expCmd:       cmd,
				}
			},
			givenContext: context.Background(),
			givenCmdFn: func(cmder Commander) error {
				return nil
			},
		},
	}
	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			// Mocks
			mockPipeliner := new(mockredis.MockPipeliner)
			mockRedisClient := new(mockredis.MockUniversalClient)

			// Given
			tc.givenRedisClientArg = tc.givenMockRedisClientArg()
			mockRedisClient.ExpectedCalls = []*mock.Call{
				mockRedisClient.On(
					"Pipelined",
					tc.givenRedisClientArg.givenContext,
					mock.Anything,
				).Return(tc.givenRedisClientArg.expCmd, tc.givenRedisClientArg.expErr).Run(func(args mock.Arguments) {
					callback := args.Get(1).(func(redis.Pipeliner) error)
					_ = callback(mockPipeliner)
				}),
			}

			// When
			instance := redisClient{
				rdb: mockRedisClient,
			}
			err := instance.DoInBatch(tc.givenContext, tc.givenCmdFn)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_redisClient_Close(t *testing.T) {
	type redisClientArg struct {
		expErr error
	}
	type arg struct {
		givenMockRedisClientArg func() redisClientArg
		givenRedisClientArg     redisClientArg
		expResult               interface{}
		expErr                  error
	}
	tcs := map[string]arg{
		"error": {
			givenMockRedisClientArg: func() redisClientArg {
				return redisClientArg{
					expErr: errors.New("redis: close error"),
				}
			},
			expErr: errors.New("redis: close error"),
		},
		"success": {
			givenMockRedisClientArg: func() redisClientArg {
				return redisClientArg{}
			},
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
				mockRedisClient.On("Close").Return(tc.givenRedisClientArg.expErr),
			}

			// When
			instance := redisClient{
				rdb: mockRedisClient,
			}
			err := instance.Close()

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_redisClient_Delete(t *testing.T) {
	type redisClientArg struct {
		givenContext context.Context
		givenKey     string
		expCmd       *redis.IntCmd
	}
	type arg struct {
		givenMockRedisClientArg func() redisClientArg
		givenContext            context.Context
		givenKey                string
		givenRedisClientArg     redisClientArg
		expResult               int64
		expErr                  error
	}
	tcs := map[string]arg{
		"error": {
			givenMockRedisClientArg: func() redisClientArg {
				var cmd redis.IntCmd
				cmd.SetErr(errors.New("redis: delete error"))
				return redisClientArg{
					givenContext: context.Background(),
					givenKey:     "redis",
					expCmd:       &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "redis",
			expErr:       errors.New("redis: delete error"),
		},
		"success": {
			givenMockRedisClientArg: func() redisClientArg {
				var cmd redis.IntCmd
				cmd.SetVal(1)
				return redisClientArg{
					givenContext: context.Background(),
					givenKey:     "redis",
					expCmd:       &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "redis",
			expResult:    1,
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
				mockRedisClient.On("Del", tc.givenRedisClientArg.givenContext, tc.givenRedisClientArg.givenKey).Return(tc.givenRedisClientArg.expCmd),
			}

			// When
			instance := redisClient{
				rdb: mockRedisClient,
			}

			result, err := instance.Delete(tc.givenContext, tc.givenKey)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expResult, result)
			}
		})
	}
}

func Test_redisClient_Expire(t *testing.T) {
	type redisClientArg struct {
		givenContext context.Context
		givenKey     string
		givenExpiry  time.Duration
		expCmd       *redis.BoolCmd
	}
	type arg struct {
		givenMockRedisClientArg func() redisClientArg
		givenContext            context.Context
		givenKey                string
		givenExpiry             time.Duration
		givenRedisClientArg     redisClientArg
		expResult               bool
		expErr                  error
	}
	tcs := map[string]arg{
		"error": {
			givenMockRedisClientArg: func() redisClientArg {
				var cmd redis.BoolCmd
				cmd.SetErr(errors.New("redis: expire error"))
				return redisClientArg{
					givenContext: context.Background(),
					givenKey:     "redis",
					expCmd:       &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "redis",
			expErr:       errors.New("redis: expire error"),
		},
		"success": {
			givenMockRedisClientArg: func() redisClientArg {
				var cmd redis.BoolCmd
				cmd.SetVal(true)
				return redisClientArg{
					givenContext: context.Background(),
					givenKey:     "redis",
					expCmd:       &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "redis",
			expResult:    true,
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
				mockRedisClient.On(
					"Expire",
					tc.givenRedisClientArg.givenContext,
					tc.givenRedisClientArg.givenKey,
					tc.givenRedisClientArg.givenExpiry,
				).Return(tc.givenRedisClientArg.expCmd),
			}

			// When
			instance := redisClient{
				rdb: mockRedisClient,
			}

			result, err := instance.Expire(tc.givenContext, tc.givenKey, tc.givenExpiry)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expResult, result)
			}
		})
	}
}

package redis

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/viebiz/lit/mocks/mockredis"
	"github.com/viebiz/redis"
)

func Test_redisClient_SetString(t *testing.T) {
	type mockCmdArg struct {
		givenContext context.Context
		givenKey     string
		givenValue   interface{}
		givenArgs    redis.SetArgs
		expCmd       *redis.StatusCmd
	}

	type arg[T Type] struct {
		givenMockCmdArgFn func() mockCmdArg
		givenMockCmdArg   mockCmdArg
		givenContext      context.Context
		givenKey          string
		givenValue        T
		givenExpiration   time.Duration
		expErr            error
	}

	tcs := map[string]arg[string]{
		"error: no expiry": {
			givenMockCmdArgFn: func() mockCmdArg {
				var cmd redis.StatusCmd
				cmd.SetErr(errors.New("error: no expiry"))
				return mockCmdArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   "value",
					givenArgs: redis.SetArgs{
						KeepTTL: true,
					},
					expCmd: &cmd,
				}
			},
			givenContext:    context.Background(),
			givenKey:        "key",
			givenValue:      "value",
			givenExpiration: 0,
			expErr:          errors.New("error: no expiry"),
		},
		"error": {
			givenMockCmdArgFn: func() mockCmdArg {
				var cmd redis.StatusCmd
				cmd.SetErr(errors.New("error"))
				return mockCmdArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   "value",
					givenArgs: redis.SetArgs{
						KeepTTL: false,
						TTL:     time.Duration(1),
					},
					expCmd: &cmd,
				}
			},
			givenContext:    context.Background(),
			givenKey:        "key",
			givenValue:      "value",
			givenExpiration: time.Duration(1),
			expErr:          errors.New("error"),
		},
		"success": {
			givenMockCmdArgFn: func() mockCmdArg {
				var cmd redis.StatusCmd
				cmd.SetVal(statusOK)
				return mockCmdArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   "value",
					givenArgs: redis.SetArgs{
						KeepTTL: false,
						TTL:     time.Duration(1),
					},
					expCmd: &cmd,
				}
			},
			givenContext:    context.Background(),
			givenKey:        "key",
			givenValue:      "value",
			givenExpiration: time.Duration(1),
		},
	}

	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			// Given
			mockUniversalClient := new(mockredis.MockUniversalClient)

			// Mock
			tc.givenMockCmdArg = tc.givenMockCmdArgFn()
			mockUniversalClient.ExpectedCalls = []*mock.Call{
				mockUniversalClient.On(
					"SetArgs",
					tc.givenMockCmdArg.givenContext,
					tc.givenMockCmdArg.givenKey,
					tc.givenMockCmdArg.givenValue,
					tc.givenMockCmdArg.givenArgs,
				).Return(tc.givenMockCmdArg.expCmd),
			}

			// When
			instance := redisClient{
				rdb: mockUniversalClient,
			}
			err := instance.SetString(tc.givenContext, tc.givenKey, tc.givenValue, tc.givenExpiration)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_redisClient_SetStringIfNotExist(t *testing.T) {
	type mockCmdArg struct {
		givenContext context.Context
		givenKey     string
		givenValue   interface{}
		givenArgs    redis.SetArgs
		expCmd       *redis.StatusCmd
	}

	type arg[T Type] struct {
		givenMockCmdArgFn func() mockCmdArg
		givenMockCmdArg   mockCmdArg
		givenContext      context.Context
		givenKey          string
		givenValue        T
		givenExpiration   time.Duration
		expErr            error
	}

	tcs := map[string]arg[string]{
		"error: no expiry": {
			givenMockCmdArgFn: func() mockCmdArg {
				var cmd redis.StatusCmd
				cmd.SetErr(errors.New("error: no expiry"))
				return mockCmdArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   "value",
					givenArgs: redis.SetArgs{
						KeepTTL: true,
						Mode:    setModeNX.String(),
					},
					expCmd: &cmd,
				}
			},
			givenContext:    context.Background(),
			givenKey:        "key",
			givenValue:      "value",
			givenExpiration: 0,
			expErr:          errors.New("error: no expiry"),
		},
		"error": {
			givenMockCmdArgFn: func() mockCmdArg {
				var cmd redis.StatusCmd
				cmd.SetErr(errors.New("error"))
				return mockCmdArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   "value",
					givenArgs: redis.SetArgs{
						KeepTTL: false,
						TTL:     time.Duration(1),
						Mode:    setModeNX.String(),
					},
					expCmd: &cmd,
				}
			},
			givenContext:    context.Background(),
			givenKey:        "key",
			givenValue:      "value",
			givenExpiration: time.Duration(1),
			expErr:          errors.New("error"),
		},
		"error: ErrFailToSetValue": {
			givenMockCmdArgFn: func() mockCmdArg {
				var cmd redis.StatusCmd
				cmd.SetVal("ok")
				return mockCmdArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   "value",
					givenArgs: redis.SetArgs{
						KeepTTL: false,
						TTL:     time.Duration(1),
						Mode:    setModeNX.String(),
					},
					expCmd: &cmd,
				}
			},
			givenContext:    context.Background(),
			givenKey:        "key",
			givenValue:      "value",
			givenExpiration: time.Duration(1),
			expErr:          ErrFailToSetValue,
		},
		"success": {
			givenMockCmdArgFn: func() mockCmdArg {
				var cmd redis.StatusCmd
				cmd.SetVal(statusOK)
				return mockCmdArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   "value",
					givenArgs: redis.SetArgs{
						KeepTTL: false,
						TTL:     time.Duration(1),
						Mode:    setModeNX.String(),
					},
					expCmd: &cmd,
				}
			},
			givenContext:    context.Background(),
			givenKey:        "key",
			givenValue:      "value",
			givenExpiration: time.Duration(1),
		},
	}

	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			// Given
			mockUniversalClient := new(mockredis.MockUniversalClient)

			// Mock
			tc.givenMockCmdArg = tc.givenMockCmdArgFn()
			mockUniversalClient.ExpectedCalls = []*mock.Call{
				mockUniversalClient.On(
					"SetArgs",
					tc.givenMockCmdArg.givenContext,
					tc.givenMockCmdArg.givenKey,
					tc.givenMockCmdArg.givenValue,
					tc.givenMockCmdArg.givenArgs,
				).Return(tc.givenMockCmdArg.expCmd),
			}

			// When
			instance := redisClient{
				rdb: mockUniversalClient,
			}
			err := instance.SetStringIfNotExist(tc.givenContext, tc.givenKey, tc.givenValue, tc.givenExpiration)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_redisClient_SetStringIfExist(t *testing.T) {
	type mockCmdArg struct {
		givenContext context.Context
		givenKey     string
		givenValue   interface{}
		givenArgs    redis.SetArgs
		expCmd       *redis.StatusCmd
	}

	type arg[T Type] struct {
		givenMockCmdArgFn func() mockCmdArg
		givenMockCmdArg   mockCmdArg
		givenContext      context.Context
		givenKey          string
		givenValue        T
		givenExpiration   time.Duration
		expErr            error
	}

	tcs := map[string]arg[string]{
		"error: no expiry": {
			givenMockCmdArgFn: func() mockCmdArg {
				var cmd redis.StatusCmd
				cmd.SetErr(errors.New("error: no expiry"))
				return mockCmdArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   "value",
					givenArgs: redis.SetArgs{
						KeepTTL: true,
						Mode:    setModeXX.String(),
					},
					expCmd: &cmd,
				}
			},
			givenContext:    context.Background(),
			givenKey:        "key",
			givenValue:      "value",
			givenExpiration: 0,
			expErr:          errors.New("error: no expiry"),
		},
		"error": {
			givenMockCmdArgFn: func() mockCmdArg {
				var cmd redis.StatusCmd
				cmd.SetErr(errors.New("error"))
				return mockCmdArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   "value",
					givenArgs: redis.SetArgs{
						KeepTTL: false,
						TTL:     time.Duration(1),
						Mode:    setModeXX.String(),
					},
					expCmd: &cmd,
				}
			},
			givenContext:    context.Background(),
			givenKey:        "key",
			givenValue:      "value",
			givenExpiration: time.Duration(1),
			expErr:          errors.New("error"),
		},
		"error: ErrFailToSetValue": {
			givenMockCmdArgFn: func() mockCmdArg {
				var cmd redis.StatusCmd
				cmd.SetVal("ok")
				return mockCmdArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   "value",
					givenArgs: redis.SetArgs{
						KeepTTL: false,
						TTL:     time.Duration(1),
						Mode:    setModeXX.String(),
					},
					expCmd: &cmd,
				}
			},
			givenContext:    context.Background(),
			givenKey:        "key",
			givenValue:      "value",
			givenExpiration: time.Duration(1),
			expErr:          ErrFailToSetValue,
		},
		"success": {
			givenMockCmdArgFn: func() mockCmdArg {
				var cmd redis.StatusCmd
				cmd.SetVal(statusOK)
				return mockCmdArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   "value",
					givenArgs: redis.SetArgs{
						KeepTTL: false,
						TTL:     time.Duration(1),
						Mode:    setModeXX.String(),
					},
					expCmd: &cmd,
				}
			},
			givenContext:    context.Background(),
			givenKey:        "key",
			givenValue:      "value",
			givenExpiration: time.Duration(1),
		},
	}

	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			// Given
			mockUniversalClient := new(mockredis.MockUniversalClient)

			// Mock
			tc.givenMockCmdArg = tc.givenMockCmdArgFn()
			mockUniversalClient.ExpectedCalls = []*mock.Call{
				mockUniversalClient.On(
					"SetArgs",
					tc.givenMockCmdArg.givenContext,
					tc.givenMockCmdArg.givenKey,
					tc.givenMockCmdArg.givenValue,
					tc.givenMockCmdArg.givenArgs,
				).Return(tc.givenMockCmdArg.expCmd),
			}

			// When
			instance := redisClient{
				rdb: mockUniversalClient,
			}
			err := instance.SetStringIfExist(tc.givenContext, tc.givenKey, tc.givenValue, tc.givenExpiration)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_redisClient_GetString(t *testing.T) {
	type mockCmdArg struct {
		givenContext context.Context
		givenKey     string
		expCmd       *redis.StringCmd
	}

	type arg[T Type] struct {
		givenMockCmdArgFn func() mockCmdArg
		givenMockCmdArg   mockCmdArg
		givenContext      context.Context
		givenKey          string
		expResult         T
		expErr            error
	}

	tcs := map[string]arg[string]{
		"nil": {
			givenMockCmdArgFn: func() mockCmdArg {
				var cmd redis.StringCmd
				cmd.SetErr(redis.Nil)
				return mockCmdArg{
					givenContext: context.Background(),
					givenKey:     "key",
					expCmd:       &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
		},
		"error": {
			givenMockCmdArgFn: func() mockCmdArg {
				var cmd redis.StringCmd
				cmd.SetErr(errors.New("error"))
				return mockCmdArg{
					givenContext: context.Background(),
					givenKey:     "key",
					expCmd:       &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			expErr:       errors.New("error"),
		},
		"success": {
			givenMockCmdArgFn: func() mockCmdArg {
				var cmd redis.StringCmd
				cmd.SetVal("value")
				return mockCmdArg{
					givenContext: context.Background(),
					givenKey:     "key",
					expCmd:       &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			expResult:    "value",
		},
	}

	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			// Given
			mockUniversalClient := new(mockredis.MockUniversalClient)

			// Mock
			tc.givenMockCmdArg = tc.givenMockCmdArgFn()
			mockUniversalClient.ExpectedCalls = []*mock.Call{
				mockUniversalClient.On(
					"Get",
					tc.givenMockCmdArg.givenContext,
					tc.givenMockCmdArg.givenKey,
				).Return(tc.givenMockCmdArg.expCmd),
			}

			// When
			instance := redisClient{
				rdb: mockUniversalClient,
			}
			result, err := instance.GetString(tc.givenContext, tc.givenKey)

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

func Test_redisClient_SetInt(t *testing.T) {
	type mockCmdArg struct {
		givenContext context.Context
		givenKey     string
		givenValue   interface{}
		givenArgs    redis.SetArgs
		expCmd       *redis.StatusCmd
	}

	type arg[T Type] struct {
		givenMockCmdArgFn func() mockCmdArg
		givenMockCmdArg   mockCmdArg
		givenContext      context.Context
		givenKey          string
		givenValue        T
		givenExpiration   time.Duration
		expErr            error
	}

	tcs := map[string]arg[int64]{
		"error: no expiry": {
			givenMockCmdArgFn: func() mockCmdArg {
				var cmd redis.StatusCmd
				cmd.SetErr(errors.New("error: no expiry"))
				return mockCmdArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   int64(0),
					givenArgs: redis.SetArgs{
						KeepTTL: true,
					},
					expCmd: &cmd,
				}
			},
			givenContext:    context.Background(),
			givenKey:        "key",
			givenValue:      0,
			givenExpiration: 0,
			expErr:          errors.New("error: no expiry"),
		},
		"error": {
			givenMockCmdArgFn: func() mockCmdArg {
				var cmd redis.StatusCmd
				cmd.SetErr(errors.New("error"))
				return mockCmdArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   int64(0),
					givenArgs: redis.SetArgs{
						KeepTTL: false,
						TTL:     time.Duration(1),
					},
					expCmd: &cmd,
				}
			},
			givenContext:    context.Background(),
			givenKey:        "key",
			givenValue:      0,
			givenExpiration: time.Duration(1),
			expErr:          errors.New("error"),
		},
		"success": {
			givenMockCmdArgFn: func() mockCmdArg {
				var cmd redis.StatusCmd
				cmd.SetVal(statusOK)
				return mockCmdArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   int64(1),
					givenArgs: redis.SetArgs{
						KeepTTL: false,
						TTL:     time.Duration(1),
					},
					expCmd: &cmd,
				}
			},
			givenContext:    context.Background(),
			givenKey:        "key",
			givenValue:      1,
			givenExpiration: time.Duration(1),
		},
	}

	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			// Given
			mockUniversalClient := new(mockredis.MockUniversalClient)

			// Mock
			tc.givenMockCmdArg = tc.givenMockCmdArgFn()
			mockUniversalClient.ExpectedCalls = []*mock.Call{
				mockUniversalClient.On(
					"SetArgs",
					tc.givenMockCmdArg.givenContext,
					tc.givenMockCmdArg.givenKey,
					tc.givenMockCmdArg.givenValue,
					tc.givenMockCmdArg.givenArgs,
				).Return(tc.givenMockCmdArg.expCmd),
			}

			// When
			instance := redisClient{
				rdb: mockUniversalClient,
			}
			err := instance.SetInt(tc.givenContext, tc.givenKey, tc.givenValue, tc.givenExpiration)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_redisClient_SetIntIfNotExist(t *testing.T) {
	type mockCmdArg struct {
		givenContext context.Context
		givenKey     string
		givenValue   interface{}
		givenArgs    redis.SetArgs
		expCmd       *redis.StatusCmd
	}

	type arg[T Type] struct {
		givenMockCmdArgFn func() mockCmdArg
		givenMockCmdArg   mockCmdArg
		givenContext      context.Context
		givenKey          string
		givenValue        T
		givenExpiration   time.Duration
		expErr            error
	}

	tcs := map[string]arg[int64]{
		"error: no expiry": {
			givenMockCmdArgFn: func() mockCmdArg {
				var cmd redis.StatusCmd
				cmd.SetErr(errors.New("error: no expiry"))
				return mockCmdArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   int64(0),
					givenArgs: redis.SetArgs{
						KeepTTL: true,
						Mode:    setModeNX.String(),
					},
					expCmd: &cmd,
				}
			},
			givenContext:    context.Background(),
			givenKey:        "key",
			givenValue:      int64(0),
			givenExpiration: 0,
			expErr:          errors.New("error: no expiry"),
		},
		"error": {
			givenMockCmdArgFn: func() mockCmdArg {
				var cmd redis.StatusCmd
				cmd.SetErr(errors.New("error"))
				return mockCmdArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   int64(0),
					givenArgs: redis.SetArgs{
						KeepTTL: false,
						TTL:     time.Duration(1),
						Mode:    setModeNX.String(),
					},
					expCmd: &cmd,
				}
			},
			givenContext:    context.Background(),
			givenKey:        "key",
			givenValue:      0,
			givenExpiration: time.Duration(1),
			expErr:          errors.New("error"),
		},
		"error: ErrFailToSetValue": {
			givenMockCmdArgFn: func() mockCmdArg {
				var cmd redis.StatusCmd
				cmd.SetVal("ok")
				return mockCmdArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   int64(0),
					givenArgs: redis.SetArgs{
						KeepTTL: false,
						TTL:     time.Duration(1),
						Mode:    setModeNX.String(),
					},
					expCmd: &cmd,
				}
			},
			givenContext:    context.Background(),
			givenKey:        "key",
			givenValue:      0,
			givenExpiration: time.Duration(1),
			expErr:          ErrFailToSetValue,
		},
		"success": {
			givenMockCmdArgFn: func() mockCmdArg {
				var cmd redis.StatusCmd
				cmd.SetVal(statusOK)
				return mockCmdArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   int64(1),
					givenArgs: redis.SetArgs{
						KeepTTL: false,
						TTL:     time.Duration(1),
						Mode:    setModeNX.String(),
					},
					expCmd: &cmd,
				}
			},
			givenContext:    context.Background(),
			givenKey:        "key",
			givenValue:      1,
			givenExpiration: time.Duration(1),
		},
	}

	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			// Given
			mockUniversalClient := new(mockredis.MockUniversalClient)

			// Mock
			tc.givenMockCmdArg = tc.givenMockCmdArgFn()
			mockUniversalClient.ExpectedCalls = []*mock.Call{
				mockUniversalClient.On(
					"SetArgs",
					tc.givenMockCmdArg.givenContext,
					tc.givenMockCmdArg.givenKey,
					tc.givenMockCmdArg.givenValue,
					tc.givenMockCmdArg.givenArgs,
				).Return(tc.givenMockCmdArg.expCmd),
			}

			// When
			instance := redisClient{
				rdb: mockUniversalClient,
			}
			err := instance.SetIntIfNotExist(tc.givenContext, tc.givenKey, tc.givenValue, tc.givenExpiration)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_redisClient_SetIntIfExist(t *testing.T) {
	type mockCmdArg struct {
		givenContext context.Context
		givenKey     string
		givenValue   interface{}
		givenArgs    redis.SetArgs
		expCmd       *redis.StatusCmd
	}

	type arg[T Type] struct {
		givenMockCmdArgFn func() mockCmdArg
		givenMockCmdArg   mockCmdArg
		givenContext      context.Context
		givenKey          string
		givenValue        T
		givenExpiration   time.Duration
		expErr            error
	}

	tcs := map[string]arg[int64]{
		"error: no expiry": {
			givenMockCmdArgFn: func() mockCmdArg {
				var cmd redis.StatusCmd
				cmd.SetErr(errors.New("error: no expiry"))
				return mockCmdArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   int64(0),
					givenArgs: redis.SetArgs{
						KeepTTL: true,
						Mode:    setModeXX.String(),
					},
					expCmd: &cmd,
				}
			},
			givenContext:    context.Background(),
			givenKey:        "key",
			givenValue:      0,
			givenExpiration: 0,
			expErr:          errors.New("error: no expiry"),
		},
		"error": {
			givenMockCmdArgFn: func() mockCmdArg {
				var cmd redis.StatusCmd
				cmd.SetErr(errors.New("error"))
				return mockCmdArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   int64(0),
					givenArgs: redis.SetArgs{
						KeepTTL: false,
						TTL:     time.Duration(1),
						Mode:    setModeXX.String(),
					},
					expCmd: &cmd,
				}
			},
			givenContext:    context.Background(),
			givenKey:        "key",
			givenValue:      0,
			givenExpiration: time.Duration(1),
			expErr:          errors.New("error"),
		},
		"error: ErrFailToSetValue": {
			givenMockCmdArgFn: func() mockCmdArg {
				var cmd redis.StatusCmd
				cmd.SetVal("ok")
				return mockCmdArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   int64(0),
					givenArgs: redis.SetArgs{
						KeepTTL: false,
						TTL:     time.Duration(1),
						Mode:    setModeXX.String(),
					},
					expCmd: &cmd,
				}
			},
			givenContext:    context.Background(),
			givenKey:        "key",
			givenValue:      0,
			givenExpiration: time.Duration(1),
			expErr:          ErrFailToSetValue,
		},
		"success": {
			givenMockCmdArgFn: func() mockCmdArg {
				var cmd redis.StatusCmd
				cmd.SetVal(statusOK)
				return mockCmdArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   int64(1),
					givenArgs: redis.SetArgs{
						KeepTTL: false,
						TTL:     time.Duration(1),
						Mode:    setModeXX.String(),
					},
					expCmd: &cmd,
				}
			},
			givenContext:    context.Background(),
			givenKey:        "key",
			givenValue:      1,
			givenExpiration: time.Duration(1),
		},
	}

	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			// Given
			mockUniversalClient := new(mockredis.MockUniversalClient)

			// Mock
			tc.givenMockCmdArg = tc.givenMockCmdArgFn()
			mockUniversalClient.ExpectedCalls = []*mock.Call{
				mockUniversalClient.On(
					"SetArgs",
					tc.givenMockCmdArg.givenContext,
					tc.givenMockCmdArg.givenKey,
					tc.givenMockCmdArg.givenValue,
					tc.givenMockCmdArg.givenArgs,
				).Return(tc.givenMockCmdArg.expCmd),
			}

			// When
			instance := redisClient{
				rdb: mockUniversalClient,
			}
			err := instance.SetIntIfExist(tc.givenContext, tc.givenKey, tc.givenValue, tc.givenExpiration)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_redisClient_GetInt(t *testing.T) {
	type mockCmdArg struct {
		givenContext context.Context
		givenKey     string
		expCmd       *redis.StringCmd
	}

	type arg[T Type] struct {
		givenMockCmdArgFn func() mockCmdArg
		givenMockCmdArg   mockCmdArg
		givenContext      context.Context
		givenKey          string
		expResult         T
		expErr            error
	}

	tcs := map[string]arg[int64]{
		"nil": {
			givenMockCmdArgFn: func() mockCmdArg {
				var cmd redis.StringCmd
				cmd.SetErr(redis.Nil)
				return mockCmdArg{
					givenContext: context.Background(),
					givenKey:     "key",
					expCmd:       &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
		},
		"error": {
			givenMockCmdArgFn: func() mockCmdArg {
				var cmd redis.StringCmd
				cmd.SetErr(errors.New("error"))
				return mockCmdArg{
					givenContext: context.Background(),
					givenKey:     "key",
					expCmd:       &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			expErr:       errors.New("error"),
		},
		"success": {
			givenMockCmdArgFn: func() mockCmdArg {
				var cmd redis.StringCmd
				cmd.SetVal("1")
				return mockCmdArg{
					givenContext: context.Background(),
					givenKey:     "key",
					expCmd:       &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			expResult:    1,
		},
	}

	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			// Given
			mockUniversalClient := new(mockredis.MockUniversalClient)

			// Mock
			tc.givenMockCmdArg = tc.givenMockCmdArgFn()
			mockUniversalClient.ExpectedCalls = []*mock.Call{
				mockUniversalClient.On(
					"Get",
					tc.givenMockCmdArg.givenContext,
					tc.givenMockCmdArg.givenKey,
				).Return(tc.givenMockCmdArg.expCmd),
			}

			// When
			instance := redisClient{
				rdb: mockUniversalClient,
			}
			result, err := instance.GetInt(tc.givenContext, tc.givenKey)

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

func Test_redisClient_IncrementBy(t *testing.T) {
	type mockRedisClientArg struct {
		givenContext context.Context
		givenKey     string
		givenValue   int64
		expCmd       *redis.IntCmd
	}
	type arg struct {
		givenMockRedisClientArgFn func() mockRedisClientArg
		givenMockRedisClientArg   mockRedisClientArg
		givenContext              context.Context
		givenKey                  string
		givenValue                int64
		expResult                 int64
		expErr                    error
	}
	tcs := map[string]arg{
		"error": {
			givenMockRedisClientArgFn: func() mockRedisClientArg {
				var cmd redis.IntCmd
				cmd.SetErr(errors.New("error"))
				return mockRedisClientArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   1,
					expCmd:       &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			givenValue:   1,
			expErr:       errors.New("error"),
		},
		"success": {
			givenMockRedisClientArgFn: func() mockRedisClientArg {
				var cmd redis.IntCmd
				cmd.SetVal(1)
				return mockRedisClientArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   1,
					expCmd:       &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			givenValue:   1,
			expResult:    1,
		},
	}
	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			// Given
			mockUniversalClient := new(mockredis.MockUniversalClient)

			// Mocks
			tc.givenMockRedisClientArg = tc.givenMockRedisClientArgFn()
			mockUniversalClient.ExpectedCalls = []*mock.Call{
				mockUniversalClient.On(
					"IncrBy",
					tc.givenMockRedisClientArg.givenContext,
					tc.givenMockRedisClientArg.givenKey,
					tc.givenMockRedisClientArg.givenValue,
				).Return(tc.givenMockRedisClientArg.expCmd),
			}

			// When
			instance := redisClient{
				rdb: mockUniversalClient,
			}
			result, err := instance.IncrementBy(tc.givenContext, tc.givenKey, tc.givenValue)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, result, tc.expResult)
			}
		})
	}
}

func Test_redisClient_DecrementBy(t *testing.T) {
	type mockRedisClientArg struct {
		givenContext context.Context
		givenKey     string
		givenValue   int64
		expCmd       *redis.IntCmd
	}
	type arg struct {
		givenMockRedisClientArgFn func() mockRedisClientArg
		givenMockRedisClientArg   mockRedisClientArg
		givenContext              context.Context
		givenKey                  string
		givenValue                int64
		expResult                 int64
		expErr                    error
	}
	tcs := map[string]arg{
		"error": {
			givenMockRedisClientArgFn: func() mockRedisClientArg {
				var cmd redis.IntCmd
				cmd.SetErr(errors.New("error"))
				return mockRedisClientArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   1,
					expCmd:       &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			givenValue:   1,
			expErr:       errors.New("error"),
		},
		"success": {
			givenMockRedisClientArgFn: func() mockRedisClientArg {
				var cmd redis.IntCmd
				cmd.SetVal(1)
				return mockRedisClientArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   1,
					expCmd:       &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			givenValue:   1,
			expResult:    1,
		},
	}
	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			// Given
			mockUniversalClient := new(mockredis.MockUniversalClient)

			// Mocks
			tc.givenMockRedisClientArg = tc.givenMockRedisClientArgFn()
			mockUniversalClient.ExpectedCalls = []*mock.Call{
				mockUniversalClient.On(
					"DecrBy",
					tc.givenMockRedisClientArg.givenContext,
					tc.givenMockRedisClientArg.givenKey,
					tc.givenMockRedisClientArg.givenValue,
				).Return(tc.givenMockRedisClientArg.expCmd),
			}

			// When
			instance := redisClient{
				rdb: mockUniversalClient,
			}
			result, err := instance.DecrementBy(tc.givenContext, tc.givenKey, tc.givenValue)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, result, tc.expResult)
			}
		})
	}
}

func Test_redisClient_SetFloat(t *testing.T) {
	type mockCmdArg struct {
		givenContext context.Context
		givenKey     string
		givenValue   interface{}
		givenArgs    redis.SetArgs
		expCmd       *redis.StatusCmd
	}

	type arg[T Type] struct {
		givenMockCmdArgFn func() mockCmdArg
		givenMockCmdArg   mockCmdArg
		givenContext      context.Context
		givenKey          string
		givenValue        T
		givenExpiration   time.Duration
		expErr            error
	}

	tcs := map[string]arg[float64]{
		"error: no expiry": {
			givenMockCmdArgFn: func() mockCmdArg {
				var cmd redis.StatusCmd
				cmd.SetErr(errors.New("error: no expiry"))
				return mockCmdArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   float64(0),
					givenArgs: redis.SetArgs{
						KeepTTL: true,
					},
					expCmd: &cmd,
				}
			},
			givenContext:    context.Background(),
			givenKey:        "key",
			givenValue:      0,
			givenExpiration: 0,
			expErr:          errors.New("error: no expiry"),
		},
		"error": {
			givenMockCmdArgFn: func() mockCmdArg {
				var cmd redis.StatusCmd
				cmd.SetErr(errors.New("error"))
				return mockCmdArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   float64(0),
					givenArgs: redis.SetArgs{
						KeepTTL: false,
						TTL:     time.Duration(1),
					},
					expCmd: &cmd,
				}
			},
			givenContext:    context.Background(),
			givenKey:        "key",
			givenValue:      0,
			givenExpiration: time.Duration(1),
			expErr:          errors.New("error"),
		},
		"success": {
			givenMockCmdArgFn: func() mockCmdArg {
				var cmd redis.StatusCmd
				cmd.SetVal(statusOK)
				return mockCmdArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   float64(1),
					givenArgs: redis.SetArgs{
						KeepTTL: false,
						TTL:     time.Duration(1),
					},
					expCmd: &cmd,
				}
			},
			givenContext:    context.Background(),
			givenKey:        "key",
			givenValue:      1,
			givenExpiration: time.Duration(1),
		},
	}

	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			// Given
			mockUniversalClient := new(mockredis.MockUniversalClient)

			// Mock
			tc.givenMockCmdArg = tc.givenMockCmdArgFn()
			mockUniversalClient.ExpectedCalls = []*mock.Call{
				mockUniversalClient.On(
					"SetArgs",
					tc.givenMockCmdArg.givenContext,
					tc.givenMockCmdArg.givenKey,
					tc.givenMockCmdArg.givenValue,
					tc.givenMockCmdArg.givenArgs,
				).Return(tc.givenMockCmdArg.expCmd),
			}

			// When
			instance := redisClient{
				rdb: mockUniversalClient,
			}
			err := instance.SetFloat(tc.givenContext, tc.givenKey, tc.givenValue, tc.givenExpiration)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_redisClient_SetFloatIfNotExist(t *testing.T) {
	type mockCmdArg struct {
		givenContext context.Context
		givenKey     string
		givenValue   interface{}
		givenArgs    redis.SetArgs
		expCmd       *redis.StatusCmd
	}

	type arg[T Type] struct {
		givenMockCmdArgFn func() mockCmdArg
		givenMockCmdArg   mockCmdArg
		givenContext      context.Context
		givenKey          string
		givenValue        T
		givenExpiration   time.Duration
		expErr            error
	}

	tcs := map[string]arg[float64]{
		"error: no expiry": {
			givenMockCmdArgFn: func() mockCmdArg {
				var cmd redis.StatusCmd
				cmd.SetErr(errors.New("error: no expiry"))
				return mockCmdArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   float64(0),
					givenArgs: redis.SetArgs{
						KeepTTL: true,
						Mode:    setModeNX.String(),
					},
					expCmd: &cmd,
				}
			},
			givenContext:    context.Background(),
			givenKey:        "key",
			givenValue:      0,
			givenExpiration: 0,
			expErr:          errors.New("error: no expiry"),
		},
		"error": {
			givenMockCmdArgFn: func() mockCmdArg {
				var cmd redis.StatusCmd
				cmd.SetErr(errors.New("error"))
				return mockCmdArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   float64(0),
					givenArgs: redis.SetArgs{
						KeepTTL: false,
						TTL:     time.Duration(1),
						Mode:    setModeNX.String(),
					},
					expCmd: &cmd,
				}
			},
			givenContext:    context.Background(),
			givenKey:        "key",
			givenValue:      0,
			givenExpiration: time.Duration(1),
			expErr:          errors.New("error"),
		},
		"error: ErrFailToSetValue": {
			givenMockCmdArgFn: func() mockCmdArg {
				var cmd redis.StatusCmd
				cmd.SetVal("ok")
				return mockCmdArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   float64(0),
					givenArgs: redis.SetArgs{
						KeepTTL: false,
						TTL:     time.Duration(1),
						Mode:    setModeNX.String(),
					},
					expCmd: &cmd,
				}
			},
			givenContext:    context.Background(),
			givenKey:        "key",
			givenValue:      0,
			givenExpiration: time.Duration(1),
			expErr:          ErrFailToSetValue,
		},
		"success": {
			givenMockCmdArgFn: func() mockCmdArg {
				var cmd redis.StatusCmd
				cmd.SetVal(statusOK)
				return mockCmdArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   float64(1),
					givenArgs: redis.SetArgs{
						KeepTTL: false,
						TTL:     time.Duration(1),
						Mode:    setModeNX.String(),
					},
					expCmd: &cmd,
				}
			},
			givenContext:    context.Background(),
			givenKey:        "key",
			givenValue:      1,
			givenExpiration: time.Duration(1),
		},
	}

	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			// Given
			mockUniversalClient := new(mockredis.MockUniversalClient)

			// Mock
			tc.givenMockCmdArg = tc.givenMockCmdArgFn()
			mockUniversalClient.ExpectedCalls = []*mock.Call{
				mockUniversalClient.On(
					"SetArgs",
					tc.givenMockCmdArg.givenContext,
					tc.givenMockCmdArg.givenKey,
					tc.givenMockCmdArg.givenValue,
					tc.givenMockCmdArg.givenArgs,
				).Return(tc.givenMockCmdArg.expCmd),
			}

			// When
			instance := redisClient{
				rdb: mockUniversalClient,
			}
			err := instance.SetFloatIfNotExist(tc.givenContext, tc.givenKey, tc.givenValue, tc.givenExpiration)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_redisClient_SetFloatIfExist(t *testing.T) {
	type mockCmdArg struct {
		givenContext context.Context
		givenKey     string
		givenValue   interface{}
		givenArgs    redis.SetArgs
		expCmd       *redis.StatusCmd
	}

	type arg[T Type] struct {
		givenMockCmdArgFn func() mockCmdArg
		givenMockCmdArg   mockCmdArg
		givenContext      context.Context
		givenKey          string
		givenValue        T
		givenExpiration   time.Duration
		expErr            error
	}

	tcs := map[string]arg[float64]{
		"error: no expiry": {
			givenMockCmdArgFn: func() mockCmdArg {
				var cmd redis.StatusCmd
				cmd.SetErr(errors.New("error: no expiry"))
				return mockCmdArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   float64(0),
					givenArgs: redis.SetArgs{
						KeepTTL: true,
						Mode:    setModeXX.String(),
					},
					expCmd: &cmd,
				}
			},
			givenContext:    context.Background(),
			givenKey:        "key",
			givenValue:      0,
			givenExpiration: 0,
			expErr:          errors.New("error: no expiry"),
		},
		"error": {
			givenMockCmdArgFn: func() mockCmdArg {
				var cmd redis.StatusCmd
				cmd.SetErr(errors.New("error"))
				return mockCmdArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   float64(0),
					givenArgs: redis.SetArgs{
						KeepTTL: false,
						TTL:     time.Duration(1),
						Mode:    setModeXX.String(),
					},
					expCmd: &cmd,
				}
			},
			givenContext:    context.Background(),
			givenKey:        "key",
			givenValue:      0,
			givenExpiration: time.Duration(1),
			expErr:          errors.New("error"),
		},
		"error: ErrFailToSetValue": {
			givenMockCmdArgFn: func() mockCmdArg {
				var cmd redis.StatusCmd
				cmd.SetVal("ok")
				return mockCmdArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   float64(0),
					givenArgs: redis.SetArgs{
						KeepTTL: false,
						TTL:     time.Duration(1),
						Mode:    setModeXX.String(),
					},
					expCmd: &cmd,
				}
			},
			givenContext:    context.Background(),
			givenKey:        "key",
			givenValue:      0,
			givenExpiration: time.Duration(1),
			expErr:          ErrFailToSetValue,
		},
		"success": {
			givenMockCmdArgFn: func() mockCmdArg {
				var cmd redis.StatusCmd
				cmd.SetVal(statusOK)
				return mockCmdArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   float64(1),
					givenArgs: redis.SetArgs{
						KeepTTL: false,
						TTL:     time.Duration(1),
						Mode:    setModeXX.String(),
					},
					expCmd: &cmd,
				}
			},
			givenContext:    context.Background(),
			givenKey:        "key",
			givenValue:      1,
			givenExpiration: time.Duration(1),
		},
	}

	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			// Given
			mockUniversalClient := new(mockredis.MockUniversalClient)

			// Mock
			tc.givenMockCmdArg = tc.givenMockCmdArgFn()
			mockUniversalClient.ExpectedCalls = []*mock.Call{
				mockUniversalClient.On(
					"SetArgs",
					tc.givenMockCmdArg.givenContext,
					tc.givenMockCmdArg.givenKey,
					tc.givenMockCmdArg.givenValue,
					tc.givenMockCmdArg.givenArgs,
				).Return(tc.givenMockCmdArg.expCmd),
			}

			// When
			instance := redisClient{
				rdb: mockUniversalClient,
			}
			err := instance.SetFloatIfExist(tc.givenContext, tc.givenKey, tc.givenValue, tc.givenExpiration)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_redisClient_IncrementFloatBy(t *testing.T) {
	type mockRedisClientArg struct {
		givenContext context.Context
		givenKey     string
		givenValue   float64
		expCmd       *redis.FloatCmd
	}
	type arg struct {
		givenMockRedisClientArgFn func() mockRedisClientArg
		givenMockRedisClientArg   mockRedisClientArg
		givenContext              context.Context
		givenKey                  string
		givenValue                float64
		expResult                 float64
		expErr                    error
	}
	tcs := map[string]arg{
		"error": {
			givenMockRedisClientArgFn: func() mockRedisClientArg {
				var cmd redis.FloatCmd
				cmd.SetErr(errors.New("error"))
				return mockRedisClientArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   1,
					expCmd:       &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			givenValue:   1,
			expErr:       errors.New("error"),
		},
		"success": {
			givenMockRedisClientArgFn: func() mockRedisClientArg {
				var cmd redis.FloatCmd
				cmd.SetVal(1)
				return mockRedisClientArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   1,
					expCmd:       &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			givenValue:   1,
			expResult:    1,
		},
	}
	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			// Given
			mockUniversalClient := new(mockredis.MockUniversalClient)

			// Mocks
			tc.givenMockRedisClientArg = tc.givenMockRedisClientArgFn()
			mockUniversalClient.ExpectedCalls = []*mock.Call{
				mockUniversalClient.On(
					"IncrByFloat",
					tc.givenMockRedisClientArg.givenContext,
					tc.givenMockRedisClientArg.givenKey,
					tc.givenMockRedisClientArg.givenValue,
				).Return(tc.givenMockRedisClientArg.expCmd),
			}

			// When
			instance := redisClient{
				rdb: mockUniversalClient,
			}
			result, err := instance.IncrementFloatBy(tc.givenContext, tc.givenKey, tc.givenValue)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, result, tc.expResult)
			}
		})
	}
}

func Test_redisClient_GetFloat(t *testing.T) {
	type mockCmdArg struct {
		givenContext context.Context
		givenKey     string
		expCmd       *redis.StringCmd
	}

	type arg[T Type] struct {
		givenMockCmdArgFn func() mockCmdArg
		givenMockCmdArg   mockCmdArg
		givenContext      context.Context
		givenKey          string
		expResult         T
		expErr            error
	}

	tcs := map[string]arg[float64]{
		"nil": {
			givenMockCmdArgFn: func() mockCmdArg {
				var cmd redis.StringCmd
				cmd.SetErr(redis.Nil)
				return mockCmdArg{
					givenContext: context.Background(),
					givenKey:     "key",
					expCmd:       &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
		},
		"error": {
			givenMockCmdArgFn: func() mockCmdArg {
				var cmd redis.StringCmd
				cmd.SetErr(errors.New("error"))
				return mockCmdArg{
					givenContext: context.Background(),
					givenKey:     "key",
					expCmd:       &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			expErr:       errors.New("error"),
		},
		"success": {
			givenMockCmdArgFn: func() mockCmdArg {
				var cmd redis.StringCmd
				cmd.SetVal("1")
				return mockCmdArg{
					givenContext: context.Background(),
					givenKey:     "key",
					expCmd:       &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			expResult:    1,
		},
	}

	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			// Given
			mockUniversalClient := new(mockredis.MockUniversalClient)

			// Mock
			tc.givenMockCmdArg = tc.givenMockCmdArgFn()
			mockUniversalClient.ExpectedCalls = []*mock.Call{
				mockUniversalClient.On(
					"Get",
					tc.givenMockCmdArg.givenContext,
					tc.givenMockCmdArg.givenKey,
				).Return(tc.givenMockCmdArg.expCmd),
			}

			// When
			instance := redisClient{
				rdb: mockUniversalClient,
			}
			result, err := instance.GetFloat(tc.givenContext, tc.givenKey)

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

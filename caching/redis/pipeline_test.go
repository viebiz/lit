package redis

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/viebiz/lit/mocks/mockredis"
)

func Test_commander_Discard(t *testing.T) {
	type mockPipelinerArg struct {
	}
	type arg struct {
		givenMockPipelinerArgFn func() mockPipelinerArg
		givenPipelinerArg       mockPipelinerArg
	}
	tcs := map[string]arg{
		"success": {
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				return mockPipelinerArg{}
			},
		},
	}
	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()
			// Mocks
			mockPipeliner := new(mockredis.MockPipeliner)

			// Given
			tc.givenPipelinerArg = tc.givenMockPipelinerArgFn()
			mockPipeliner.ExpectedCalls = []*mock.Call{
				mockPipeliner.On("Discard").Return(),
			}

			// When
			instance := commander{
				pipeliner: mockPipeliner,
			}
			instance.Discard()

			// Then

		})
	}
}

func Test_commander_Execute(t *testing.T) {
	type mockPipelinerArg struct {
		givenContext context.Context
		expCmd       []redis.Cmder
		expErr       error
	}
	type arg struct {
		givenMockPipelinerArgFn func() mockPipelinerArg
		givenPipelinerArg       mockPipelinerArg
		givenContext            context.Context
		expErr                  error
	}
	tcs := map[string]arg{
		"error": {
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd []redis.Cmder
				return mockPipelinerArg{
					givenContext: context.Background(),
					expCmd:       cmd,
					expErr:       errors.New("error"),
				}
			},
			givenContext: context.Background(),
			expErr:       errors.New("error"),
		},
		"success": {
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd []redis.Cmder
				return mockPipelinerArg{
					givenContext: context.Background(),
					expCmd:       cmd,
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
			mockPipeliner := new(mockredis.MockPipeliner)

			// Given
			tc.givenPipelinerArg = tc.givenMockPipelinerArgFn()
			mockPipeliner.ExpectedCalls = []*mock.Call{
				mockPipeliner.On("Exec", tc.givenPipelinerArg.givenContext).Return(tc.givenPipelinerArg.expCmd, tc.givenPipelinerArg.expErr),
			}

			// When
			instance := commander{
				pipeliner: mockPipeliner,
			}
			err := instance.Execute(tc.givenContext)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_commander_Delete(t *testing.T) {
	type mockPipelinerArg struct {
		givenContext context.Context
		givenKey     string
		expCmd       *redis.IntCmd
	}
	type arg struct {
		givenMockPipelinerArgFn func() mockPipelinerArg
		givenPipelinerArg       mockPipelinerArg
		givenContext            context.Context
		givenKey                string
		expResult               int64
		expErr                  error
	}
	tcs := map[string]arg{
		"error": {
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.IntCmd
				cmd.SetErr(errors.New("error"))
				return mockPipelinerArg{
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
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.IntCmd
				cmd.SetVal(1)
				return mockPipelinerArg{
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
			// Mocks
			mockPipeliner := new(mockredis.MockPipeliner)

			// Given
			tc.givenPipelinerArg = tc.givenMockPipelinerArgFn()
			mockPipeliner.ExpectedCalls = []*mock.Call{
				mockPipeliner.On("Del", tc.givenPipelinerArg.givenContext, tc.givenPipelinerArg.givenKey).Return(tc.givenPipelinerArg.expCmd),
			}

			// When
			instance := commander{
				pipeliner: mockPipeliner,
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

func Test_commander_Expire(t *testing.T) {
	type mockPipelinerArg struct {
		givenContext context.Context
		givenKey     string
		givenExpiry  time.Duration
		expCmd       *redis.BoolCmd
	}
	type arg struct {
		givenMockPipelinerArgFn func() mockPipelinerArg
		givenPipelinerArg       mockPipelinerArg
		givenContext            context.Context
		givenKey                string
		givenExpiry             time.Duration
		expResult               bool
		expErr                  error
	}
	tcs := map[string]arg{
		"error": {
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.BoolCmd
				cmd.SetErr(errors.New("error"))
				return mockPipelinerArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenExpiry:  time.Duration(0),
					expCmd:       &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			givenExpiry:  time.Duration(0),
			expErr:       errors.New("error"),
		},
		"success": {
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.BoolCmd
				cmd.SetVal(true)
				return mockPipelinerArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenExpiry:  time.Duration(1),
					expCmd:       &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			givenExpiry:  time.Duration(1),
			expResult:    true,
		},
	}
	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()
			// Mocks
			mockPipeliner := new(mockredis.MockPipeliner)

			// Given
			tc.givenPipelinerArg = tc.givenMockPipelinerArgFn()
			mockPipeliner.ExpectedCalls = []*mock.Call{
				mockPipeliner.On("Expire", tc.givenPipelinerArg.givenContext, tc.givenPipelinerArg.givenKey, tc.givenPipelinerArg.givenExpiry).Return(tc.givenPipelinerArg.expCmd),
			}

			// When
			instance := commander{
				pipeliner: mockPipeliner,
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

func Test_commander_SetString(t *testing.T) {
	type mockPipelinerArg struct {
		givenContext context.Context
		givenKey     string
		givenValue   interface{}
		givenArg     redis.SetArgs
		expCmd       *redis.StatusCmd
	}
	type arg struct {
		givenMockPipelinerArgFn func() mockPipelinerArg
		givenPipelinerArg       mockPipelinerArg
		givenContext            context.Context
		givenKey                string
		givenValue              string
		givenExpiry             time.Duration
		expErr                  error
	}
	tcs := map[string]arg{
		"error: no expiry": {
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.StatusCmd
				cmd.SetErr(errors.New("error: no expiry"))
				return mockPipelinerArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   "value",
					givenArg: redis.SetArgs{
						KeepTTL: true,
					},
					expCmd: &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			givenValue:   "value",
			givenExpiry:  time.Duration(0),
			expErr:       errors.New("error: no expiry"),
		},
		"error": {
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.StatusCmd
				cmd.SetErr(errors.New("error"))
				return mockPipelinerArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   "value",
					givenArg: redis.SetArgs{
						TTL: time.Duration(1),
					},
					expCmd: &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			givenValue:   "value",
			givenExpiry:  time.Duration(1),
			expErr:       errors.New("error"),
		},
		"success": {
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.StatusCmd
				cmd.SetVal("OK")
				return mockPipelinerArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   "value",
					givenArg: redis.SetArgs{
						TTL: time.Duration(1),
					},
					expCmd: &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			givenValue:   "value",
			givenExpiry:  time.Duration(1),
		},
	}
	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()
			// Mocks
			mockPipeliner := new(mockredis.MockPipeliner)

			// Given
			tc.givenPipelinerArg = tc.givenMockPipelinerArgFn()
			mockPipeliner.ExpectedCalls = []*mock.Call{
				mockPipeliner.On(
					"SetArgs",
					tc.givenPipelinerArg.givenContext,
					tc.givenPipelinerArg.givenKey,
					tc.givenPipelinerArg.givenValue,
					tc.givenPipelinerArg.givenArg,
				).Return(tc.givenPipelinerArg.expCmd),
			}

			// When
			instance := commander{
				pipeliner: mockPipeliner,
			}
			err := instance.SetString(tc.givenContext, tc.givenKey, tc.givenValue, tc.givenExpiry)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_commander_SetStringIfNotExist(t *testing.T) {
	type mockPipelinerArg struct {
		givenContext context.Context
		givenKey     string
		givenValue   interface{}
		givenArg     redis.SetArgs
		expCmd       *redis.StatusCmd
	}
	type arg struct {
		givenMockPipelinerArgFn func() mockPipelinerArg
		givenPipelinerArg       mockPipelinerArg
		givenContext            context.Context
		givenKey                string
		givenValue              string
		givenExpiry             time.Duration
		expErr                  error
	}
	tcs := map[string]arg{
		"error: no expiry": {
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.StatusCmd
				cmd.SetErr(errors.New("error: no expiry"))
				return mockPipelinerArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   "value",
					givenArg: redis.SetArgs{
						KeepTTL: true,
						Mode:    setModeNX.String(),
					},
					expCmd: &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			givenValue:   "value",
			givenExpiry:  time.Duration(0),
			expErr:       errors.New("error: no expiry"),
		},
		"error": {
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.StatusCmd
				cmd.SetErr(errors.New("error"))
				return mockPipelinerArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   "value",
					givenArg: redis.SetArgs{
						TTL:  time.Duration(1),
						Mode: setModeNX.String(),
					},
					expCmd: &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			givenValue:   "value",
			givenExpiry:  time.Duration(1),
			expErr:       errors.New("error"),
		},
		"error: ErrFailToSetValue": {
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.StatusCmd
				cmd.SetVal("ok")
				return mockPipelinerArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   "value",
					givenArg: redis.SetArgs{
						TTL:  time.Duration(1),
						Mode: setModeNX.String(),
					},
					expCmd: &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			givenValue:   "value",
			givenExpiry:  time.Duration(1),
			expErr:       ErrFailToSetValue,
		},
		"success": {
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.StatusCmd
				cmd.SetVal("OK")
				return mockPipelinerArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   "value",
					givenArg: redis.SetArgs{
						TTL:  time.Duration(1),
						Mode: setModeNX.String(),
					},
					expCmd: &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			givenValue:   "value",
			givenExpiry:  time.Duration(1),
		},
	}
	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()
			// Mocks
			mockPipeliner := new(mockredis.MockPipeliner)

			// Given
			tc.givenPipelinerArg = tc.givenMockPipelinerArgFn()
			mockPipeliner.ExpectedCalls = []*mock.Call{
				mockPipeliner.On(
					"SetArgs",
					tc.givenPipelinerArg.givenContext,
					tc.givenPipelinerArg.givenKey,
					tc.givenPipelinerArg.givenValue,
					tc.givenPipelinerArg.givenArg,
				).Return(tc.givenPipelinerArg.expCmd),
			}

			// When
			instance := commander{
				pipeliner: mockPipeliner,
			}
			err := instance.SetStringIfNotExist(tc.givenContext, tc.givenKey, tc.givenValue, tc.givenExpiry)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_commander_SetStringIfExist(t *testing.T) {
	type mockPipelinerArg struct {
		givenContext context.Context
		givenKey     string
		givenValue   interface{}
		givenArg     redis.SetArgs
		expCmd       *redis.StatusCmd
	}
	type arg struct {
		givenMockPipelinerArgFn func() mockPipelinerArg
		givenPipelinerArg       mockPipelinerArg
		givenContext            context.Context
		givenKey                string
		givenValue              string
		givenExpiry             time.Duration
		expErr                  error
	}
	tcs := map[string]arg{
		"error: no expiry": {
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.StatusCmd
				cmd.SetErr(errors.New("error: no expiry"))
				return mockPipelinerArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   "value",
					givenArg: redis.SetArgs{
						KeepTTL: true,
						Mode:    setModeXX.String(),
					},
					expCmd: &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			givenValue:   "value",
			givenExpiry:  time.Duration(0),
			expErr:       errors.New("error: no expiry"),
		},
		"error": {
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.StatusCmd
				cmd.SetErr(errors.New("error"))
				return mockPipelinerArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   "value",
					givenArg: redis.SetArgs{
						TTL:  time.Duration(1),
						Mode: setModeXX.String(),
					},
					expCmd: &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			givenValue:   "value",
			givenExpiry:  time.Duration(1),
			expErr:       errors.New("error"),
		},
		"error: ErrFailToSetValue": {
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.StatusCmd
				cmd.SetVal("ok")
				return mockPipelinerArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   "value",
					givenArg: redis.SetArgs{
						TTL:  time.Duration(1),
						Mode: setModeXX.String(),
					},
					expCmd: &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			givenValue:   "value",
			givenExpiry:  time.Duration(1),
			expErr:       ErrFailToSetValue,
		},
		"success": {
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.StatusCmd
				cmd.SetVal("OK")
				return mockPipelinerArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   "value",
					givenArg: redis.SetArgs{
						TTL:  time.Duration(1),
						Mode: setModeXX.String(),
					},
					expCmd: &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			givenValue:   "value",
			givenExpiry:  time.Duration(1),
		},
	}
	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()
			// Mocks
			mockPipeliner := new(mockredis.MockPipeliner)

			// Given
			tc.givenPipelinerArg = tc.givenMockPipelinerArgFn()
			mockPipeliner.ExpectedCalls = []*mock.Call{
				mockPipeliner.On(
					"SetArgs",
					tc.givenPipelinerArg.givenContext,
					tc.givenPipelinerArg.givenKey,
					tc.givenPipelinerArg.givenValue,
					tc.givenPipelinerArg.givenArg,
				).Return(tc.givenPipelinerArg.expCmd),
			}

			// When
			instance := commander{
				pipeliner: mockPipeliner,
			}
			err := instance.SetStringIfExist(tc.givenContext, tc.givenKey, tc.givenValue, tc.givenExpiry)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_commander_GetString(t *testing.T) {
	type mockPipelinerArg struct {
		givenContext context.Context
		givenKey     string
		expCmd       *redis.StringCmd
	}
	type arg struct {
		givenMockPipelinerArgFn func() mockPipelinerArg
		givenPipelinerArg       mockPipelinerArg
		givenContext            context.Context
		givenKey                string
		expResult               string
		expErr                  error
	}
	tcs := map[string]arg{
		"nil": {
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.StringCmd
				cmd.SetErr(redis.Nil)
				return mockPipelinerArg{
					givenContext: context.Background(),
					givenKey:     "key",
					expCmd:       &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			expResult:    "",
		},
		"error": {
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.StringCmd
				cmd.SetErr(errors.New("error"))
				return mockPipelinerArg{
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
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.StringCmd
				cmd.SetVal("value")
				return mockPipelinerArg{
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
			// Mocks
			mockPipeliner := new(mockredis.MockPipeliner)

			// Given
			tc.givenPipelinerArg = tc.givenMockPipelinerArgFn()
			mockPipeliner.ExpectedCalls = []*mock.Call{
				mockPipeliner.On(
					"Get",
					tc.givenPipelinerArg.givenContext,
					tc.givenPipelinerArg.givenKey,
				).Return(tc.givenPipelinerArg.expCmd),
			}

			// When
			instance := commander{
				pipeliner: mockPipeliner,
			}
			result, err := instance.GetString(tc.givenContext, tc.givenKey)

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

func Test_commander_SetInt(t *testing.T) {
	type mockPipelinerArg struct {
		givenContext context.Context
		givenKey     string
		givenValue   interface{}
		givenArg     redis.SetArgs
		expCmd       *redis.StatusCmd
	}
	type arg struct {
		givenMockPipelinerArgFn func() mockPipelinerArg
		givenPipelinerArg       mockPipelinerArg
		givenContext            context.Context
		givenKey                string
		givenValue              int64
		givenExpiry             time.Duration
		expErr                  error
	}
	tcs := map[string]arg{
		"error: no expiry": {
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.StatusCmd
				cmd.SetErr(errors.New("error: no expiry"))
				return mockPipelinerArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   int64(1),
					givenArg: redis.SetArgs{
						KeepTTL: true,
					},
					expCmd: &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			givenValue:   1,
			givenExpiry:  time.Duration(0),
			expErr:       errors.New("error: no expiry"),
		},
		"error": {
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.StatusCmd
				cmd.SetErr(errors.New("error"))
				return mockPipelinerArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   int64(1),
					givenArg: redis.SetArgs{
						TTL: time.Duration(1),
					},
					expCmd: &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			givenValue:   1,
			givenExpiry:  time.Duration(1),
			expErr:       errors.New("error"),
		},
		"success": {
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.StatusCmd
				cmd.SetVal("OK")
				return mockPipelinerArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   int64(1),
					givenArg: redis.SetArgs{
						TTL: time.Duration(1),
					},
					expCmd: &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			givenValue:   1,
			givenExpiry:  time.Duration(1),
		},
	}
	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()
			// Mocks
			mockPipeliner := new(mockredis.MockPipeliner)

			// Given
			tc.givenPipelinerArg = tc.givenMockPipelinerArgFn()
			mockPipeliner.ExpectedCalls = []*mock.Call{
				mockPipeliner.On(
					"SetArgs",
					tc.givenPipelinerArg.givenContext,
					tc.givenPipelinerArg.givenKey,
					tc.givenPipelinerArg.givenValue,
					tc.givenPipelinerArg.givenArg,
				).Return(tc.givenPipelinerArg.expCmd),
			}

			// When
			instance := commander{
				pipeliner: mockPipeliner,
			}
			err := instance.SetInt(tc.givenContext, tc.givenKey, tc.givenValue, tc.givenExpiry)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_commander_SetIntIfNotExist(t *testing.T) {
	type mockPipelinerArg struct {
		givenContext context.Context
		givenKey     string
		givenValue   interface{}
		givenArg     redis.SetArgs
		expCmd       *redis.StatusCmd
	}
	type arg struct {
		givenMockPipelinerArgFn func() mockPipelinerArg
		givenPipelinerArg       mockPipelinerArg
		givenContext            context.Context
		givenKey                string
		givenValue              int64
		givenExpiry             time.Duration
		expErr                  error
	}
	tcs := map[string]arg{
		"error: no expiry": {
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.StatusCmd
				cmd.SetErr(errors.New("error: no expiry"))
				return mockPipelinerArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   int64(1),
					givenArg: redis.SetArgs{
						KeepTTL: true,
						Mode:    setModeNX.String(),
					},
					expCmd: &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			givenValue:   1,
			givenExpiry:  time.Duration(0),
			expErr:       errors.New("error: no expiry"),
		},
		"error": {
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.StatusCmd
				cmd.SetErr(errors.New("error"))
				return mockPipelinerArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   int64(1),
					givenArg: redis.SetArgs{
						TTL:  time.Duration(1),
						Mode: setModeNX.String(),
					},
					expCmd: &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			givenValue:   1,
			givenExpiry:  time.Duration(1),
			expErr:       errors.New("error"),
		},
		"error: ErrFailToSetValue": {
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.StatusCmd
				cmd.SetVal("ok")
				return mockPipelinerArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   int64(1),
					givenArg: redis.SetArgs{
						TTL:  time.Duration(1),
						Mode: setModeNX.String(),
					},
					expCmd: &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			givenValue:   1,
			givenExpiry:  time.Duration(1),
			expErr:       ErrFailToSetValue,
		},
		"success": {
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.StatusCmd
				cmd.SetVal("OK")
				return mockPipelinerArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   int64(1),
					givenArg: redis.SetArgs{
						TTL:  time.Duration(1),
						Mode: setModeNX.String(),
					},
					expCmd: &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			givenValue:   1,
			givenExpiry:  time.Duration(1),
		},
	}
	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()
			// Mocks
			mockPipeliner := new(mockredis.MockPipeliner)

			// Given
			tc.givenPipelinerArg = tc.givenMockPipelinerArgFn()
			mockPipeliner.ExpectedCalls = []*mock.Call{
				mockPipeliner.On(
					"SetArgs",
					tc.givenPipelinerArg.givenContext,
					tc.givenPipelinerArg.givenKey,
					tc.givenPipelinerArg.givenValue,
					tc.givenPipelinerArg.givenArg,
				).Return(tc.givenPipelinerArg.expCmd),
			}

			// When
			instance := commander{
				pipeliner: mockPipeliner,
			}
			err := instance.SetIntIfNotExist(tc.givenContext, tc.givenKey, tc.givenValue, tc.givenExpiry)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_commander_SetIntIfExist(t *testing.T) {
	type mockPipelinerArg struct {
		givenContext context.Context
		givenKey     string
		givenValue   interface{}
		givenArg     redis.SetArgs
		expCmd       *redis.StatusCmd
	}
	type arg struct {
		givenMockPipelinerArgFn func() mockPipelinerArg
		givenPipelinerArg       mockPipelinerArg
		givenContext            context.Context
		givenKey                string
		givenValue              int64
		givenExpiry             time.Duration
		expErr                  error
	}
	tcs := map[string]arg{
		"error: no expiry": {
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.StatusCmd
				cmd.SetErr(errors.New("error: no expiry"))
				return mockPipelinerArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   int64(1),
					givenArg: redis.SetArgs{
						KeepTTL: true,
						Mode:    setModeXX.String(),
					},
					expCmd: &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			givenValue:   1,
			givenExpiry:  time.Duration(0),
			expErr:       errors.New("error: no expiry"),
		},
		"error": {
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.StatusCmd
				cmd.SetErr(errors.New("error"))
				return mockPipelinerArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   int64(1),
					givenArg: redis.SetArgs{
						TTL:  time.Duration(1),
						Mode: setModeXX.String(),
					},
					expCmd: &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			givenValue:   1,
			givenExpiry:  time.Duration(1),
			expErr:       errors.New("error"),
		},
		"error: ErrFailToSetValue": {
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.StatusCmd
				cmd.SetVal("ok")
				return mockPipelinerArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   int64(1),
					givenArg: redis.SetArgs{
						TTL:  time.Duration(1),
						Mode: setModeXX.String(),
					},
					expCmd: &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			givenValue:   1,
			givenExpiry:  time.Duration(1),
			expErr:       ErrFailToSetValue,
		},
		"success": {
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.StatusCmd
				cmd.SetVal("OK")
				return mockPipelinerArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   int64(1),
					givenArg: redis.SetArgs{
						TTL:  time.Duration(1),
						Mode: setModeXX.String(),
					},
					expCmd: &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			givenValue:   1,
			givenExpiry:  time.Duration(1),
		},
	}
	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()
			// Mocks
			mockPipeliner := new(mockredis.MockPipeliner)

			// Given
			tc.givenPipelinerArg = tc.givenMockPipelinerArgFn()
			mockPipeliner.ExpectedCalls = []*mock.Call{
				mockPipeliner.On(
					"SetArgs",
					tc.givenPipelinerArg.givenContext,
					tc.givenPipelinerArg.givenKey,
					tc.givenPipelinerArg.givenValue,
					tc.givenPipelinerArg.givenArg,
				).Return(tc.givenPipelinerArg.expCmd),
			}

			// When
			instance := commander{
				pipeliner: mockPipeliner,
			}
			err := instance.SetIntIfExist(tc.givenContext, tc.givenKey, tc.givenValue, tc.givenExpiry)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_commander_GetInt(t *testing.T) {
	type mockPipelinerArg struct {
		givenContext context.Context
		givenKey     string
		expCmd       *redis.StringCmd
	}
	type arg struct {
		givenMockPipelinerArgFn func() mockPipelinerArg
		givenPipelinerArg       mockPipelinerArg
		givenContext            context.Context
		givenKey                string
		expResult               int64
		expErr                  error
	}
	tcs := map[string]arg{
		"nil": {
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.StringCmd
				cmd.SetErr(redis.Nil)
				return mockPipelinerArg{
					givenContext: context.Background(),
					givenKey:     "key",
					expCmd:       &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			expResult:    0,
		},
		"error": {
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.StringCmd
				cmd.SetErr(errors.New("error"))
				return mockPipelinerArg{
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
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.StringCmd
				cmd.SetVal("1")
				return mockPipelinerArg{
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
			// Mocks
			mockPipeliner := new(mockredis.MockPipeliner)

			// Given
			tc.givenPipelinerArg = tc.givenMockPipelinerArgFn()
			mockPipeliner.ExpectedCalls = []*mock.Call{
				mockPipeliner.On(
					"Get",
					tc.givenPipelinerArg.givenContext,
					tc.givenPipelinerArg.givenKey,
				).Return(tc.givenPipelinerArg.expCmd),
			}

			// When
			instance := commander{
				pipeliner: mockPipeliner,
			}
			result, err := instance.GetInt(tc.givenContext, tc.givenKey)

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

func Test_commander_SetFloat(t *testing.T) {
	type mockPipelinerArg struct {
		givenContext context.Context
		givenKey     string
		givenValue   interface{}
		givenArg     redis.SetArgs
		expCmd       *redis.StatusCmd
	}
	type arg struct {
		givenMockPipelinerArgFn func() mockPipelinerArg
		givenPipelinerArg       mockPipelinerArg
		givenContext            context.Context
		givenKey                string
		givenValue              float64
		givenExpiry             time.Duration
		expErr                  error
	}
	tcs := map[string]arg{
		"error: no expiry": {
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.StatusCmd
				cmd.SetErr(errors.New("error: no expiry"))
				return mockPipelinerArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   float64(1),
					givenArg: redis.SetArgs{
						KeepTTL: true,
					},
					expCmd: &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			givenValue:   1,
			givenExpiry:  time.Duration(0),
			expErr:       errors.New("error: no expiry"),
		},
		"error": {
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.StatusCmd
				cmd.SetErr(errors.New("error"))
				return mockPipelinerArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   float64(1),
					givenArg: redis.SetArgs{
						TTL: time.Duration(1),
					},
					expCmd: &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			givenValue:   1,
			givenExpiry:  time.Duration(1),
			expErr:       errors.New("error"),
		},
		"success": {
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.StatusCmd
				cmd.SetVal("OK")
				return mockPipelinerArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   float64(1),
					givenArg: redis.SetArgs{
						TTL: time.Duration(1),
					},
					expCmd: &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			givenValue:   1,
			givenExpiry:  time.Duration(1),
		},
	}
	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()
			// Mocks
			mockPipeliner := new(mockredis.MockPipeliner)

			// Given
			tc.givenPipelinerArg = tc.givenMockPipelinerArgFn()
			mockPipeliner.ExpectedCalls = []*mock.Call{
				mockPipeliner.On(
					"SetArgs",
					tc.givenPipelinerArg.givenContext,
					tc.givenPipelinerArg.givenKey,
					tc.givenPipelinerArg.givenValue,
					tc.givenPipelinerArg.givenArg,
				).Return(tc.givenPipelinerArg.expCmd),
			}

			// When
			instance := commander{
				pipeliner: mockPipeliner,
			}
			err := instance.SetFloat(tc.givenContext, tc.givenKey, tc.givenValue, tc.givenExpiry)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_commander_SetFloatIfNotExist(t *testing.T) {
	type mockPipelinerArg struct {
		givenContext context.Context
		givenKey     string
		givenValue   interface{}
		givenArg     redis.SetArgs
		expCmd       *redis.StatusCmd
	}
	type arg struct {
		givenMockPipelinerArgFn func() mockPipelinerArg
		givenPipelinerArg       mockPipelinerArg
		givenContext            context.Context
		givenKey                string
		givenValue              float64
		givenExpiry             time.Duration
		expErr                  error
	}
	tcs := map[string]arg{
		"error: no expiry": {
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.StatusCmd
				cmd.SetErr(errors.New("error: no expiry"))
				return mockPipelinerArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   float64(1),
					givenArg: redis.SetArgs{
						KeepTTL: true,
						Mode:    setModeNX.String(),
					},
					expCmd: &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			givenValue:   1,
			givenExpiry:  time.Duration(0),
			expErr:       errors.New("error: no expiry"),
		},
		"error": {
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.StatusCmd
				cmd.SetErr(errors.New("error"))
				return mockPipelinerArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   float64(1),
					givenArg: redis.SetArgs{
						TTL:  time.Duration(1),
						Mode: setModeNX.String(),
					},
					expCmd: &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			givenValue:   1,
			givenExpiry:  time.Duration(1),
			expErr:       errors.New("error"),
		},
		"error: ErrFailToSetValue": {
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.StatusCmd
				cmd.SetVal("ok")
				return mockPipelinerArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   float64(1),
					givenArg: redis.SetArgs{
						TTL:  time.Duration(1),
						Mode: setModeNX.String(),
					},
					expCmd: &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			givenValue:   1,
			givenExpiry:  time.Duration(1),
			expErr:       ErrFailToSetValue,
		},
		"success": {
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.StatusCmd
				cmd.SetVal("OK")
				return mockPipelinerArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   float64(1),
					givenArg: redis.SetArgs{
						TTL:  time.Duration(1),
						Mode: setModeNX.String(),
					},
					expCmd: &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			givenValue:   1,
			givenExpiry:  time.Duration(1),
		},
	}
	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()
			// Mocks
			mockPipeliner := new(mockredis.MockPipeliner)

			// Given
			tc.givenPipelinerArg = tc.givenMockPipelinerArgFn()
			mockPipeliner.ExpectedCalls = []*mock.Call{
				mockPipeliner.On(
					"SetArgs",
					tc.givenPipelinerArg.givenContext,
					tc.givenPipelinerArg.givenKey,
					tc.givenPipelinerArg.givenValue,
					tc.givenPipelinerArg.givenArg,
				).Return(tc.givenPipelinerArg.expCmd),
			}

			// When
			instance := commander{
				pipeliner: mockPipeliner,
			}
			err := instance.SetFloatIfNotExist(tc.givenContext, tc.givenKey, tc.givenValue, tc.givenExpiry)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_commander_SetFloatIfExist(t *testing.T) {
	type mockPipelinerArg struct {
		givenContext context.Context
		givenKey     string
		givenValue   interface{}
		givenArg     redis.SetArgs
		expCmd       *redis.StatusCmd
	}
	type arg struct {
		givenMockPipelinerArgFn func() mockPipelinerArg
		givenPipelinerArg       mockPipelinerArg
		givenContext            context.Context
		givenKey                string
		givenValue              float64
		givenExpiry             time.Duration
		expErr                  error
	}
	tcs := map[string]arg{
		"error: no expiry": {
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.StatusCmd
				cmd.SetErr(errors.New("error: no expiry"))
				return mockPipelinerArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   float64(1),
					givenArg: redis.SetArgs{
						KeepTTL: true,
						Mode:    setModeXX.String(),
					},
					expCmd: &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			givenValue:   1,
			givenExpiry:  time.Duration(0),
			expErr:       errors.New("error: no expiry"),
		},
		"error": {
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.StatusCmd
				cmd.SetErr(errors.New("error"))
				return mockPipelinerArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   float64(1),
					givenArg: redis.SetArgs{
						TTL:  time.Duration(1),
						Mode: setModeXX.String(),
					},
					expCmd: &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			givenValue:   1,
			givenExpiry:  time.Duration(1),
			expErr:       errors.New("error"),
		},
		"error: ErrFailToSetValue": {
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.StatusCmd
				cmd.SetVal("ok")
				return mockPipelinerArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   float64(1),
					givenArg: redis.SetArgs{
						TTL:  time.Duration(1),
						Mode: setModeXX.String(),
					},
					expCmd: &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			givenValue:   1,
			givenExpiry:  time.Duration(1),
			expErr:       ErrFailToSetValue,
		},
		"success": {
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.StatusCmd
				cmd.SetVal("OK")
				return mockPipelinerArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   float64(1),
					givenArg: redis.SetArgs{
						TTL:  time.Duration(1),
						Mode: setModeXX.String(),
					},
					expCmd: &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			givenValue:   1,
			givenExpiry:  time.Duration(1),
		},
	}
	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()
			// Mocks
			mockPipeliner := new(mockredis.MockPipeliner)

			// Given
			tc.givenPipelinerArg = tc.givenMockPipelinerArgFn()
			mockPipeliner.ExpectedCalls = []*mock.Call{
				mockPipeliner.On(
					"SetArgs",
					tc.givenPipelinerArg.givenContext,
					tc.givenPipelinerArg.givenKey,
					tc.givenPipelinerArg.givenValue,
					tc.givenPipelinerArg.givenArg,
				).Return(tc.givenPipelinerArg.expCmd),
			}

			// When
			instance := commander{
				pipeliner: mockPipeliner,
			}
			err := instance.SetFloatIfExist(tc.givenContext, tc.givenKey, tc.givenValue, tc.givenExpiry)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_commander_GetFloat(t *testing.T) {
	type mockPipelinerArg struct {
		givenContext context.Context
		givenKey     string
		expCmd       *redis.StringCmd
	}
	type arg struct {
		givenMockPipelinerArgFn func() mockPipelinerArg
		givenPipelinerArg       mockPipelinerArg
		givenContext            context.Context
		givenKey                string
		expResult               float64
		expErr                  error
	}
	tcs := map[string]arg{
		"nil": {
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.StringCmd
				cmd.SetErr(redis.Nil)
				return mockPipelinerArg{
					givenContext: context.Background(),
					givenKey:     "key",
					expCmd:       &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			expResult:    0,
		},
		"error": {
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.StringCmd
				cmd.SetErr(errors.New("error"))
				return mockPipelinerArg{
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
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.StringCmd
				cmd.SetVal("1")
				return mockPipelinerArg{
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
			// Mocks
			mockPipeliner := new(mockredis.MockPipeliner)

			// Given
			tc.givenPipelinerArg = tc.givenMockPipelinerArgFn()
			mockPipeliner.ExpectedCalls = []*mock.Call{
				mockPipeliner.On(
					"Get",
					tc.givenPipelinerArg.givenContext,
					tc.givenPipelinerArg.givenKey,
				).Return(tc.givenPipelinerArg.expCmd),
			}

			// When
			instance := commander{
				pipeliner: mockPipeliner,
			}
			result, err := instance.GetFloat(tc.givenContext, tc.givenKey)

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

func Test_commander_IncrementBy(t *testing.T) {
	type mockPipelinerArg struct {
		givenContext context.Context
		givenKey     string
		givenValue   int64
		expCmd       *redis.IntCmd
	}
	type arg struct {
		givenMockPipelinerArgFn func() mockPipelinerArg
		givenPipelinerArg       mockPipelinerArg
		givenContext            context.Context
		givenKey                string
		givenValue              int64
		expResult               int64
		expErr                  error
	}
	tcs := map[string]arg{
		"error": {
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.IntCmd
				cmd.SetErr(errors.New("error"))
				return mockPipelinerArg{
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
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.IntCmd
				cmd.SetVal(1)
				return mockPipelinerArg{
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
			// Mocks
			mockPipeliner := new(mockredis.MockPipeliner)

			// Given
			tc.givenPipelinerArg = tc.givenMockPipelinerArgFn()
			mockPipeliner.ExpectedCalls = []*mock.Call{
				mockPipeliner.On(
					"IncrBy",
					tc.givenPipelinerArg.givenContext,
					tc.givenPipelinerArg.givenKey,
					tc.givenPipelinerArg.givenValue,
				).Return(tc.givenPipelinerArg.expCmd),
			}

			// When
			instance := commander{
				pipeliner: mockPipeliner,
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

func Test_commander_DecrementBy(t *testing.T) {
	type mockPipelinerArg struct {
		givenContext context.Context
		givenKey     string
		givenValue   int64
		expCmd       *redis.IntCmd
	}
	type arg struct {
		givenMockPipelinerArgFn func() mockPipelinerArg
		givenPipelinerArg       mockPipelinerArg
		givenContext            context.Context
		givenKey                string
		givenValue              int64
		expResult               int64
		expErr                  error
	}
	tcs := map[string]arg{
		"error": {
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.IntCmd
				cmd.SetErr(errors.New("error"))
				return mockPipelinerArg{
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
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.IntCmd
				cmd.SetVal(1)
				return mockPipelinerArg{
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
			// Mocks
			mockPipeliner := new(mockredis.MockPipeliner)

			// Given
			tc.givenPipelinerArg = tc.givenMockPipelinerArgFn()
			mockPipeliner.ExpectedCalls = []*mock.Call{
				mockPipeliner.On(
					"DecrBy",
					tc.givenPipelinerArg.givenContext,
					tc.givenPipelinerArg.givenKey,
					tc.givenPipelinerArg.givenValue,
				).Return(tc.givenPipelinerArg.expCmd),
			}

			// When
			instance := commander{
				pipeliner: mockPipeliner,
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

func Test_commander_IncrementFloatBy(t *testing.T) {
	type mockPipelinerArg struct {
		givenContext context.Context
		givenKey     string
		givenValue   float64
		expCmd       *redis.FloatCmd
	}
	type arg struct {
		givenMockPipelinerArgFn func() mockPipelinerArg
		givenPipelinerArg       mockPipelinerArg
		givenContext            context.Context
		givenKey                string
		givenValue              float64
		expResult               float64
		expErr                  error
	}
	tcs := map[string]arg{
		"error": {
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.FloatCmd
				cmd.SetErr(errors.New("error"))
				return mockPipelinerArg{
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
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.FloatCmd
				cmd.SetVal(1)
				return mockPipelinerArg{
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
			// Mocks
			mockPipeliner := new(mockredis.MockPipeliner)

			// Given
			tc.givenPipelinerArg = tc.givenMockPipelinerArgFn()
			mockPipeliner.ExpectedCalls = []*mock.Call{
				mockPipeliner.On(
					"IncrByFloat",
					tc.givenPipelinerArg.givenContext,
					tc.givenPipelinerArg.givenKey,
					tc.givenPipelinerArg.givenValue,
				).Return(tc.givenPipelinerArg.expCmd),
			}

			// When
			instance := commander{
				pipeliner: mockPipeliner,
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

func Test_commander_HashSet(t *testing.T) {
	type mockPipelinerArg struct {
		givenContext context.Context
		givenKey     string
		givenValue   interface{}
		expCmd       *redis.IntCmd
	}
	type arg struct {
		givenMockPipelinerArgFn func() mockPipelinerArg
		givenPipelinerArg       mockPipelinerArg
		givenContext            context.Context
		givenKey                string
		givenValue              interface{}
		expErr                  error
	}
	tcs := map[string]arg{
		"error: ErrUnsupportedInputType": {
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.IntCmd
				cmd.SetErr(errors.New("error"))
				return mockPipelinerArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   1,
					expCmd:       &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			givenValue:   1,
			expErr:       ErrUnsupportedInputType,
		},
		"error": {
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.IntCmd
				cmd.SetErr(errors.New("error"))
				return mockPipelinerArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenValue:   map[string]string{"key1": "value1"},
					expCmd:       &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			givenValue:   map[string]string{"key1": "value1"},
			expErr:       errors.New("error"),
		},
		"error: ErrFailToSetValue": {
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.IntCmd
				cmd.SetVal(0)
				return mockPipelinerArg{
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
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.IntCmd
				cmd.SetVal(1)
				return mockPipelinerArg{
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
			mockPipeliner := new(mockredis.MockPipeliner)

			// Given
			tc.givenPipelinerArg = tc.givenMockPipelinerArgFn()
			mockPipeliner.ExpectedCalls = []*mock.Call{
				mockPipeliner.On("HSet",
					tc.givenPipelinerArg.givenContext,
					tc.givenPipelinerArg.givenKey,
					tc.givenPipelinerArg.givenValue,
				).Return(tc.givenPipelinerArg.expCmd),
			}

			// When
			instance := commander{
				pipeliner: mockPipeliner,
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

func Test_commander_HashGetAll(t *testing.T) {
	type mockPipelinerArg struct {
		givenContext context.Context
		givenKey     string
		expCmd       *redis.MapStringStringCmd
	}
	type arg struct {
		givenMockPipelinerArgFn func() mockPipelinerArg
		givenPipelinerArg       mockPipelinerArg
		givenContext            context.Context
		givenKey                string
		givenValue              interface{}
		expErr                  error
	}
	tcs := map[string]arg{
		"error: ErrUnsupportedInputType": {
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.MapStringStringCmd
				cmd.SetErr(errors.New("error"))
				return mockPipelinerArg{
					givenContext: context.Background(),
					givenKey:     "key",
					expCmd:       &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			givenValue:   1,
			expErr:       ErrUnsupportedInputType,
		},
		"error": {
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.MapStringStringCmd
				cmd.SetErr(errors.New("error"))
				return mockPipelinerArg{
					givenContext: context.Background(),
					givenKey:     "key",
					expCmd:       &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			givenValue:   &map[string]string{},
			expErr:       errors.New("error"),
		},
		"success": {
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.MapStringStringCmd
				cmd.SetVal(map[string]string{"key1": "value1"})
				return mockPipelinerArg{
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
			mockPipeliner := new(mockredis.MockPipeliner)

			// Given
			tc.givenPipelinerArg = tc.givenMockPipelinerArgFn()
			mockPipeliner.ExpectedCalls = []*mock.Call{
				mockPipeliner.On("HGetAll",
					tc.givenPipelinerArg.givenContext,
					tc.givenPipelinerArg.givenKey,
				).Return(tc.givenPipelinerArg.expCmd),
			}

			// When
			instance := commander{
				pipeliner: mockPipeliner,
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

func Test_commander_HashGetField(t *testing.T) {
	type mockPipelinerArg struct {
		givenContext context.Context
		givenKey     string
		givenField   string
		expCmd       *redis.StringCmd
	}
	type arg struct {
		givenMockPipelinerArgFn func() mockPipelinerArg
		givenPipelinerArg       mockPipelinerArg
		givenContext            context.Context
		givenKey                string
		givenField              string
		givenValue              interface{}
		expErr                  error
	}
	tcs := map[string]arg{
		"error: ErrUnsupportedInputType": {
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.StringCmd
				cmd.SetErr(errors.New("error"))
				return mockPipelinerArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenField:   "givenKey",
					expCmd:       &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			givenField:   "givenKey",
			givenValue:   1,
			expErr:       ErrUnsupportedInputType,
		},
		"error": {
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.StringCmd
				cmd.SetErr(errors.New("error"))
				return mockPipelinerArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenField:   "givenKey",
					expCmd:       &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			givenField:   "givenKey",
			givenValue:   &map[string]string{},
			expErr:       errors.New("error"),
		},
		"success": {
			givenMockPipelinerArgFn: func() mockPipelinerArg {
				var cmd redis.StringCmd
				cmd.SetVal("ok")
				return mockPipelinerArg{
					givenContext: context.Background(),
					givenKey:     "key",
					givenField:   "givenKey",
					expCmd:       &cmd,
				}
			},
			givenContext: context.Background(),
			givenKey:     "key",
			givenField:   "givenKey",
			givenValue:   stringPtr(""),
		},
	}
	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()
			// Mocks
			mockPipeliner := new(mockredis.MockPipeliner)

			// Given
			tc.givenPipelinerArg = tc.givenMockPipelinerArgFn()
			mockPipeliner.ExpectedCalls = []*mock.Call{
				mockPipeliner.On("HGet",
					tc.givenPipelinerArg.givenContext,
					tc.givenPipelinerArg.givenKey,
					tc.givenPipelinerArg.givenField,
				).Return(tc.givenPipelinerArg.expCmd),
			}

			// When
			instance := commander{
				pipeliner: mockPipeliner,
			}
			err := instance.HashGetField(tc.givenContext, tc.givenKey, tc.givenField, tc.givenValue)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

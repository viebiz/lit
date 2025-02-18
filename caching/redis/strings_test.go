package redis

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/viebiz/lit/mocks/mockredis"
)

func Test_setSingleValue(t *testing.T) {
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
		givenMode         setMode
		expErr            error
	}

	dataTypes := map[string]interface{}{
		"string": map[string]arg[string]{
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
				givenMode:       setModeNone,
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
				givenMode:       setModeNone,
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
				givenMode:       setModeNX,
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
				givenMode:       setModeNX,
			},
		},
		"int64": map[string]arg[int64]{
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
				givenMode:       setModeNone,
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
				givenMode:       setModeNone,
				expErr:          errors.New("error"),
			},
			"error: ErrFailToSetValue": {
				givenMockCmdArgFn: func() mockCmdArg {
					var cmd redis.StatusCmd
					cmd.SetVal("ok")
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
				givenMode:       setModeNX,
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
				givenMode:       setModeNX,
			},
		},
		"uint64": map[string]arg[uint64]{
			"error: no expiry": {
				givenMockCmdArgFn: func() mockCmdArg {
					var cmd redis.StatusCmd
					cmd.SetErr(errors.New("error: no expiry"))
					return mockCmdArg{
						givenContext: context.Background(),
						givenKey:     "key",
						givenValue:   uint64(0),
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
				givenMode:       setModeNone,
				expErr:          errors.New("error: no expiry"),
			},
			"error": {
				givenMockCmdArgFn: func() mockCmdArg {
					var cmd redis.StatusCmd
					cmd.SetErr(errors.New("error"))
					return mockCmdArg{
						givenContext: context.Background(),
						givenKey:     "key",
						givenValue:   uint64(0),
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
				givenMode:       setModeNone,
				expErr:          errors.New("error"),
			},
			"error: ErrFailToSetValue": {
				givenMockCmdArgFn: func() mockCmdArg {
					var cmd redis.StatusCmd
					cmd.SetVal("ok")
					return mockCmdArg{
						givenContext: context.Background(),
						givenKey:     "key",
						givenValue:   uint64(1),
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
				givenMode:       setModeNX,
				expErr:          ErrFailToSetValue,
			},
			"success": {
				givenMockCmdArgFn: func() mockCmdArg {
					var cmd redis.StatusCmd
					cmd.SetVal(statusOK)
					return mockCmdArg{
						givenContext: context.Background(),
						givenKey:     "key",
						givenValue:   uint64(1),
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
				givenMode:       setModeNX,
			},
		},
		"float64": map[string]arg[float64]{
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
				givenMode:       setModeNone,
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
				givenMode:       setModeNone,
				expErr:          errors.New("error"),
			},
			"error: ErrFailToSetValue": {
				givenMockCmdArgFn: func() mockCmdArg {
					var cmd redis.StatusCmd
					cmd.SetVal("ok")
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
				givenMode:       setModeNX,
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
				givenMode:       setModeNX,
			},
		},
	}
	for dataType, tcs := range dataTypes {
		switch dataType {
		case "string":
			tcs := tcs.(map[string]arg[string])
			for scenario, tc := range tcs {
				tc := tc
				t.Run(fmt.Sprintf("[%s] %s", dataType, scenario), func(t *testing.T) {
					// Given
					mockCmd := new(mockredis.MockCmdable)

					// Mock
					tc.givenMockCmdArg = tc.givenMockCmdArgFn()
					mockCmd.ExpectedCalls = []*mock.Call{
						mockCmd.On(
							"SetArgs",
							tc.givenMockCmdArg.givenContext,
							tc.givenMockCmdArg.givenKey,
							tc.givenMockCmdArg.givenValue,
							tc.givenMockCmdArg.givenArgs,
						).Return(tc.givenMockCmdArg.expCmd),
					}

					// When
					err := setSingleValue(tc.givenContext, mockCmd, tc.givenKey, tc.givenValue, tc.givenExpiration, tc.givenMode)

					// Then
					if tc.expErr != nil {
						require.EqualError(t, err, tc.expErr.Error())
					} else {
						require.NoError(t, err)
					}
				})
			}
		case "int64":
			tcs := tcs.(map[string]arg[int64])
			for scenario, tc := range tcs {
				tc := tc
				t.Run(fmt.Sprintf("[%s] %s", dataType, scenario), func(t *testing.T) {
					// Given
					mockCmd := new(mockredis.MockCmdable)

					// Mock
					tc.givenMockCmdArg = tc.givenMockCmdArgFn()
					mockCmd.ExpectedCalls = []*mock.Call{
						mockCmd.On(
							"SetArgs",
							tc.givenMockCmdArg.givenContext,
							tc.givenMockCmdArg.givenKey,
							tc.givenMockCmdArg.givenValue,
							tc.givenMockCmdArg.givenArgs,
						).Return(tc.givenMockCmdArg.expCmd),
					}

					// When
					err := setSingleValue(tc.givenContext, mockCmd, tc.givenKey, tc.givenValue, tc.givenExpiration, tc.givenMode)

					// Then
					if tc.expErr != nil {
						require.EqualError(t, err, tc.expErr.Error())
					} else {
						require.NoError(t, err)
					}
				})
			}
		case "uint64":
			tcs := tcs.(map[string]arg[uint64])
			for scenario, tc := range tcs {
				tc := tc
				t.Run(fmt.Sprintf("[%s] %s", dataType, scenario), func(t *testing.T) {
					// Given
					mockCmd := new(mockredis.MockCmdable)

					// Mock
					tc.givenMockCmdArg = tc.givenMockCmdArgFn()
					mockCmd.ExpectedCalls = []*mock.Call{
						mockCmd.On(
							"SetArgs",
							tc.givenMockCmdArg.givenContext,
							tc.givenMockCmdArg.givenKey,
							tc.givenMockCmdArg.givenValue,
							tc.givenMockCmdArg.givenArgs,
						).Return(tc.givenMockCmdArg.expCmd),
					}

					// When
					err := setSingleValue(tc.givenContext, mockCmd, tc.givenKey, tc.givenValue, tc.givenExpiration, tc.givenMode)

					// Then
					if tc.expErr != nil {
						require.EqualError(t, err, tc.expErr.Error())
					} else {
						require.NoError(t, err)
					}
				})
			}
		case "float64":
			tcs := tcs.(map[string]arg[float64])
			for scenario, tc := range tcs {
				tc := tc
				t.Run(fmt.Sprintf("[%s] %s", dataType, scenario), func(t *testing.T) {
					// Given
					mockCmd := new(mockredis.MockCmdable)

					// Mock
					tc.givenMockCmdArg = tc.givenMockCmdArgFn()
					mockCmd.ExpectedCalls = []*mock.Call{
						mockCmd.On(
							"SetArgs",
							tc.givenMockCmdArg.givenContext,
							tc.givenMockCmdArg.givenKey,
							tc.givenMockCmdArg.givenValue,
							tc.givenMockCmdArg.givenArgs,
						).Return(tc.givenMockCmdArg.expCmd),
					}

					// When
					err := setSingleValue(tc.givenContext, mockCmd, tc.givenKey, tc.givenValue, tc.givenExpiration, tc.givenMode)

					// Then
					if tc.expErr != nil {
						require.EqualError(t, err, tc.expErr.Error())
					} else {
						require.NoError(t, err)
					}
				})
			}
		}
	}
}

func Test_getSingleValue(t *testing.T) {
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

	dataTypes := map[string]interface{}{
		"string": map[string]arg[string]{
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
		},
		"int64": map[string]arg[int64]{
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
		},
		"uint64": map[string]arg[uint64]{
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
		},
		"float64": map[string]arg[float64]{
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
		},
	}
	for dataType, tcs := range dataTypes {
		switch dataType {
		case "string":
			tcs := tcs.(map[string]arg[string])
			for scenario, tc := range tcs {
				tc := tc
				t.Run(fmt.Sprintf("[%s] %s", dataType, scenario), func(t *testing.T) {
					// Given
					mockCmd := new(mockredis.MockCmdable)

					// Mock
					tc.givenMockCmdArg = tc.givenMockCmdArgFn()
					mockCmd.ExpectedCalls = []*mock.Call{
						mockCmd.On(
							"Get",
							tc.givenMockCmdArg.givenContext,
							tc.givenMockCmdArg.givenKey,
						).Return(tc.givenMockCmdArg.expCmd),
					}

					// When
					result, err := getSingleValue[string](tc.givenContext, mockCmd, tc.givenKey)

					// Then
					if tc.expErr != nil {
						require.EqualError(t, err, tc.expErr.Error())
					} else {
						require.NoError(t, err)
						require.Equal(t, tc.expResult, result)
					}
				})
			}
		case "int64":
			tcs := tcs.(map[string]arg[int64])
			for scenario, tc := range tcs {
				tc := tc
				t.Run(fmt.Sprintf("[%s] %s", dataType, scenario), func(t *testing.T) {
					// Given
					mockCmd := new(mockredis.MockCmdable)

					// Mock
					tc.givenMockCmdArg = tc.givenMockCmdArgFn()
					mockCmd.ExpectedCalls = []*mock.Call{
						mockCmd.On(
							"Get",
							tc.givenMockCmdArg.givenContext,
							tc.givenMockCmdArg.givenKey,
						).Return(tc.givenMockCmdArg.expCmd),
					}

					// When
					result, err := getSingleValue[int64](tc.givenContext, mockCmd, tc.givenKey)

					// Then
					if tc.expErr != nil {
						require.EqualError(t, err, tc.expErr.Error())
					} else {
						require.NoError(t, err)
						require.Equal(t, tc.expResult, result)
					}
				})
			}
		case "uint64":
			tcs := tcs.(map[string]arg[uint64])
			for scenario, tc := range tcs {
				tc := tc
				t.Run(fmt.Sprintf("[%s] %s", dataType, scenario), func(t *testing.T) {
					// Given
					mockCmd := new(mockredis.MockCmdable)

					// Mock
					tc.givenMockCmdArg = tc.givenMockCmdArgFn()
					mockCmd.ExpectedCalls = []*mock.Call{
						mockCmd.On(
							"Get",
							tc.givenMockCmdArg.givenContext,
							tc.givenMockCmdArg.givenKey,
						).Return(tc.givenMockCmdArg.expCmd),
					}

					// When
					result, err := getSingleValue[uint64](tc.givenContext, mockCmd, tc.givenKey)

					// Then
					if tc.expErr != nil {
						require.EqualError(t, err, tc.expErr.Error())
					} else {
						require.NoError(t, err)
						require.Equal(t, tc.expResult, result)
					}
				})
			}
		case "float64":
			tcs := tcs.(map[string]arg[float64])
			for scenario, tc := range tcs {
				tc := tc
				t.Run(fmt.Sprintf("[%s] %s", dataType, scenario), func(t *testing.T) {
					// Given
					mockCmd := new(mockredis.MockCmdable)

					// Mock
					tc.givenMockCmdArg = tc.givenMockCmdArgFn()
					mockCmd.ExpectedCalls = []*mock.Call{
						mockCmd.On(
							"Get",
							tc.givenMockCmdArg.givenContext,
							tc.givenMockCmdArg.givenKey,
						).Return(tc.givenMockCmdArg.expCmd),
					}

					// When
					result, err := getSingleValue[float64](tc.givenContext, mockCmd, tc.givenKey)

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
	}
}

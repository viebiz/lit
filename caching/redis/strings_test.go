package redis

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func Test_setSingleValue(t *testing.T) {
	type mockCmdArg struct {
		givenContext    context.Context
		givenKey        string
		givenValue      interface{}
		givenExpiration time.Duration
		givenArgs       redis.SetArgs
		expCmd          redis.StatusCmd
	}

	type arg[T Type] struct {
		givenMockCmdArgFn func() mockCmdArg
		givenMockCmdArg   mockCmdArg
		givenContext      context.Context
		givenCmd          redis.Cmdable
		givenKey          string
		givenValue        T
		givenExpiration   time.Duration
		givenMode         setMode
		expErr            error
	}

	dataTypes := map[string]interface{}{
		"string":  map[string]arg[string]{},
		"int64":   map[string]arg[int64]{},
		"uint64":  map[string]arg[uint64]{},
		"float64": map[string]arg[float64]{},
	}
	for dataType, tcs := range dataTypes {
		switch dataType {
		case "string":
			for scenario, tc := range tcs {
				tc := tc.(arg[string])
				t.Run(fmt.Sprintf("[%s] %s", dataType, scenario), func(t *testing.T) {
					// Given

					// Mock

					// When

					// Then
				})
			}
		case "int64":
			for scenario, tc := range tcs {
				tc := tc.(arg[int64])
				t.Run(fmt.Sprintf("[%s] %s", dataType, scenario), func(t *testing.T) {
					// Given

					// Mock

					// When

					// Then
				})
			}
		case "uint64":
			for scenario, tc := range tcs {
				tc := tc.(arg[uint64])
				t.Run(fmt.Sprintf("[%s] %s", dataType, scenario), func(t *testing.T) {
					// Given

					// Mock

					// When

					// Then
				})
			}
		case "float64":
			for scenario, tc := range tcs {
				tc := tc.(arg[float64])
				t.Run(fmt.Sprintf("[%s] %s", dataType, scenario), func(t *testing.T) {
					// Given

					// Mock

					// When

					// Then
				})
			}
		}

	}
}

package snowflake

import (
	"fmt"
	"testing"
	"time"

	pkgerrors "github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	future := time.Now().Add(time.Hour)

	type arg struct {
		givenOpts []Option
		expErr    error
		expResult bool
	}
	tcs := map[string]arg{
		"success no ops": {
			expResult: true,
		},
		"success start time": {
			givenOpts: []Option{StartTime(time.Now())},
			expResult: true,
		},
		"success machine id": {
			givenOpts: []Option{MachineID(12345)},
			expResult: true,
		},
		"success start time+machine id": {
			givenOpts: []Option{StartTime(time.Now()), MachineID(12345)},
			expResult: true,
		},
		"err start time": {
			givenOpts: []Option{StartTime(future)},
			expErr:    fmt.Errorf("invalid start time provided: %s", future),
		},
		"err machine id": {
			givenOpts: []Option{MachineID(0)},
			expErr:    fmt.Errorf("invalid machine ID provided: %d", 0),
		},
		"err start time+machine id": {
			givenOpts: []Option{StartTime(future), MachineID(0)},
			expErr:    fmt.Errorf("invalid start time provided: %s", future),
		},
	}
	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given & When:
			g, err := New(tc.givenOpts...)

			// Then:
			require.Equal(t, tc.expErr, pkgerrors.Cause(err))
			require.Equal(t, tc.expResult, g != nil)
		})
	}
}

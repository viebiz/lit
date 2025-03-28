package aes

import (
	"crypto/rand"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	type arg struct {
		mockRandReadFn func([]byte) (int, error)
		expErr         error
	}
	tcs := map[string]arg{
		"err": {
			mockRandReadFn: func([]byte) (int, error) {
				return 0, errors.New("err")
			},
			expErr: errors.New("err"),
		},
		"ok": {
			mockRandReadFn: rand.Read,
		},
	}
	for scenario, tc := range tcs {
		t.Run(scenario, func(t *testing.T) {
			// Given
			randReadFn = tc.mockRandReadFn
			defer func() {
				randReadFn = rand.Read
			}()

			// When
			instance, err := New()

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

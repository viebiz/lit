package redis

import (
	"crypto/tls"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
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

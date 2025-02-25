package redis

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/viebiz/lit/monitoring"
)

func Test_newTracingHook(t *testing.T) {
	type arg struct {
		givenInfo monitoring.ExternalServiceInfo
		expResult tracingHook
	}
	tcs := map[string]arg{
		"ok": {
			givenInfo: monitoring.ExternalServiceInfo{
				Hostname: "localhost",
				Port:     "3000",
			},
			expResult: tracingHook{
				info: monitoring.ExternalServiceInfo{
					Hostname: "localhost",
					Port:     "3000",
				},
			},
		},
	}
	for scenario, tc := range tcs {
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			// Given

			// When
			result := newTracingHook(tc.givenInfo)

			// Then
			require.Equal(t, tc.expResult, result)
		})
	}
}

func Test_tracingHook_DialHook(t *testing.T) {
	type arg struct {
		givenInfo monitoring.ExternalServiceInfo
		expResult tracingHook
	}
	tcs := map[string]arg{
		"ok": {
			givenInfo: monitoring.ExternalServiceInfo{
				Hostname: "localhost",
				Port:     "3000",
			},
			expResult: tracingHook{
				info: monitoring.ExternalServiceInfo{
					Hostname: "localhost",
					Port:     "3000",
				},
			},
		},
	}
	for scenario, tc := range tcs {
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			// Given

			// When
			instance := tracingHook{}
			result := instance.DialHook(tc.givenInfo)

			// Then
			require.Equal(t, tc.expResult, result)
		})
	}
}

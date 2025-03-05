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
				Hostname: "redis",
				Port:     "6379",
			},
			expResult: tracingHook{
				info: monitoring.ExternalServiceInfo{
					Hostname: "redis",
					Port:     "6379",
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

//func Test_tracingHook_DialHook(t *testing.T) {
//	type arg struct {
//		givenFn   func(ctx context.Context, network string, addr string) (net.Conn, error)
//		expResult tracingHook
//	}
//	tcs := map[string]arg{
//		"ok": {
//			givenFn: func(ctx context.Context, network string, addr string) (net.Conn, error) {
//				return nil, nil
//			},
//			expResult: tracingHook{
//				info: monitoring.ExternalServiceInfo{
//					Hostname: "redis",
//					Port:     "6379",
//				},
//			},
//		},
//	}
//	for scenario, tc := range tcs {
//		t.Run(scenario, func(t *testing.T) {
//			t.Parallel()
//
//			// Given
//
//			// When
//			instance := tracingHook{}
//			result := instance.DialHook(tc.givenFn)
//
//			// Then
//			require.Equal(t, tc.expResult, result)
//		})
//	}
//}

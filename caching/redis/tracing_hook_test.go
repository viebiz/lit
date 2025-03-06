package redis

import (
	"context"
	"errors"
	"net"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/viebiz/lit/mocks/mocknet"
	"github.com/viebiz/lit/mocks/mockredis"
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

func Test_tracingHook_DialHook(t *testing.T) {
	type mockDialHookArg struct {
		givenContext context.Context
		givenNetwork string
		givenAddress string
		expConn      *mocknet.MockConn
		expErr       error
	}
	type arg struct {
		givenDialHook        func(ctx context.Context, network string, addr string) (net.Conn, error)
		givenMockDialHookArg mockDialHookArg
		expConn              *mocknet.MockConn
		expErr               error
	}
	tcs := map[string]arg{
		"error": {
			givenDialHook: func(ctx context.Context, network string, addr string) (net.Conn, error) {
				return nil, errors.New("error")
			},
			givenMockDialHookArg: mockDialHookArg{
				givenContext: context.Background(),
				givenNetwork: "tcp",
				givenAddress: "localhost:6379",
				expErr:       errors.New("error"),
			},
			expErr: errors.New("error"),
		},
		"success": {
			givenDialHook: func(ctx context.Context, network string, addr string) (net.Conn, error) {
				return new(mocknet.MockConn), nil
			},
			givenMockDialHookArg: mockDialHookArg{
				givenContext: context.Background(),
				givenNetwork: "tcp",
				givenAddress: "localhost:6379",
				expConn:      &mocknet.MockConn{},
			},
			expConn: &mocknet.MockConn{},
		},
	}
	for scenario, tc := range tcs {
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()
			// Given

			// Mock
			mockDialHook := mockredis.MockDialHook{}
			mockDialHook.ExpectedCalls = []*mock.Call{
				mockDialHook.On(
					"Execute",
					tc.givenMockDialHookArg.givenContext,
					tc.givenMockDialHookArg.givenNetwork,
					tc.givenMockDialHookArg.givenAddress,
				).Return(tc.givenMockDialHookArg.expConn, tc.givenMockDialHookArg.expErr),
			}

			// When
			instance := tracingHook{}
			hookFn := instance.DialHook(tc.givenDialHook)
			result, err := hookFn(tc.givenMockDialHookArg.givenContext, tc.givenMockDialHookArg.givenNetwork, tc.givenMockDialHookArg.givenAddress)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expConn, result)
			}
		})
	}
}

func Test_tracingHook_ProcessHook(t *testing.T) {
	//type mockDialHookArg struct {
	//	givenContext context.Context
	//	givenNetwork string
	//	givenAddress string
	//	expConn      *mocknet.MockConn
	//	expErr       error
	//}
	//type arg struct {
	//	givenDialHook        func(ctx context.Context, network string, addr string) (net.Conn, error)
	//	givenMockDialHookArg mockDialHookArg
	//	expConn              *mocknet.MockConn
	//	expErr               error
	//}
	//tcs := map[string]arg{
	//	"error": {
	//		givenDialHook: func(ctx context.Context, network string, addr string) (net.Conn, error) {
	//			return nil, errors.New("error")
	//		},
	//		givenMockDialHookArg: mockDialHookArg{
	//			givenContext: context.Background(),
	//			givenNetwork: "tcp",
	//			givenAddress: "localhost:6379",
	//			expErr:       errors.New("error"),
	//		},
	//		expErr: errors.New("error"),
	//	},
	//	"success": {
	//		givenDialHook: func(ctx context.Context, network string, addr string) (net.Conn, error) {
	//			return new(mocknet.MockConn), nil
	//		},
	//		givenMockDialHookArg: mockDialHookArg{
	//			givenContext: context.Background(),
	//			givenNetwork: "tcp",
	//			givenAddress: "localhost:6379",
	//			expConn:      &mocknet.MockConn{},
	//		},
	//		expConn: &mocknet.MockConn{},
	//	},
	//}
	//for scenario, tc := range tcs {
	//	t.Run(scenario, func(t *testing.T) {
	//		t.Parallel()
	//		// Given
	//
	//		// Mock
	//		mockDialHook := mockredis.MockDialHook{}
	//		mockDialHook.ExpectedCalls = []*mock.Call{
	//			mockDialHook.On(
	//				"Execute",
	//				tc.givenMockDialHookArg.givenContext,
	//				tc.givenMockDialHookArg.givenNetwork,
	//				tc.givenMockDialHookArg.givenAddress,
	//			).Return(tc.givenMockDialHookArg.expConn, tc.givenMockDialHookArg.expErr),
	//		}
	//
	//		// When
	//		instance := tracingHook{}
	//		hookFn := instance.DialHook(tc.givenDialHook)
	//		result, err := hookFn(tc.givenMockDialHookArg.givenContext, tc.givenMockDialHookArg.givenNetwork, tc.givenMockDialHookArg.givenAddress)
	//
	//		// Then
	//		if tc.expErr != nil {
	//			require.EqualError(t, err, tc.expErr.Error())
	//		} else {
	//			require.NoError(t, err)
	//			require.Equal(t, tc.expConn, result)
	//		}
	//	})
	//}
}

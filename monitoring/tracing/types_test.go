package tracing

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransportType_String(t *testing.T) {
	type arg struct {
		givenTransportType TransportType
		expString          string
	}
	tcs := map[string]arg{
		"grpc transport": {
			givenTransportType: TransportGRPC,
			expString:          "grpc",
		},
		"http transport": {
			givenTransportType: TransportHTTP,
			expString:          "http",
		},
		"custom transport": {
			givenTransportType: TransportType("custom"),
			expString:          "custom",
		},
		"empty transport": {
			givenTransportType: TransportType(""),
			expString:          "",
		},
	}

	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			// Given
			// When
			got := tc.givenTransportType.String()

			// Then
			require.Equal(t, tc.expString, got)
		})
	}
}

func TestTransportType_IsValid(t *testing.T) {
	type arg struct {
		givenTransportType TransportType
		expValid           bool
	}
	tcs := map[string]arg{
		"grpc is valid": {
			givenTransportType: TransportGRPC,
			expValid:           true,
		},
		"http is valid": {
			givenTransportType: TransportHTTP,
			expValid:           true,
		},
		"uppercase grpc is invalid": {
			givenTransportType: TransportType("GRPC"),
			expValid:           false,
		},
		"uppercase http is invalid": {
			givenTransportType: TransportType("HTTP"),
			expValid:           false,
		},
		"invalid transport type": {
			givenTransportType: TransportType("invalid"),
			expValid:           false,
		},
		"empty transport is invalid": {
			givenTransportType: TransportType(""),
			expValid:           false,
		},
		"tcp transport is invalid": {
			givenTransportType: TransportType("tcp"),
			expValid:           false,
		},
		"websocket transport is invalid": {
			givenTransportType: TransportType("websocket"),
			expValid:           false,
		},
		"mixed case grpc is invalid": {
			givenTransportType: TransportType("Grpc"),
			expValid:           false,
		},
		"mixed case http is invalid": {
			givenTransportType: TransportType("Http"),
			expValid:           false,
		},
	}

	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			// Given
			// When
			got := tc.givenTransportType.IsValid()

			// Then
			require.Equal(t, tc.expValid, got)
		})
	}
}

func TestTransportType_Constants(t *testing.T) {
	// Verify the constant values are as expected
	require.Equal(t, "grpc", string(TransportGRPC))
	require.Equal(t, "http", string(TransportHTTP))

	// Verify constants are valid
	require.True(t, TransportGRPC.IsValid())
	require.True(t, TransportHTTP.IsValid())
}

func TestTransportType_StringRoundTrip(t *testing.T) {
	// Given
	transportTypes := []TransportType{
		TransportGRPC,
		TransportHTTP,
	}

	for _, tt := range transportTypes {
		t.Run(string(tt), func(t *testing.T) {
			// When
			str := tt.String()
			reconstructed := TransportType(str)

			// Then
			require.Equal(t, tt, reconstructed)
			require.True(t, reconstructed.IsValid())
		})
	}
}

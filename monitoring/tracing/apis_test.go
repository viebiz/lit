package tracing

import (
	"context"
	"crypto/tls"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	type arg struct {
		givenCfg         Config
		contextCancelled bool
		expErr           error
	}
	tcs := map[string]arg{
		"success - grpc": {
			givenCfg: Config{
				ExporterURL: "localhost:4317",
			},
		},
		"success - grpc with tls": {
			givenCfg: Config{
				ExporterURL: "localhost:4317",
				TLSConfig:   &tls.Config{},
			},
		},
		"success - http": {
			givenCfg: Config{
				ExporterURL:   "localhost:4318",
				TransportType: TransportHTTP,
			},
		},
		"success - http with tls": {
			givenCfg: Config{
				ExporterURL:   "localhost:4318",
				TransportType: TransportHTTP,
				TLSConfig:     &tls.Config{},
			},
		},
		"error - missing URL": {
			expErr: ErrMissingExporterURL,
		},
		"error - invalid transport type": {
			givenCfg: Config{
				ExporterURL:   "localhost:4318",
				TransportType: "INVALID",
			},
			expErr: fmt.Errorf("%w: INVALID", ErrInvalidTransportType),
		},
	}
	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()
			// Given

			// When
			err := Init(context.Background(), tc.givenCfg)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

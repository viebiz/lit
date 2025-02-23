package tracing

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_createExporter(t *testing.T) {
	type args struct {
		givenCfg Config
		expErr   error
	}
	tcs := map[string]args{
		"success": {
			givenCfg: Config{
				ExporterURL:   "localhost:4317",
				TransportType: TransportGRPC,
			},
		},
		"error": {
			givenCfg: Config{
				ExporterURL:   "localhost:4317",
				TransportType: "unsupported",
			},
			expErr: errors.New("unknown transport type: unsupported"),
		},
	}

	for scenario, tc := range tcs {
		t.Run(scenario, func(t *testing.T) {
			tc := tc
			t.Run(scenario, func(t *testing.T) {
				t.Parallel()
				// Given

				// When
				exp, err := createExporter(context.Background(), tc.givenCfg)

				// Then
				if tc.expErr != nil {
					require.EqualError(t, err, tc.expErr.Error())
				} else {
					require.NoError(t, err)
					require.NotNil(t, exp)
				}
			})
		})
	}
}

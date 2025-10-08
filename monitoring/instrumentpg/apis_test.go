package instrumentpg

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/viebiz/lit/mocks/mocktracer"
	"github.com/viebiz/lit/postgres"
)

func TestWithInstrumentation(t *testing.T) {
	tp := mocktracer.NewTestTracerProvider(t)
	defer tp.Shutdown(t)

	tcs := map[string]struct {
		givenPool postgres.BeginnerExecutor
	}{
		"nil pool": {
			givenPool: nil,
		},
	}

	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			// When
			instrumented := WithInstrumentation(tc.givenPool)

			// Then
			require.NotNil(t, instrumented)
			// Should return instrumentedDB type
			_, ok := instrumented.(instrumentedDB)
			require.True(t, ok)
		})
	}
}

func TestWithInstrumentationTx(t *testing.T) {
	tp := mocktracer.NewTestTracerProvider(t)
	defer tp.Shutdown(t)

	tcs := map[string]struct {
		givenTx postgres.ContextExecutor
	}{
		"nil tx": {
			givenTx: nil,
		},
	}

	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			// When
			instrumented := WithInstrumentationTx(tc.givenTx)

			// Then
			require.NotNil(t, instrumented)
			// Should return instrumentedTx type
			_, ok := instrumented.(instrumentedTx)
			require.True(t, ok)
		})
	}
}

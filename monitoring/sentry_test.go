package monitoring

import (
	"sync"
	"testing"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func Test_initSentry(t *testing.T) {
	type args struct {
		givenCfg    sentryConfig
		clientExist bool
		expErr      error
	}
	tcs := map[string]args{
		"success": {
			givenCfg: sentryConfig{
				DSN:         "https://whatever@example.com/16042000",
				ServerName:  "lightning",
				Environment: "dev",
				Version:     "1.0.0",
			},
			clientExist: true,
		},
		"success - skip if DSN empty": {
			givenCfg: sentryConfig{},
		},
	}

	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()
			// Given

			// When
			c, err := initSentry(tc.givenCfg, zap.NewNop())

			// Then
			if tc.expErr != nil {
				require.ErrorContains(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
			}
			if tc.clientExist {
				require.NotNil(t, c)
			}
		})
	}
}

// TransportMock is a mock implementation of the HTTP transport for the Sentry client.
// It stores captured events for testing purposes.
//
//	Keep it here to avoid impact test coverage
type TransportMock struct {
	mu        sync.Mutex
	events    []*sentry.Event
	lastEvent *sentry.Event
}

func (t *TransportMock) Configure(_ sentry.ClientOptions) {}
func (t *TransportMock) SendEvent(event *sentry.Event) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.events = append(t.events, event)
	t.lastEvent = event
}
func (t *TransportMock) Flush(_ time.Duration) bool {
	return true
}
func (t *TransportMock) Events() []*sentry.Event {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.events
}
func (t *TransportMock) Close() {}

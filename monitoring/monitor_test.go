package monitoring

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/getsentry/sentry-go"
	"github.com/google/go-cmp/cmp/cmpopts"
	pkgerrors "github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/viebiz/lit/internal/testutil"
)

func TestMonitorLogger(t *testing.T) {
	type args struct {
		doLogging func(w io.Writer)
		expected  []map[string]string
	}
	tcs := map[string]args{
		"infof": {
			doLogging: func(w io.Writer) {
				m, err := New(Config{ServerName: "lightning", Environment: "dev", Version: "1.0.0", Writer: w})
				require.NoError(t, err)

				m.Infof("Hello %s project", "lightning")
				m.Flush(DefaultFlushWait)
			},
			expected: []map[string]string{
				{"level": "INFO", "ts": "2025-02-23T13:34:56.185+0700", "msg": "Sentry DSN not provided. Not using Sentry Error Reporting", "server.name": "lightning", "environment": "dev", "version": "1.0.0"},
				{"level": "INFO", "ts": "2025-02-23T13:34:56.185+0700", "msg": "OTelExporter URL not provided. Not using Distributed Tracing", "server.name": "lightning", "environment": "dev", "version": "1.0.0"},
				{"level": "INFO", "ts": "2025-02-23T13:34:56.185+0700", "msg": "Hello lightning project", "server.name": "lightning", "environment": "dev", "version": "1.0.0"},
			},
		},
		"errorf": {
			doLogging: func(w io.Writer) {
				m, err := New(Config{ServerName: "lightning", Environment: "dev", Version: "1.0.0", Writer: w})
				require.NoError(t, err)

				m.Errorf(errors.New("simulated error for unit test"), "")
				m.Flush(DefaultFlushWait)
			},
			expected: []map[string]string{
				{"level": "INFO", "ts": "2025-02-23T13:34:56.185+0700", "msg": "Sentry DSN not provided. Not using Sentry Error Reporting", "server.name": "lightning", "environment": "dev", "version": "1.0.0"},
				{"level": "INFO", "ts": "2025-02-23T13:34:56.185+0700", "msg": "OTelExporter URL not provided. Not using Distributed Tracing", "server.name": "lightning", "environment": "dev", "version": "1.0.0"},
				{"level": "ERROR", "ts": "2025-02-23T13:34:56.185+0700", "msg": "Err: simulated error for unit test", "error.kind": "*errors.errorString", "error.message": "simulated error for unit test", "server.name": "lightning", "environment": "dev", "version": "1.0.0"},
			},
		},
		"errorf - with extra message": {
			doLogging: func(w io.Writer) {
				m, err := New(Config{ServerName: "lightning", Environment: "dev", Version: "1.0.0", Writer: w})
				require.NoError(t, err)

				m.WithTag("request_id", "123").Errorf(errors.New("simulated error"), "Unit test exception on %s project", "lightning")
				m.Flush(DefaultFlushWait)
			},
			expected: []map[string]string{
				{"level": "INFO", "ts": "2025-02-23T13:34:56.185+0700", "msg": "Sentry DSN not provided. Not using Sentry Error Reporting", "server.name": "lightning", "environment": "dev", "version": "1.0.0"},
				{"level": "INFO", "ts": "2025-02-23T13:34:56.185+0700", "msg": "OTelExporter URL not provided. Not using Distributed Tracing", "server.name": "lightning", "environment": "dev", "version": "1.0.0"},
				{"level": "ERROR", "ts": "2025-02-23T13:34:56.185+0700", "msg": "Unit test exception on lightning project. Err: simulated error", "error.kind": "*errors.errorString", "error.message": "simulated error", "server.name": "lightning", "environment": "dev", "version": "1.0.0", "request_id": "123"},
			},
		},
		"errorf - with stacktrace error": {
			doLogging: func(w io.Writer) {
				m, err := New(Config{ServerName: "lightning", Environment: "dev", Version: "1.0.0", Writer: w})
				require.NoError(t, err)

				m.WithTag("request_id", "123").Errorf(pkgerrors.WithStack(errors.New("simulated error")), "Unit test exception on %s project", "lightning")
				m.Flush(DefaultFlushWait)
			},
			expected: []map[string]string{
				{"level": "INFO", "ts": "2025-02-23T13:34:56.185+0700", "msg": "Sentry DSN not provided. Not using Sentry Error Reporting", "server.name": "lightning", "environment": "dev", "version": "1.0.0"},
				{"level": "INFO", "ts": "2025-02-23T13:34:56.185+0700", "msg": "OTelExporter URL not provided. Not using Distributed Tracing", "server.name": "lightning", "environment": "dev", "version": "1.0.0"},
				{"level": "ERROR", "ts": "2025-02-23T13:34:56.185+0700", "msg": "Unit test exception on lightning project. Err: simulated error", "error.kind": "*errors.withStack", "error.message": "simulated error", "error.stack": "github.com/viebiz/lit/monitoring.TestMonitorLogger.func", "server.name": "lightning", "environment": "dev", "version": "1.0.0", "request_id": "123"},
			},
		},
		"infof with tags - child does not affect parent": {
			doLogging: func(w io.Writer) {
				m, err := New(Config{ServerName: "lightning", Environment: "dev", Version: "1.0.0", Writer: w})
				require.NoError(t, err)

				m.WithTag("request_id", "123").Infof("Child Monitor hello world")
				m.Infof("Parent Monitor hello world")
				m.Flush(DefaultFlushWait)
			},
			expected: []map[string]string{
				{"level": "INFO", "ts": "2025-02-23T13:34:56.185+0700", "msg": "Sentry DSN not provided. Not using Sentry Error Reporting", "server.name": "lightning", "environment": "dev", "version": "1.0.0"},
				{"level": "INFO", "ts": "2025-02-23T13:34:56.185+0700", "msg": "OTelExporter URL not provided. Not using Distributed Tracing", "server.name": "lightning", "environment": "dev", "version": "1.0.0"},
				{"level": "INFO", "ts": "2025-02-23T13:34:56.185+0700", "msg": "Child Monitor hello world", "server.name": "lightning", "environment": "dev", "version": "1.0.0", "request_id": "123"},
				{"level": "INFO", "ts": "2025-02-23T13:34:56.185+0700", "msg": "Parent Monitor hello world", "server.name": "lightning", "environment": "dev", "version": "1.0.0"},
			},
		},
		"errorf with tags - child does not affect parent": {
			doLogging: func(w io.Writer) {
				m, err := New(Config{ServerName: "lightning", Environment: "dev", Version: "1.0.0", Writer: w})
				require.NoError(t, err)

				simulatedErr := errors.New("simulated error for unit testing")
				m.WithTag("request_id", "123").Errorf(simulatedErr, "Child")
				m.Errorf(simulatedErr, "Parent")
				m.Flush(DefaultFlushWait)
			},
			expected: []map[string]string{
				{"level": "INFO", "ts": "2025-02-23T13:34:56.185+0700", "msg": "Sentry DSN not provided. Not using Sentry Error Reporting", "server.name": "lightning", "environment": "dev", "version": "1.0.0"},
				{"level": "INFO", "ts": "2025-02-23T13:34:56.185+0700", "msg": "OTelExporter URL not provided. Not using Distributed Tracing", "server.name": "lightning", "environment": "dev", "version": "1.0.0"},
				{"level": "ERROR", "ts": "2025-02-23T13:34:56.185+0700", "msg": "Child. Err: simulated error for unit testing", "error.kind": "*errors.errorString", "error.message": "simulated error for unit testing", "server.name": "lightning", "environment": "dev", "version": "1.0.0", "request_id": "123"},
				{"level": "ERROR", "ts": "2025-02-23T13:34:56.185+0700", "msg": "Parent. Err: simulated error for unit testing", "error.kind": "*errors.errorString", "error.message": "simulated error for unit testing", "server.name": "lightning", "environment": "dev", "version": "1.0.0"},
			},
		},
		"infof - not in ctx": {
			doLogging: func(w io.Writer) {
				m := FromContext(context.Background())

				m.Infof("Skip this log because it is nil")
				m.Flush(DefaultFlushWait)
			},
		},
		"infof with tags - not in ctx": {
			doLogging: func(w io.Writer) {
				m := FromContext(context.Background())

				m.WithTag("name", "lightning").Infof("Skip this log because it is nil")
				m.Flush(DefaultFlushWait)
			},
		},
		"errorf - not in ctx": {
			doLogging: func(w io.Writer) {
				m := FromContext(context.Background())

				m.Errorf(errors.New("simulated error"), "Skip this log because it is nil")
				m.Flush(DefaultFlushWait)
			},
		},
	}

	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			// Given
			logBuffer := bytes.NewBuffer(nil)

			// When
			tc.doLogging(logBuffer)

			// Then
			parsedLog, err := parseLog(logBuffer.Bytes())
			require.NoError(t, err)
			testutil.Equal(t, tc.expected, parsedLog, cmpopts.IgnoreMapEntries(func(k string, v any) bool {
				// Ignore timestamp field as it updates dynamically
				if k == "ts" {
					return true
				}

				// The error stack quite big, so ignore it first
				if k == "error.stack" {
					return true
				}

				return false
			}))
		})
	}
}

func TestMonitor_ReportError(t *testing.T) {
	type args struct {
		givenErr  error
		givenMsg  string
		givenArgs []any
		expected  []map[string]string
	}
	tcs := map[string]args{
		"errorf - capture error": {
			givenErr: errors.New("simulated error"),
			expected: []map[string]string{
				{"level": "INFO", "ts": "2025-02-23T13:34:56.185+0700", "msg": "Sentry DSN not provided. Not using Sentry Error Reporting", "server.name": "lightning", "environment": "dev", "version": "1.0.0"},
				{"level": "INFO", "ts": "2025-02-23T13:34:56.185+0700", "msg": "OTelExporter URL not provided. Not using Distributed Tracing", "server.name": "lightning", "environment": "dev", "version": "1.0.0"},
				{"level": "ERROR", "ts": "2025-02-23T13:34:56.185+0700", "msg": "Err: simulated error", "error.kind": "*errors.errorString", "error.message": "simulated error", "server.name": "lightning", "environment": "dev", "version": "1.0.0"},
			},
		},
	}

	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()
			// Given
			logBuffer := bytes.NewBuffer(nil)
			m, err := New(Config{ServerName: "lightning", Environment: "dev", Version: "1.0.0", Writer: logBuffer})
			defer m.Flush(DefaultFlushWait)
			require.NoError(t, err)

			sentryClient, transport := setupClientTest()
			m.sentryClient = sentryClient // Inject sentryClient to capture error report

			// When
			m.Errorf(tc.givenErr, tc.givenMsg, tc.givenArgs...)

			// Then
			parsedLog, err := parseLog(logBuffer.Bytes())
			require.NoError(t, err)
			testutil.Equal(t, tc.expected, parsedLog, cmpopts.IgnoreMapEntries(func(k string, v any) bool {
				// Ignore timestamp field as it updates dynamically
				if k == "ts" {
					return true
				}

				// The error stack quite big, so ignore it first
				if k == "error.stack" {
					return true
				}

				return false
			}))
			exceptions := transport.lastEvent.Exception
			require.True(t, len(exceptions) > 0)
			lastException := exceptions[len(exceptions)-1]
			require.Equal(t, lastException.Type, "*errors.errorString")
			require.Equal(t, lastException.Value, "simulated error")
			require.True(t, len(lastException.Stacktrace.Frames) > 0)
		})
	}
}

func setupClientTest() (*sentry.Client, *TransportMock) {
	transport := &TransportMock{}
	client, _ := sentry.NewClient(sentry.ClientOptions{
		Dsn:       "https://whatever@example.com/16042000",
		Transport: transport,
		Integrations: func(i []sentry.Integration) []sentry.Integration {
			return []sentry.Integration{}
		},
	})

	return client, transport
}

// parseLog first converts []byte into string and then to map.
// Idea is to mimic actual log line of key value pairs
func parseLog(b []byte) ([]map[string]string, error) {
	var result []map[string]string
	for _, s := range strings.Split(string(b), "\n") {
		if s == "" {
			break
		}
		var r map[string]string
		if err := json.Unmarshal([]byte(s), &r); err != nil {
			if strings.HasSuffix(s, "Initializing Logger") {
				continue
			}
			return nil, err
		}
		result = append(result, r)
	}
	return result, nil
}

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
	pkgerrors "github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"github.com/viebiz/lit/testutil"
)

func TestMonitorLogger(t *testing.T) {
	type args struct {
		doLogging func(w io.Writer)
		expected  []map[string]interface{}
	}
	tcs := map[string]args{
		"infof": {
			doLogging: func(w io.Writer) {
				m, err := New(Config{ServerName: "lightning", Environment: "dev", Version: "1.0.0", Writer: w})
				require.NoError(t, err)

				m.Infof("Hello %s project", "lightning")
				m.Flush(DefaultFlushWait)
			},
			expected: []map[string]interface{}{
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
			expected: []map[string]interface{}{
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
			expected: []map[string]interface{}{
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
			expected: []map[string]interface{}{
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
			expected: []map[string]interface{}{
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
			expected: []map[string]interface{}{
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
		"info": {
			doLogging: func(w io.Writer) {
				m, err := New(Config{ServerName: "lightning", Environment: "dev", Version: "1.0.0", Writer: w})
				require.NoError(t, err)

				m.Info("Hello lightning project")
				m.Flush(DefaultFlushWait)
			},
			expected: []map[string]interface{}{
				{"level": "INFO", "ts": "2025-02-23T13:34:56.185+0700", "msg": "Hello lightning project", "server.name": "lightning", "environment": "dev", "version": "1.0.0"},
			},
		},
		"info with fields": {
			doLogging: func(w io.Writer) {
				m, err := New(Config{ServerName: "lightning", Environment: "dev", Version: "1.0.0", Writer: w})
				require.NoError(t, err)

				m.Info("Hello", StringField("name", "lightning"), IntField("release", 20250710), JSONField("context", []byte(`{"dev":"the-witcher-knight","role":"admin"}`)))
				m.Flush(DefaultFlushWait)
			},
			expected: []map[string]interface{}{
				{"level": "INFO", "ts": "2025-02-23T13:34:56.185+0700", "msg": "Hello", "name": "lightning", "release": 20250710.0, "context": map[string]any{"dev": "the-witcher-knight", "role": "admin"}, "server.name": "lightning", "environment": "dev", "version": "1.0.0"},
			},
		},
		"info - not in ctx": {
			doLogging: func(w io.Writer) {
				m := FromContext(context.Background())

				m.Info("Skip this log because it is nil")
				m.Flush(DefaultFlushWait)
			},
		},
		"info - nil monitor": {
			doLogging: func(w io.Writer) {
				var m *Monitor
				m.Info("Skip this log because it is nil")
			},
		},
		"infof - nil monitor": {
			doLogging: func(w io.Writer) {
				var m *Monitor
				m.Infof("Skip this log because it is nil")
			},
		},
		"error - nil monitor": {
			doLogging: func(w io.Writer) {
				var m *Monitor
				m.Error(errors.New("test error"), "Skip this log because it is nil")
			},
		},
		"errorf - nil monitor": {
			doLogging: func(w io.Writer) {
				var m *Monitor
				m.Errorf(errors.New("test error"), "Skip this log because it is nil")
			},
		},
		"error": {
			doLogging: func(w io.Writer) {
				m, err := New(Config{ServerName: "lightning", Environment: "dev", Version: "1.0.0", Writer: w})
				require.NoError(t, err)

				m.Error(errors.New("simulated error for unit test"), "", StringField("name", "lightning"))
				m.Flush(DefaultFlushWait)
			},
			expected: []map[string]interface{}{
				{"level": "ERROR", "ts": "2025-02-23T13:34:56.185+0700", "msg": "Err: simulated error for unit test", "error.kind": "*errors.errorString", "error.message": "simulated error for unit test", "name": "lightning", "server.name": "lightning", "environment": "dev", "version": "1.0.0"},
			},
		},
		"error - with extra message": {
			doLogging: func(w io.Writer) {
				m, err := New(Config{ServerName: "lightning", Environment: "dev", Version: "1.0.0", Writer: w})
				require.NoError(t, err)

				m.WithTag("request_id", "123").Error(errors.New("simulated error"), "Unit test exception", StringField("name", "lightning"))
				m.Flush(DefaultFlushWait)
			},
			expected: []map[string]interface{}{
				{"level": "ERROR", "ts": "2025-02-23T13:34:56.185+0700", "msg": "Unit test exception. Err: simulated error", "error.kind": "*errors.errorString", "error.message": "simulated error", "name": "lightning", "server.name": "lightning", "environment": "dev", "version": "1.0.0", "request_id": "123"},
			},
		},
		"error - with stacktrace error": {
			doLogging: func(w io.Writer) {
				m, err := New(Config{ServerName: "lightning", Environment: "dev", Version: "1.0.0", Writer: w})
				require.NoError(t, err)

				m.WithTag("request_id", "123").Error(pkgerrors.WithStack(errors.New("simulated error")), "Unit test exception", StringField("name", "lightning"))
				m.Flush(DefaultFlushWait)
			},
			expected: []map[string]interface{}{
				{"level": "ERROR", "ts": "2025-02-23T13:34:56.185+0700", "msg": "Unit test exception. Err: simulated error", "error.kind": "*errors.withStack", "error.message": "simulated error", "name": "lightning", "error.stack": "github.com/viebiz/lit/monitoring.TestMonitorLogger.func", "server.name": "lightning", "environment": "dev", "version": "1.0.0", "request_id": "123"},
			},
		},
		"error - not in ctx": {
			doLogging: func(w io.Writer) {
				m := FromContext(context.Background())

				m.Error(errors.New("simulated error"), "Skip this log because it is nil")
				m.Flush(DefaultFlushWait)
			},
		},
		"with": {
			doLogging: func(w io.Writer) {
				m, err := New(Config{ServerName: "lightning", Environment: "dev", Version: "1.0.0", Writer: w})
				require.NoError(t, err)

				extraTags := map[string]string{
					"request_id": "123",
					"user_id":    "ID-00001",
				}

				m.With(extraTags).Info("Hello lightning project")
				m.Flush(DefaultFlushWait)
			},
			expected: []map[string]interface{}{
				{"level": "INFO", "ts": "2025-02-23T13:34:56.185+0700", "msg": "Hello lightning project", "request_id": "123", "user_id": "ID-00001", "server.name": "lightning", "environment": "dev", "version": "1.0.0"},
			},
		},
		"getLogField - not in ctx": {
			doLogging: func(w io.Writer) {
				m := FromContext(context.Background())

				m.getLogFields()
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
			parsedLog, err := parseLog(logBuffer.Bytes(), 2)
			require.NoError(t, err)
			testutil.Equal(t, tc.expected, parsedLog, testutil.IgnoreSliceMapEntries(func(k string, v interface{}) bool {
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
		expected  []map[string]interface{}
	}
	tcs := map[string]args{
		"errorf - capture error": {
			givenErr: errors.New("simulated error"),
			expected: []map[string]interface{}{
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
			parsedLog, err := parseLog(logBuffer.Bytes(), 2)
			require.NoError(t, err)
			testutil.Equal(t, tc.expected, parsedLog, testutil.IgnoreSliceMapEntries(func(k string, v interface{}) bool {
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

			require.Greater(t, len(transport.Events()), 0)
			exceptions := transport.Events()[len(transport.Events())-1].Exception
			require.True(t, len(exceptions) > 0)
			lastException := exceptions[len(exceptions)-1]
			require.Equal(t, lastException.Type, "*errors.errorString")
			require.Equal(t, lastException.Value, "simulated error")
			require.True(t, len(lastException.Stacktrace.Frames) > 0)
		})
	}
}

func setupClientTest() (*sentry.Client, *sentry.MockTransport) {
	transport := &sentry.MockTransport{}
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
func TestMonitor_WithTag(t *testing.T) {
	type args struct {
		key   string
		value string
	}
	tcs := map[string]args{
		"add new tag": {
			key:   "request_id",
			value: "123",
		},
		"update existing tag": {
			key:   "server.name",
			value: "updated-server",
		},
		"empty value": {
			key:   "empty_tag",
			value: "",
		},
	}

	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			// Given
			m, err := New(Config{ServerName: "test", Environment: "test", Version: "1.0.0"})
			require.NoError(t, err)

			// When
			child := m.WithTag(tc.key, tc.value)

			// Then
			require.NotNil(t, child)
			require.NotEqual(t, m, child) // Should be different instance
			require.Contains(t, child.logTags, tc.key)
			require.Equal(t, tc.value, child.logTags[tc.key])
			// Parent should remain unchanged
			if tc.key != "server.name" { // server.name is set during New()
				require.NotContains(t, m.logTags, tc.key)
			}
		})
	}
}

func TestMonitor_WithTag_NilMonitor(t *testing.T) {
	// When
	child := (*Monitor)(nil).WithTag("key", "value")

	// Then
	require.Nil(t, child)
}

func TestMonitor_With(t *testing.T) {
	type args struct {
		tags map[string]string
	}
	tcs := map[string]args{
		"add multiple tags": {
			tags: map[string]string{
				"request_id": "123",
				"user_id":    "456",
			},
		},
		"empty map": {
			tags: map[string]string{},
		},
		"nil map": {
			tags: nil,
		},
	}

	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			// Given
			m, err := New(Config{ServerName: "test", Environment: "test", Version: "1.0.0"})
			require.NoError(t, err)

			// When
			child := m.With(tc.tags)

			// Then
			require.NotNil(t, child)
			if len(tc.tags) > 0 {
				require.NotEqual(t, m, child) // Should be different instance only when tags are added
			}
			for k, v := range tc.tags {
				require.Contains(t, child.logTags, k)
				require.Equal(t, v, child.logTags[k])
			}
		})
	}
}

func TestMonitor_With_NilMonitor(t *testing.T) {
	// When
	child := (*Monitor)(nil).With(map[string]string{"key": "value"})

	// Then
	require.Nil(t, child)
}

func TestMonitor_getLogFields(t *testing.T) {
	tcs := map[string]struct {
		tags         map[string]string
		expectedKeys []string
	}{
		"with tags": {
			tags: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
			expectedKeys: []string{"key1", "key2"},
		},
		"empty tags": {
			tags:         map[string]string{},
			expectedKeys: []string{},
		},
	}

	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			// Given
			m := &Monitor{
				logTags: tc.tags,
			}

			// When
			fields := m.getLogFields()

			// Then
			require.Len(t, fields, len(tc.expectedKeys))
			for _, key := range tc.expectedKeys {
				found := false
				for _, field := range fields {
					if field.Key == key {
						require.Equal(t, tc.tags[key], field.String)
						found = true
						break
					}
				}
				require.True(t, found, "Expected key %s not found in fields", key)
			}
		})
	}
}

func TestMonitor_getLogFields_NilMonitor(t *testing.T) {
	// When
	fields := (*Monitor)(nil).getLogFields()

	// Then
	require.Nil(t, fields)
}

func TestMonitor_Flush(t *testing.T) {
	// Given
	m, err := New(Config{ServerName: "test", Environment: "test", Version: "1.0.0"})
	require.NoError(t, err)

	// When
	m.Flush(DefaultFlushWait)

	// Then
	// No panic should occur, and function should complete
}

func TestMonitor_Flush_NilMonitor(t *testing.T) {
	// When
	(*Monitor)(nil).Flush(DefaultFlushWait)

	// Then
	// No panic should occur
}

func parseLog(b []byte, skip int) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	for idx, s := range strings.Split(string(b), "\n") {
		if s == "" {
			break
		}
		if idx < skip {
			continue // Go to next row
		}
		var r map[string]interface{}
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

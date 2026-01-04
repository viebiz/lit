package grpcclient

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/viebiz/lit/grpcclient/testdata"
	"github.com/viebiz/lit/monitoring"
)

func TestSerializeProtoMessage_ReturnsJSON(t *testing.T) {
	t.Parallel()

	// Given
	req := &testdata.WeatherRequest{
		Date: "M41.993.32",
	}

	// When
	got := serializeProtoMessage(req)

	// Then
	require.Equal(t, `{"date":"M41.993.32"}`, got)
}

func TestSerializeProtoMessage_NonProto_ReturnsEmpty(t *testing.T) {
	t.Parallel()

	// Given
	req := struct {
		Foo string
	}{Foo: "bar"}

	// When
	got := serializeProtoMessage(req)

	// Then
	require.Equal(t, "", got)
}

func TestLogRequestBody_WritesExpectedLog(t *testing.T) {
	t.Parallel()

	// Given
	logBuffer := new(bytes.Buffer)
	m, err := monitoring.New(monitoring.Config{
		ServerName:  "lightning",
		Environment: "dev",
		Version:     "1.0.0",
		Writer:      logBuffer,
	})
	require.NoError(t, err)
	ctx := monitoring.SetInContext(context.Background(), m)

	req := &testdata.WeatherRequest{
		Date: "M41.993.32",
	}

	// When
	logRequestBody(ctx, req)

	// Then
	parsed, perr := parseLog(logBuffer.Bytes(), 2) // Skip init logs (sentry + otel)
	require.NoError(t, perr)
	require.Len(t, parsed, 1)
	require.Equal(t, "grpc.outgoing_request", parsed[0]["msg"])
	require.Equal(t, `{"date":"M41.993.32"}`, parsed[0]["grpc.request"])
	require.Equal(t, "INFO", parsed[0]["level"])
}

func TestLogRequestBody_NoMonitorDoesNotPanic(t *testing.T) {
	t.Parallel()

	// Given
	ctx := context.Background()
	req := &testdata.WeatherRequest{Date: "M41.874.21"}

	// When/Then
	require.NotPanics(t, func() {
		logRequestBody(ctx, req)
	})
}

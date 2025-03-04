package testutil

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
)

func NewTraceID(t *testing.T, v string) trace.TraceID {
	t.Helper()

	traceID, err := trace.TraceIDFromHex(v)
	require.NoError(t, err)
	return traceID
}

func NewSpanID(t *testing.T, v string) trace.SpanID {
	t.Helper()

	spanID, err := trace.SpanIDFromHex(v)
	require.NoError(t, err)
	return spanID
}

func NewTraceState(t *testing.T, value string) trace.TraceState {
	t.Helper()

	st, err := trace.ParseTraceState(value)
	require.NoError(t, err)
	return st
}

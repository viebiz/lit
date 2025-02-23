package testutil

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
)

func Compare(t *testing.T, expected, actual trace.SpanContext) {
	t.Helper()
	if !actual.Equal(expected) {
		t.Errorf("\n mismatched. \n expected: %+v \n got: %+v \n diff:\n%s",
			expected, actual,
			cmp.Diff(expected, actual))
		t.FailNow()
	}
}

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

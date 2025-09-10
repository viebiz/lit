package monitoring

import (
	"context"
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"github.com/viebiz/lit/monitoring/tracing/mocktracer"
	"go.opentelemetry.io/otel/attribute"
)

func TestStartSegment_NoCurSpan(t *testing.T) {
	_, end := StartSegment(context.Background(), "test")
	end()
}

func TestStartSegment_WithCurSpan(t *testing.T) {
	ctx, _ := tracer.Start(context.Background(), "parent_span")

	_, end := StartSegment(ctx, "test")
	end()
}

func TestStartSegmentWithTags_NoCurSpan(t *testing.T) {
	_, end := StartSegmentWithTags(context.Background(), "test", map[string]string{"key": "value"})
	end()
}

func TestStartSegmentWithTags_WithCurSpan(t *testing.T) {
	ctx, _ := tracer.Start(context.Background(), "parent_span")

	_, end := StartSegmentWithTags(ctx, "test", map[string]string{"key": "value"})
	end()
}

func TestInjectField(t *testing.T) {
	tp := mocktracer.Start()
	defer tp.Stop()

	type args struct {
		key   string
		value interface{}
	}
	tcs := map[string]args{
		"string value": {
			key:   "user_id",
			value: "123",
		},
		"int value": {
			key:   "count",
			value: 42,
		},
		"bool value": {
			key:   "enabled",
			value: true,
		},
	}

	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			// Given
			ctx, span := tracer.Start(context.Background(), "test-span")
			m, err := New(Config{ServerName: "test", Environment: "test", Version: "1.0.0"})
			require.NoError(t, err)
			ctx = SetInContext(ctx, m)

			// When
			newCtx := InjectField(ctx, tc.key, tc.value)

			// Then
			require.NotEqual(t, ctx, newCtx)
			// Check that span has the attribute
			spanAttrs := span.(interface{ Attributes() []attribute.KeyValue }).Attributes()
			found := false
			for _, attr := range spanAttrs {
				if string(attr.Key) == tc.key {
					require.Equal(t, fmt.Sprintf("%v", tc.value), attr.Value.AsString())
					found = true
					break
				}
			}
			require.True(t, found, "Attribute %s not found in span", tc.key)

			// Check that monitor has the tag
			newMonitor := FromContext(newCtx)
			require.Contains(t, newMonitor.logTags, tc.key)
			require.Equal(t, fmt.Sprintf("%v", tc.value), newMonitor.logTags[tc.key])
		})
	}
}

func TestInjectFields(t *testing.T) {
	tp := mocktracer.Start()
	defer tp.Stop()

	// Given
	ctx, span := tracer.Start(context.Background(), "test-span")
	m, err := New(Config{ServerName: "test", Environment: "test", Version: "1.0.0"})
	require.NoError(t, err)
	ctx = SetInContext(ctx, m)

	tags := map[string]string{
		"request_id": "123",
		"user_id":    "456",
	}

	// When
	newCtx := InjectFields(ctx, tags)

	// Then
	require.NotEqual(t, ctx, newCtx)

	// Check span attributes
	spanAttrs := span.(interface{ Attributes() []attribute.KeyValue }).Attributes()
	for key, expectedValue := range tags {
		found := false
		for _, attr := range spanAttrs {
			if string(attr.Key) == key {
				require.Equal(t, expectedValue, attr.Value.AsString())
				found = true
				break
			}
		}
		require.True(t, found, "Attribute %s not found in span", key)
	}

	// Check monitor tags
	newMonitor := FromContext(newCtx)
	for key, expectedValue := range tags {
		require.Contains(t, newMonitor.logTags, key)
		require.Equal(t, expectedValue, newMonitor.logTags[key])
	}
}

func TestInjectTracingInfo(t *testing.T) {
	tp := mocktracer.Start()
	defer tp.Stop()

	// Given
	_, span := tracer.Start(context.Background(), "test-span")
	m, err := New(Config{ServerName: "test", Environment: "test", Version: "1.0.0"})
	require.NoError(t, err)

	extraTags := map[string]string{
		"custom_tag": "custom_value",
	}

	// When
	newMonitor := InjectTracingInfo(m, span.SpanContext(), extraTags)

	// Then
	require.NotNil(t, newMonitor)
	require.NotEqual(t, m, newMonitor)

	// Check trace ID and span ID are injected
	require.Contains(t, newMonitor.logTags, traceIDKey)
	require.Contains(t, newMonitor.logTags, spanIDKey)
	require.Equal(t, span.SpanContext().TraceID().String(), newMonitor.logTags[traceIDKey])
	require.Equal(t, span.SpanContext().SpanID().String(), newMonitor.logTags[spanIDKey])

	// Check extra tags
	for key, value := range extraTags {
		require.Contains(t, newMonitor.logTags, key)
		require.Equal(t, value, newMonitor.logTags[key])
	}
}

func TestInjectOutgoingTracingInfo(t *testing.T) {
	tp := mocktracer.Start()
	defer tp.Stop()

	// Given
	_, span := tracer.Start(context.Background(), "test-span")
	m, err := New(Config{ServerName: "test", Environment: "test", Version: "1.0.0"})
	require.NoError(t, err)

	// When
	newMonitor := InjectOutgoingTracingInfo(m, span.SpanContext())

	// Then
	require.NotNil(t, newMonitor)
	require.NotEqual(t, m, newMonitor)

	// Check outgoing trace ID and span ID are injected
	require.Contains(t, newMonitor.logTags, outgoingTraceIDKey)
	require.Contains(t, newMonitor.logTags, outgoingSpanIDKey)
	require.Equal(t, span.SpanContext().TraceID().String(), newMonitor.logTags[outgoingTraceIDKey])
	require.Equal(t, span.SpanContext().SpanID().String(), newMonitor.logTags[outgoingSpanIDKey])
}

func TestNotifyErrorToInstrumentation(t *testing.T) {
	tp := mocktracer.Start()
	defer tp.Stop()

	// Given
	ctx, _ := tracer.Start(context.Background(), "test-span")
	testErr := errors.New("test error")

	// When
	NotifyErrorToInstrumentation(ctx, testErr)

	// Then
	// The function should not panic and should handle the error notification
	// We can't easily test the span status with mocktracer, so we just ensure no panic
}

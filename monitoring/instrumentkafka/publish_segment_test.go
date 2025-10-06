package instrumentkafka

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/viebiz/lit/mocks/mocktracer"
)

func TestPublishSegment_SpanID(t *testing.T) {
	tp := mocktracer.NewTestTracerProvider(t)
	defer tp.Shutdown(t)

	// GIVEN
	_, span := tracer.Start(context.Background(), "test-span")
	expectedSpanID := span.SpanContext().SpanID().String()
	ps := buildPublishSegment(span)

	// WHEN
	spanID := ps.SpanID()

	// THEN
	require.Equal(t, expectedSpanID, spanID)
}

func TestPublishSegment_SpanContext(t *testing.T) {
	tp := mocktracer.NewTestTracerProvider(t)
	defer tp.Shutdown(t)

	// GIVEN
	_, span := tracer.Start(context.Background(), "test-span")
	expectedSpanContext := span.SpanContext()
	ps := buildPublishSegment(span)

	// WHEN
	spanContext := ps.SpanContext()

	// THEN
	require.Equal(t, expectedSpanContext, spanContext)
	require.True(t, spanContext.IsValid())
	require.Equal(t, expectedSpanContext.TraceID(), spanContext.TraceID())
	require.Equal(t, expectedSpanContext.SpanID(), spanContext.SpanID())
}

func TestPublishSegment_End_Success(t *testing.T) {
	tp := mocktracer.NewTestTracerProvider(t)
	defer tp.Shutdown(t)

	// GIVEN
	_, span := tracer.Start(context.Background(), "test-span")
	ps := buildPublishSegment(span)
	partition := int32(5)
	offset := int64(100)

	// WHEN
	ps.End(partition, offset, nil)
}

func TestPublishSegment_End_WithError(t *testing.T) {
	tp := mocktracer.NewTestTracerProvider(t)
	defer tp.Shutdown(t)

	// GIVEN
	_, span := tracer.Start(context.Background(), "test-span")
	ps := buildPublishSegment(span)
	partition := int32(5)
	offset := int64(100)
	sendErr := errors.New("failed to send message")

	// WHEN
	ps.End(partition, offset, sendErr)
}

func TestPublishSegment_BuildPublishSegment(t *testing.T) {
	tp := mocktracer.NewTestTracerProvider(t)
	defer tp.Shutdown(t)

	// GIVEN
	_, span := tracer.Start(context.Background(), "test-span")
	expectedSpanID := span.SpanContext().SpanID().String()

	// WHEN
	ps := buildPublishSegment(span)

	// THEN
	require.NotNil(t, ps)
	require.Equal(t, expectedSpanID, ps.SpanID())
	require.Equal(t, span.SpanContext(), ps.SpanContext())
}

package mocktracer

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

// TracerProviderMock represents a mock TracerProvider for unit test
type TracerProviderMock struct {
	spanExporter *tracetest.InMemoryExporter
}

func (tp TracerProviderMock) GetLatestSpan() SpanStub {
	spans := tp.spanExporter.GetSpans()
	if len(spans) == 0 {
		return SpanStub{}
	}

	return SpanStub(spans[len(spans)-1])
}

func (tp TracerProviderMock) Reset() {
	tp.spanExporter.Reset()
}

func (tp TracerProviderMock) Stop() {
	otel.SetTracerProvider(noop.NewTracerProvider()) // Reset global tracer provider
}

type SpanStub tracetest.SpanStub

// staticIDGenerator supports generate static trace_id & span_id for unit test
type staticIDGenerator struct{}

var _ sdktrace.IDGenerator = (*staticIDGenerator)(nil)

func (gen staticIDGenerator) NewIDs(ctx context.Context) (trace.TraceID, trace.SpanID) {
	traceIDHex := fmt.Sprintf("%032x", 1)
	traceID, _ := trace.TraceIDFromHex(traceIDHex)

	spanID := gen.NewSpanID(ctx, traceID)

	return traceID, spanID
}

func (gen staticIDGenerator) NewSpanID(ctx context.Context, traceID trace.TraceID) trace.SpanID {
	spanIDHex := fmt.Sprintf("%016x", 1)
	spanID, _ := trace.SpanIDFromHex(spanIDHex)
	return spanID
}

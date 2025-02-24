package mocktracer

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

// Start setup TracerProviderMock for testing
func Start() TracerProviderMock {
	// Create in-memory exporter support collect span for unit test
	exp := tracetest.NewInMemoryExporter()

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(resource.Default()),
		sdktrace.WithIDGenerator(&staticIDGenerator{}),
		sdktrace.WithSpanProcessor(
			// To process span immediately
			sdktrace.NewSimpleSpanProcessor(exp),
		),
	)

	otel.SetTracerProvider(tp) // Override tracer provider
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return TracerProviderMock{
		spanExporter: exp,
	}
}

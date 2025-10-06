package mocktracer

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

// TestTracerProvider provides a test-friendly tracer provider with in-memory span collection.
// Use this instead of the mocktracer package for better testing.
type TestTracerProvider struct {
	provider *sdktrace.TracerProvider
	exporter *tracetest.InMemoryExporter
}

// NewTestTracerProvider creates a new test tracer provider with in-memory span collection.
// It automatically sets the global tracer provider and propagator.
//
// Example usage:
//
//	func TestMyFunction(t *testing.T) {
//	    tp := mocktracer.NewTestTracerProvider(t)
//	    defer tp.Shutdown(t)
//
//	    // Your test code here
//	    ctx, span := otel.Tracer("test").Start(context.Background(), "test-span")
//	    defer span.End()
//
//	    // Assert on collected spans
//	    spans := tp.GetSpans()
//	    require.Len(t, spans, 1)
//	    require.Equal(t, "test-span", spans[0].Name())
//	}
func NewTestTracerProvider(t *testing.T, opts ...TestTracerOption) *TestTracerProvider {
	t.Helper()

	// Create in-memory exporter
	exporter := tracetest.NewInMemoryExporter()

	// Default configuration
	cfg := &testTracerConfig{
		serviceName: "test-service",
		environment: "test",
		version:     "1.0.0-test",
		sampler:     sdktrace.AlwaysSample(),
		idGenerator: NewStaticIDGenerator([16]byte{}, [8]byte{}),
	}

	// Apply options
	for _, opt := range opts {
		opt(cfg)
	}

	// Build resource
	rsc := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(cfg.serviceName),
		semconv.ServiceVersion(cfg.version),
		semconv.DeploymentEnvironment(cfg.environment),
	)

	// Create provider with simple span processor for immediate processing
	provider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(cfg.sampler),
		sdktrace.WithResource(rsc),
		sdktrace.WithSpanProcessor(sdktrace.NewSimpleSpanProcessor(exporter)),
		sdktrace.WithIDGenerator(cfg.idGenerator),
	)

	// Set global tracer provider
	otel.SetTracerProvider(provider)

	// Setup propagators
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return &TestTracerProvider{
		provider: provider,
		exporter: exporter,
	}
}

// GetSpans returns all spans collected by the in-memory exporter.
func (tp *TestTracerProvider) GetSpans() []tracetest.SpanStub {
	return tp.exporter.GetSpans()
}

// GetLatestSpan returns the most recently collected span.
// Returns an empty SpanStub if no spans have been collected.
func (tp *TestTracerProvider) GetLatestSpan() tracetest.SpanStub {
	spans := tp.exporter.GetSpans()
	if len(spans) == 0 {
		return tracetest.SpanStub{}
	}
	return spans[len(spans)-1]
}

// Reset clears all collected spans.
// Useful for testing multiple scenarios in a single test.
func (tp *TestTracerProvider) Reset() {
	tp.exporter.Reset()
}

// Shutdown shuts down the tracer provider and resets the global tracer provider.
// Should be called with defer in tests.
func (tp *TestTracerProvider) Shutdown(t *testing.T) {
	t.Helper()

	ctx := context.Background()
	if err := tp.provider.Shutdown(ctx); err != nil {
		t.Errorf("failed to shutdown tracer provider: %v", err)
	}

	// Reset global tracer provider to noop
	otel.SetTracerProvider(noop.NewTracerProvider())
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator())
}

// testTracerConfig holds configuration for test tracer provider.
type testTracerConfig struct {
	serviceName string
	environment string
	version     string
	sampler     sdktrace.Sampler
	idGenerator sdktrace.IDGenerator
}

// TestTracerOption is a function that configures a test tracer provider.
type TestTracerOption func(*testTracerConfig)

// WithTestServiceName sets the service name for the test tracer.
func WithTestServiceName(name string) TestTracerOption {
	return func(cfg *testTracerConfig) {
		cfg.serviceName = name
	}
}

// WithTestEnvironment sets the environment for the test tracer.
func WithTestEnvironment(env string) TestTracerOption {
	return func(cfg *testTracerConfig) {
		cfg.environment = env
	}
}

// WithTestVersion sets the version for the test tracer.
func WithTestVersion(version string) TestTracerOption {
	return func(cfg *testTracerConfig) {
		cfg.version = version
	}
}

// WithTestSampler sets the sampler for the test tracer.
// Default is AlwaysSample().
func WithTestSampler(sampler sdktrace.Sampler) TestTracerOption {
	return func(cfg *testTracerConfig) {
		cfg.sampler = sampler
	}
}

// WithTestIDGenerator sets a custom ID generator for the test tracer.
// Useful for generating deterministic trace and span IDs in tests.
func WithTestIDGenerator(gen sdktrace.IDGenerator) TestTracerOption {
	return func(cfg *testTracerConfig) {
		cfg.idGenerator = gen
	}
}

// StaticIDGenerator generates static, deterministic trace and span IDs for testing.
// This is useful when you need predictable IDs in your tests.
type StaticIDGenerator struct {
	traceID trace.TraceID
	spanID  trace.SpanID
}

// NewStaticIDGenerator creates a new static ID generator with the provided IDs.
// If traceID is zero, it defaults to "00000000000000000000000000000001".
// If spanID is zero, it defaults to "0000000000000001".
func NewStaticIDGenerator(traceID trace.TraceID, spanID trace.SpanID) *StaticIDGenerator {
	if traceID.IsValid() == false {
		traceID, _ = trace.TraceIDFromHex("00000000000000000000000000000001")
	}
	if spanID.IsValid() == false {
		spanID, _ = trace.SpanIDFromHex("0000000000000001")
	}
	return &StaticIDGenerator{
		traceID: traceID,
		spanID:  spanID,
	}
}

// NewIDs returns the static trace and span IDs.
func (g *StaticIDGenerator) NewIDs(ctx context.Context) (trace.TraceID, trace.SpanID) {
	return g.traceID, g.spanID
}

// NewSpanID returns the static span ID.
func (g *StaticIDGenerator) NewSpanID(ctx context.Context, traceID trace.TraceID) trace.SpanID {
	return g.spanID
}

// IncrementingIDGenerator generates incrementing trace and span IDs for testing.
// Each call to NewIDs or NewSpanID increments the counter.
type IncrementingIDGenerator struct {
	traceCounter uint64
	spanCounter  uint64
}

// NewIncrementingIDGenerator creates a new incrementing ID generator.
func NewIncrementingIDGenerator() *IncrementingIDGenerator {
	return &IncrementingIDGenerator{
		traceCounter: 1,
		spanCounter:  1,
	}
}

// NewIDs returns a new trace ID and span ID, incrementing both counters.
func (g *IncrementingIDGenerator) NewIDs(ctx context.Context) (trace.TraceID, trace.SpanID) {
	traceID := g.nextTraceID()
	spanID := g.nextSpanID()
	return traceID, spanID
}

// NewSpanID returns a new span ID, incrementing the span counter.
func (g *IncrementingIDGenerator) NewSpanID(ctx context.Context, traceID trace.TraceID) trace.SpanID {
	return g.nextSpanID()
}

func (g *IncrementingIDGenerator) nextTraceID() trace.TraceID {
	var id trace.TraceID
	// Use the counter as the last 8 bytes of the trace ID
	for i := 0; i < 8; i++ {
		id[8+i] = byte(g.traceCounter >> (56 - uint(i)*8))
	}
	g.traceCounter++
	return id
}

func (g *IncrementingIDGenerator) nextSpanID() trace.SpanID {
	var id trace.SpanID
	for i := 0; i < 8; i++ {
		id[i] = byte(g.spanCounter >> (56 - uint(i)*8))
	}
	g.spanCounter++
	return id
}

// WithStaticIDs is a convenience option that creates a test tracer with static IDs.
// Useful when you need predictable trace and span IDs in tests.
//
// Example:
//
//	tp := mocktracer.NewTestTracerProvider(t,
//	    tracing.WithStaticIDs("00000000000000000000000000000001", "0000000000000001"),
//	)
func WithStaticIDs(traceIDHex, spanIDHex string) TestTracerOption {
	return func(cfg *testTracerConfig) {
		traceID, _ := trace.TraceIDFromHex(traceIDHex)
		spanID, _ := trace.SpanIDFromHex(spanIDHex)
		cfg.idGenerator = NewStaticIDGenerator(traceID, spanID)
	}
}

// WithIncrementingIDs is a convenience option that creates a test tracer with incrementing IDs.
// Useful when you need unique but predictable trace and span IDs across multiple spans.
func WithIncrementingIDs() TestTracerOption {
	return func(cfg *testTracerConfig) {
		cfg.idGenerator = NewIncrementingIDGenerator()
	}
}

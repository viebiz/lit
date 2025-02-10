package monitoring

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

type MonitorTest struct {
	logBuffer    *bytes.Buffer
	logger       *Logger
	spanExporter *tracetest.InMemoryExporter
	tp           trace.TracerProvider
	propagator   propagation.TextMapPropagator
}

func NewMonitorTest() (MonitorTest, func()) {
	logBuffer := new(bytes.Buffer)
	logger := NewLoggerWithWriter(logBuffer)

	exporter := tracetest.NewInMemoryExporter()

	m := MonitorTest{
		logBuffer:    logBuffer,
		logger:       logger,
		spanExporter: exporter,
		tp: sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithSyncer(exporter),
			sdktrace.WithIDGenerator(&staticIDGenerator{}),
		),
		propagator: propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{}, // Support W3C TraceContext.
			propagation.Baggage{},      // Support baggage propagation.
		),
	}

	originGetTracer := getTracer
	originGetPropagator := getTextMapPropagator

	// Override getTracer, getPropagator for testing
	getTracer = m.GetTracer
	getTextMapPropagator = m.GetPropagator

	return m, func() {
		getTracer = originGetTracer
		getTextMapPropagator = originGetPropagator
	}
}

func (m MonitorTest) Context() context.Context {
	return SetInContext(context.Background(), m.logger)
}

func (m MonitorTest) GetLogger() *Logger {
	return m.logger
}

func (m MonitorTest) GetLogs(t *testing.T) []map[string]interface{} {
	t.Helper()
	var logs []map[string]interface{}

	lines := bytes.Split(m.logBuffer.Bytes(), []byte("\n")) // \n is end of line
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		var msg map[string]interface{}
		require.NoError(t, json.Unmarshal(line, &msg))
		logs = append(logs, msg)
	}

	return logs
}

func (m MonitorTest) GetTracer() trace.Tracer {
	return m.tp.Tracer(tracerName)
}

func (m MonitorTest) GetPropagator() propagation.TextMapPropagator {
	return m.propagator
}

func (m MonitorTest) GetSpans() tracetest.SpanStubs {
	return m.spanExporter.GetSpans()
}

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

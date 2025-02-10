package monitoring

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const (
	tracerName = "gitlab.com/bizgroup2/lightning/monitoring"
)

var (
	// getTracer creates tracer by global TracerProvider
	getTracer = func() trace.Tracer {
		return otel.Tracer(tracerName)
	}

	// getTextMapPropagator returns global TextMapPropagator configs
	getTextMapPropagator = func() propagation.TextMapPropagator {
		return otel.GetTextMapPropagator()
	}
)

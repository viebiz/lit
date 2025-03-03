package monitoring

import (
	"go.opentelemetry.io/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	tracerName = "github.com/viebiz/lit/monitoring"
)

var (
	tracer = otel.Tracer(tracerName, trace.WithSchemaURL(semconv.SchemaURL))
)

package monitoring

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
	"go.opentelemetry.io/otel/trace"
)

// StartVaultSegment starts a trace.Span with vault information
func StartVaultSegment(ctx context.Context, info ExternalServiceInfo, operation string) func(error) {
	opts := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			semconv.ServerAddress(info.Hostname),
		),
	}

	_, span := getTracer().Start(ctx, fmt.Sprintf("vault.%s", operation), opts...)

	return func(err error) {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err, trace.WithStackTrace(true))
		}

		span.End()
	}
}

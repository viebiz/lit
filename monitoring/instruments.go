package monitoring

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const (
	traceIDKey         = "trace_id"
	spanIDKey          = "span_id"
	outgoingTraceIDKey = "outgoing_trace_id"
	outgoingSpanIDKey  = "outgoing_span_id"
)

// InjectField injects a field to Logger and trace.Span in context
func InjectField[T any](ctx context.Context, key string, value T) context.Context {
	trace.SpanFromContext(ctx).SetAttributes(attribute.String(key, fmt.Sprintf("%v", value)))
	return SetInContext(ctx, FromContext(ctx).With(Field(key, value)))
}

func InjectFields(ctx context.Context, fields map[string]string) context.Context {
	attrs := make([]attribute.KeyValue, 0, len(fields))
	logFields := make([]LogField, 0, len(fields))
	for k, v := range fields {
		attrs = append(attrs, attribute.String(k, v))
		logFields = append(logFields, Field(k, v))
	}

	trace.SpanFromContext(ctx).SetAttributes(attrs...)
	return SetInContext(ctx, FromContext(ctx).With(logFields...))
}

func injectTracingInfo(ctx context.Context, spanCtx trace.SpanContext) context.Context {
	return SetInContext(ctx, FromContext(ctx).
		With(
			Field(traceIDKey, spanCtx.TraceID().String()),
			Field(spanIDKey, spanCtx.SpanID().String()),
		),
	)
}

func injectOutgoingTracingInfo(ctx context.Context, spanCtx trace.SpanContext) context.Context {
	return SetInContext(ctx, FromContext(ctx).
		With(
			Field(outgoingTraceIDKey, spanCtx.TraceID().String()),
			Field(outgoingSpanIDKey, spanCtx.SpanID().String()),
		),
	)
}

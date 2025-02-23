package monitoring

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
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
	valStr := fmt.Sprintf("%v", value) // TODO: Optimize it
	trace.SpanFromContext(ctx).SetAttributes(attribute.String(key, valStr))
	return SetInContext(ctx, FromContext(ctx).WithTag(key, valStr))
}

func InjectFields(ctx context.Context, tags map[string]string) context.Context {
	attrs := make([]attribute.KeyValue, 0, len(tags))
	for k, v := range tags {
		attrs = append(attrs, attribute.String(k, v))
	}

	trace.SpanFromContext(ctx).SetAttributes(attrs...)
	return SetInContext(ctx, FromContext(ctx).With(tags))
}

func InjectTracingInfo(m *Monitor, spanCtx trace.SpanContext) *Monitor {
	return m.With(map[string]string{
		traceIDKey: spanCtx.TraceID().String(),
		spanIDKey:  spanCtx.SpanID().String(),
	})
}

func InjectOutgoingTracingInfo(m *Monitor, spanCtx trace.SpanContext) *Monitor {
	return m.With(map[string]string{
		outgoingTraceIDKey: spanCtx.TraceID().String(),
		outgoingSpanIDKey:  spanCtx.SpanID().String(),
	})
}

func NotifyErrorToInstrumentation(ctx context.Context, err error) {
	span := trace.SpanFromContext(ctx)
	span.SetStatus(codes.Error, err.Error())
}

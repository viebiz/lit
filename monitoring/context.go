package monitoring

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

// monitorContextKey is the context key used to retrieve a Monitor from context.
// An empty struct is preferred for efficiency: https://github.com/golang/go/issues/17826#issuecomment-259035465
type monitorContextKey struct{}

// SetInContext sets the logger in context
func SetInContext(ctx context.Context, m *Monitor) context.Context {
	return context.WithValue(ctx, monitorContextKey{}, m)
}

// FromContext gets the logger from context
func FromContext(ctx context.Context) *Monitor {
	if m, ok := ctx.Value(monitorContextKey{}).(*Monitor); ok && m != nil {
		return m
	}

	return nil
}

// NewContext copies the logger from old to a new context
// Use this when you want to use a new context but copy the logger over from the original context
func NewContext(ctx context.Context) context.Context {
	newCtx := SetInContext(context.Background(), FromContext(ctx))

	if span := trace.SpanFromContext(ctx); span.IsRecording() {
		newCtx = trace.ContextWithSpan(newCtx, span)
	}

	return newCtx
}

package monitoring

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

type contextKey uint8

const (
	loggerContextKey contextKey = 0
)

// SetInContext sets the logger in context
func SetInContext(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, loggerContextKey, logger)
}

// FromContext gets the logger from context
func FromContext(ctx context.Context) *Logger {
	logger, ok := ctx.Value(loggerContextKey).(*Logger)
	if !ok {
		return NewNoopLogger()
	}

	return logger
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

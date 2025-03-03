package monitoring

import (
	"context"
	"testing"
)

func TestStartSegment_NoCurSpan(t *testing.T) {
	_, end := StartSegment(context.Background(), "test")
	end()
}

func TestStartSegment_WithCurSpan(t *testing.T) {
	ctx, _ := tracer.Start(context.Background(), "parent_span")

	_, end := StartSegment(ctx, "test")
	end()
}

func TestStartSegmentWithTags_NoCurSpan(t *testing.T) {
	_, end := StartSegmentWithTags(context.Background(), "test", map[string]string{"key": "value"})
	end()
}

func TestStartSegmentWithTags_WithCurSpan(t *testing.T) {
	ctx, _ := tracer.Start(context.Background(), "parent_span")

	_, end := StartSegmentWithTags(ctx, "test", map[string]string{"key": "value"})
	end()
}

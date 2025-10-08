package instrumentkafka

import (
	"context"
	"errors"
	"testing"

	"github.com/IBM/sarama"
	"github.com/viebiz/lit/mocks/mocktracer"
	"go.opentelemetry.io/otel"
)

func TestStartConsumeTxn_NoCurSpan(t *testing.T) {
	tp := mocktracer.NewTestTracerProvider(t)
	defer tp.Shutdown(t)

	// GIVEN
	ctx := context.Background()
	msg := &sarama.ConsumerMessage{
		Topic:     "test-topic",
		Partition: 1,
		Offset:    100,
		Key:       []byte("test-key"),
		Value:     []byte("test-value"),
	}

	// WHEN
	ctx, end := StartConsumeTxn(ctx, msg)

	// THEN
	requireTraceContextPresent(t, ctx)
	consumerErr := errors.New("consumer error")
	end(consumerErr)
}

func TestStartConsumeTxn_WithSpanContext(t *testing.T) {
	tp := mocktracer.NewTestTracerProvider(t)
	defer tp.Shutdown(t)

	// GIVEN
	ctx := context.Background()
	msg := &sarama.ConsumerMessage{
		Topic:     "test-topic",
		Partition: 1,
		Offset:    100,
		Key:       []byte("test-key"),
		Value:     []byte("test-value"),
		Headers:   []*sarama.RecordHeader{},
	}

	// Inject parent span context into message headers
	parentCtx, _ := tracer.Start(ctx, "parent-span")
	otel.GetTextMapPropagator().Inject(parentCtx, buildConsumeMessageCarrier(msg))

	// WHEN
	ctx, end := StartConsumeTxn(ctx, msg)

	// THEN
	requireTraceContextPresent(t, ctx)
	end(nil)
}

func TestStartCommitSegment(t *testing.T) {
	tp := mocktracer.NewTestTracerProvider(t)
	defer tp.Shutdown(t)

	// GIVEN
	ctx := context.Background()
	topic := "test-topic"
	key := "test-key"
	partition := int32(1)
	offset := int64(100)

	// WHEN
	end := StartCommitSegment(ctx, topic, key, partition, offset)

	// THEN
	end()
}

func TestStartCommitSegment_WithExistingSpan(t *testing.T) {
	tp := mocktracer.NewTestTracerProvider(t)
	defer tp.Shutdown(t)

	// GIVEN
	ctx, parentSpan := tracer.Start(context.Background(), "parent-span")
	topic := "test-topic"
	key := "test-key"
	partition := int32(1)
	offset := int64(100)

	// WHEN
	end := StartCommitSegment(ctx, topic, key, partition, offset)

	// THEN
	end()
	parentSpan.End()
}

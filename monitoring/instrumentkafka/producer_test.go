package instrumentkafka

import (
	"context"
	"errors"
	"testing"

	"github.com/IBM/sarama"
	"github.com/stretchr/testify/require"
	"github.com/viebiz/lit/mocks/mocktracer"
	"github.com/viebiz/lit/monitoring"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func TestStartSyncPublishSegment_NoCurSpan(t *testing.T) {
	tp := mocktracer.NewTestTracerProvider(t)
	defer tp.Shutdown(t)

	// GIVEN
	ctx := context.Background()
	extSvcInfo := monitoring.ExternalServiceInfo{
		Hostname: "kafka-broker",
		Port:     "9092",
	}
	clientID := "test-client"
	pm := &sarama.ProducerMessage{
		Topic:     "test-topic",
		Partition: 1,
	}
	key := "test-key"

	// WHEN
	ctx, end := StartSyncPublishSegment(ctx, extSvcInfo, clientID, pm, key)

	// THEN
	requireTraceContextPresent(t, ctx)

	partition := int32(2)
	offset := int64(100)
	sendErr := errors.New("send error")
	end(partition, offset, sendErr)
}

func TestStartSyncPublishSegment_WithCurSpan(t *testing.T) {
	tp := mocktracer.NewTestTracerProvider(t)
	defer tp.Shutdown(t)

	// GIVEN
	ctx, _ := tracer.Start(context.Background(), "test-parent-span")
	extSvcInfo := monitoring.ExternalServiceInfo{
		Hostname: "kafka-broker",
		Port:     "9092",
	}
	clientID := "test-client"
	pm := &sarama.ProducerMessage{
		Topic:     "test-topic",
		Partition: 1,
	}
	key := "test-key"

	// WHEN
	ctx, end := StartSyncPublishSegment(ctx, extSvcInfo, clientID, pm, key)

	// THEN
	requireTraceContextPresentFromHeader(t, pm)

	partition := int32(2)
	offset := int64(100)
	end(partition, offset, errors.New("send error"))
}

func TestStartAsyncEnqueueSegment_NoCurSpan(t *testing.T) {
	tp := mocktracer.NewTestTracerProvider(t)
	defer tp.Shutdown(t)

	// GIVEN
	ctx := context.Background()
	clientID := "test-client"
	pm := &sarama.ProducerMessage{
		Topic:     "test-topic",
		Partition: 1,
	}
	key := "test-key"

	// WHEN
	ctx, _, end := StartAsyncEnqueueSegment(ctx, clientID, pm, key)

	// THEN
	requireTraceContextPresent(t, ctx)
	end()
}

func TestStartAsyncEnqueueSegment_WithCurSpan(t *testing.T) {
	tp := mocktracer.NewTestTracerProvider(t)
	defer tp.Shutdown(t)

	// Create a parent span
	ctx, _ := tracer.Start(context.Background(), "test-parent-span")

	// GIVEN
	clientID := "test-client"
	pm := &sarama.ProducerMessage{
		Topic:     "test-topic",
		Partition: 1,
	}
	key := "test-key"

	// WHEN
	ctx, _, end := StartAsyncEnqueueSegment(ctx, clientID, pm, key)

	// THEN
	requireTraceContextPresent(t, ctx)
	requireTraceContextPresentFromHeader(t, pm)
	end()
}

func TestStartAsyncPublishSegment(t *testing.T) {
	tp := mocktracer.NewTestTracerProvider(t)
	defer tp.Shutdown(t)

	// GIVEN
	extSvcInfo := monitoring.ExternalServiceInfo{
		Hostname: "kafka-broker",
		Port:     "9092",
	}
	clientID := "test-client"
	pm := &sarama.ProducerMessage{
		Topic:     "test-topic",
		Partition: 1,
	}
	key := "test-key"

	ctx, _ := tracer.Start(context.Background(), "parent-span")
	otel.GetTextMapPropagator().Inject(ctx, buildProduceMessageCarrier(pm))

	// WHEN
	ps := StartAsyncPublishSegment(extSvcInfo, clientID, pm, key)

	// THEN
	requireTraceContextPresentFromHeader(t, pm)
	partition := int32(2)
	offset := int64(100)
	ps.End(partition, offset, nil)
}

func TestStartAsyncPublishSegment_NoExistingContext(t *testing.T) {
	tp := mocktracer.NewTestTracerProvider(t)
	defer tp.Shutdown(t)

	// GIVEN
	extSvcInfo := monitoring.ExternalServiceInfo{
		Hostname: "kafka-broker",
		Port:     "9092",
	}
	clientID := "test-client"
	pm := &sarama.ProducerMessage{
		Topic:     "test-topic",
		Partition: 1,
	}
	key := "test-key"

	// WHEN
	ps := StartAsyncPublishSegment(extSvcInfo, clientID, pm, key)

	// THEN
	requireTraceContextPresentFromHeader(t, pm)
	partition := int32(2)
	offset := int64(100)
	ps.End(partition, offset, nil)
}

func requireTraceContextPresentFromHeader(t *testing.T, pm *sarama.ProducerMessage) {
	t.Helper()
	ctx := context.Background()
	ctx = otel.GetTextMapPropagator().Extract(ctx, buildProduceMessageCarrier(pm))
	requireTraceContextPresent(t, ctx)
}

func requireTraceContextPresent(t *testing.T, ctx context.Context) {
	t.Helper()
	span := trace.SpanFromContext(ctx)
	require.NotNil(t, span, "span should be present")
	require.True(t, span.SpanContext().IsValid(), "span context should be valid")
}

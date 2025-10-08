package instrumentgrpc

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/viebiz/lit/mocks/mocktracer"
	"github.com/viebiz/lit/monitoring"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
)

func TestStartUnaryCallSegment_NoCurSpan(t *testing.T) {
	tp := mocktracer.NewTestTracerProvider(t)
	defer tp.Shutdown(t)

	// Given
	ctx := context.Background()
	svcInfo := monitoring.ExternalServiceInfo{
		Hostname: "example.com",
		Port:     "50051",
	}

	// When
	newCtx, endFunc := StartUnaryCallSegment(ctx, svcInfo, "/weather.WeatherService/GetWeatherInfo")

	// Then
	requireTraceContextPresent(t, newCtx)

	// Verify metadata was properly injected
	md, ok := metadata.FromOutgoingContext(newCtx)
	require.True(t, ok)
	require.NotEmpty(t, md.Get("traceparent"))
	endFunc(nil)
}

func TestStartUnaryCallSegment_WithCurSpan(t *testing.T) {
	tp := mocktracer.NewTestTracerProvider(t)
	defer tp.Shutdown(t)

	// Given
	ctx := context.Background()
	ctx, parentSpan := tracer.Start(ctx, "service-span")

	svcInfo := monitoring.ExternalServiceInfo{
		Hostname: "example.com",
		Port:     "50051",
	}

	// When
	newCtx, endFunc := StartUnaryCallSegment(ctx, svcInfo, "/weather.WeatherService/GetWeatherInfo")

	// Then
	requireTraceContextPresent(t, newCtx)
	requireTraceIDMatch(t, parentSpan.SpanContext().TraceID().String(), trace.SpanFromContext(newCtx))

	md, ok := metadata.FromOutgoingContext(newCtx)
	require.True(t, ok)
	require.NotEmpty(t, md.Get("traceparent"))

	endFunc(nil)
}

func TestStartUnaryCallSegment_EndWithError(t *testing.T) {
	tp := mocktracer.NewTestTracerProvider(t)
	defer tp.Shutdown(t)

	// Given
	ctx := context.Background()
	svcInfo := monitoring.ExternalServiceInfo{
		Hostname: "example.com",
		Port:     "50051",
	}

	// When
	newCtx, endFunc := StartUnaryCallSegment(ctx, svcInfo, "/weather.WeatherService/GetWeatherInfo")

	// Then
	requireTraceContextPresent(t, newCtx)

	endFunc(errors.New("simulated error"))
}

func TestStartUnaryCallSegment_WithExistingMetadata(t *testing.T) {
	tp := mocktracer.NewTestTracerProvider(t)
	defer tp.Shutdown(t)

	// Given
	ctx := context.Background()
	existingMD := metadata.New(map[string]string{
		"existing-key": "existing-value",
	})
	ctx = metadata.NewOutgoingContext(ctx, existingMD)

	svcInfo := monitoring.ExternalServiceInfo{
		Hostname: "example.com",
		Port:     "50051",
	}

	// When
	newCtx, endFunc := StartUnaryCallSegment(ctx, svcInfo, "/weather.WeatherService/GetWeatherInfo")

	// Then
	requireTraceContextPresent(t, newCtx)

	// Verify existing metadata was preserved
	md, ok := metadata.FromOutgoingContext(newCtx)
	require.True(t, ok)
	require.Equal(t, "existing-value", md.Get("existing-key")[0])
	require.NotEmpty(t, md.Get("traceparent"))

	endFunc(nil)
}

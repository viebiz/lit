package instrumentgrpc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/viebiz/lit/grpcclient/testdata"
	"github.com/viebiz/lit/mocks/mocktracer"
	"github.com/viebiz/lit/monitoring"
	"github.com/viebiz/lit/testutil"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

func TestStartUnaryIncomingCall_NoCurSpan(t *testing.T) {
	tp := mocktracer.NewTestTracerProvider(t)
	defer tp.Shutdown(t)

	logBuf := bytes.NewBuffer(nil)
	m, err := monitoring.New(monitoring.Config{Writer: logBuf})
	require.NoError(t, err)

	// Given
	reqCtx := context.Background()
	reqCtx = peer.NewContext(reqCtx, &peer.Peer{
		Addr: &net.TCPAddr{
			Port: 50051,
		},
	})

	// When
	ctx, reqMeta, endFunc := StartUnaryIncomingCall(reqCtx, m, "/weather.WeatherService/GetWeatherInfo", &testdata.WeatherRequest{
		Date: "M41.993.32",
	})

	// Then
	requireTraceContextPresent(t, ctx)
	require.Equal(t, "/weather.WeatherService/GetWeatherInfo", reqMeta.ServiceMethod)
	var reqBody map[string]interface{}
	require.NoError(t, json.Unmarshal(reqMeta.BodyToLog, &reqBody))
	testutil.Equal(t, map[string]interface{}{"date": "M41.993.32"}, reqBody)
	endFunc(nil)
}

func TestStartUnaryIncomingCall_WihCurSpan(t *testing.T) {
	tp := mocktracer.NewTestTracerProvider(t)
	defer tp.Shutdown(t)

	logBuf := bytes.NewBuffer(nil)
	m, err := monitoring.New(monitoring.Config{Writer: logBuf})
	require.NoError(t, err)

	// Given
	expTraceID := "deadbeefcafebabefeedfacebadc0de1"
	reqCtx := context.Background()
	reqCtx = peer.NewContext(reqCtx, &peer.Peer{
		Addr: &net.TCPAddr{
			Port: 50051,
		},
	})
	reqCtx = metadata.NewIncomingContext(reqCtx, metadata.New(map[string]string{
		"traceparent": fmt.Sprintf("00-%s-abad1dea0ddba11c-01", expTraceID),
	}))

	// When
	ctx, reqMeta, endFunc := StartUnaryIncomingCall(reqCtx, m, "/weather.WeatherService/GetWeatherInfo", &testdata.WeatherRequest{
		Date: "M41.993.32",
	})

	// Then
	requireTraceContextPresent(t, ctx)
	requireTraceIDMatch(t, expTraceID, trace.SpanFromContext(ctx))
	require.Equal(t, "/weather.WeatherService/GetWeatherInfo", reqMeta.ServiceMethod)
	var reqBody map[string]interface{}
	require.NoError(t, json.Unmarshal(reqMeta.BodyToLog, &reqBody))
	testutil.Equal(t, map[string]interface{}{"date": "M41.993.32"}, reqBody)
	endFunc(nil)
}

func TestStartUnaryIncomingCall_EndWithError(t *testing.T) {
	tp := mocktracer.NewTestTracerProvider(t)
	defer tp.Shutdown(t)

	logBuf := bytes.NewBuffer(nil)
	m, err := monitoring.New(monitoring.Config{Writer: logBuf})
	require.NoError(t, err)

	// Given
	reqCtx := context.Background()
	reqCtx = peer.NewContext(reqCtx, &peer.Peer{
		Addr: &net.TCPAddr{
			Port: 50051,
		},
	})

	// When
	ctx, reqMeta, endFunc := StartUnaryIncomingCall(reqCtx, m, "/weather.WeatherService/GetWeatherInfo", &testdata.WeatherRequest{
		Date: "M41.993.32",
	})

	// Then
	requireTraceContextPresent(t, ctx)
	require.Equal(t, "/weather.WeatherService/GetWeatherInfo", reqMeta.ServiceMethod)
	var reqBody map[string]interface{}
	require.NoError(t, json.Unmarshal(reqMeta.BodyToLog, &reqBody))
	testutil.Equal(t, map[string]interface{}{"date": "M41.993.32"}, reqBody)
	endFunc(errors.New("simulated error"))
}

func requireTraceContextPresent(t *testing.T, ctx context.Context) {
	t.Helper()
	span := trace.SpanFromContext(ctx)
	require.NotNil(t, span, "span should be present")
	require.True(t, span.SpanContext().IsValid(), "span context should be valid")
}

func requireTraceIDMatch(t *testing.T, expectedID string, span trace.Span) {
	t.Helper()
	require.Equal(t, expectedID, span.SpanContext().TraceID().String())
}

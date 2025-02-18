package monitoring

import (
	"net"
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"

	"github.com/viebiz/lit/grpcclient/testdata"
)

func TestStartUnaryIncomingCall(t *testing.T) {
	// Given
	monitor, endTest := NewMonitorTest()
	defer endTest()

	reqCtx := peer.NewContext(monitor.Context(), &peer.Peer{
		Addr: &net.TCPAddr{
			Port: 50051,
		},
	})

	reqCtx = metadata.NewIncomingContext(reqCtx, metadata.New(map[string]string{
		"traceparent": "00-deadbeefcafebabefeedfacebadc0de1-abad1dea0ddba11c-01",
		"tracestate":  "test=test-value",
		"baggage":     "user_id=1234,role=admin",
	}))

	// When
	ctx, reqMeta, end := StartUnaryIncomingCall(reqCtx, "/weather.WeatherService/GetWeatherInfo", &testdata.WeatherRequest{
		Date: "M41.993.32",
	})

	// Then
	require.NotNil(t, ctx.Value(loggerContextKey{}))
	FromContext(ctx).Infof("Got incoming request")
	expectedLogs := []map[string]interface{}{
		{
			"level":    "info",
			"msg":      "Got incoming request",
			"span_id":  "0000000000000000",                 // Random generated value
			"trace_id": "deadbeefcafebabefeedfacebadc0de1", // Should sample with incoming request
		},
	}
	requireEqual(t, expectedLogs, monitor.GetLogs(t), cmpopts.IgnoreMapEntries(func(key string, value interface{}) bool {
		return key == "ts" || key == "span_id"
	}))

	// 3.2. Validate request metadata
	requireEqual(t, GRPCRequestMetadata{
		ServiceMethod: "/weather.WeatherService/GetWeatherInfo",
		BodyToLog:     []uint8(`{"date":"M41.993.32"}`),
	}, reqMeta, cmpopts.IgnoreFields(GRPCRequestMetadata{}, "ContextData"))
	require.ElementsMatch(t, []string{"user_id=1234", "role=admin"}, reqMeta.ContextData)

	// 3.3. Simulated end instrument
	require.NotNil(t, end)
	end(nil)

	// 3.4. Validate trace attributes
	expectedAttributes := []attribute.KeyValue{
		semconv.RPCSystemGRPC,
		semconv.NetworkPeerAddress(":50051"),
		semconv.NetworkTransportTCP,
		semconv.RPCService("weather.WeatherService"),
		semconv.RPCMethod("GetWeatherInfo"),
		semconv.RPCGRPCStatusCodeOk,
	}
	require.ElementsMatch(t, expectedAttributes, monitor.GetSpans().Snapshots()[0].Attributes())
}

package monitoring

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
)

func TestStartIncomingRequest(t *testing.T) {
	// 1. Given
	monitor, endTest := NewMonitorTest()
	defer endTest()

	// Create new test request
	request, err := http.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer([]byte(`{"id":1,"name":"the-witcher-knight"}`)))
	require.NoError(t, err)
	request.Header.Set("Traceparent", "00-deadbeefcafebabefeedfacebadc0de1-abad1dea0ddba11c-01")
	request.Header.Set("Tracestate", "test=test-value")
	request.Header.Set("Baggage", "user_id=1234,role=admin")

	// 2. When
	ctx, reqMeta, end := StartIncomingRequest(monitor.GetLogger(), request)

	// 3. Then
	// 3.1. Validate logs
	require.NotNil(t, ctx.Value(loggerContextKey))
	FromContext(ctx).Infof("Got incoming request")
	expectedLogs := []map[string]interface{}{
		{
			"level":    "info",
			"msg":      "Got incoming request",
			"span_id":  "0000000000000001",
			"trace_id": "deadbeefcafebabefeedfacebadc0de1", // Same with incoming request
		},
	}
	requireEqual(t, expectedLogs, monitor.GetLogs(t), cmpopts.IgnoreMapEntries(func(key string, value interface{}) bool {
		return key == "ts" || key == "span_id"
	}))

	// 3.2. Validate request metadata
	requireEqual(t, RequestMetadata{
		Method:    http.MethodPost,
		Endpoint:  "/api/v1/users",
		BodyToLog: []byte(`{"id":1,"name":"the-witcher-knight"}`),
	}, reqMeta, cmpopts.IgnoreFields(RequestMetadata{}, "ContextData"))
	require.ElementsMatch(t, []string{"user_id=1234", "role=admin"}, reqMeta.ContextData)

	// 3.3. Simulated end instrument
	require.NotNil(t, end)
	end(http.StatusOK, nil)

	// 3.4. Validate trace attributes
	expectedAttributes := []attribute.KeyValue{
		semconv.HTTPRequestMethodKey.String("POST"),
		semconv.ServerAddressKey.String(""),
		semconv.UserAgentOriginal(""),
		semconv.URLFull("/api/v1/users"),
		semconv.NetworkPeerAddress(""),
		semconv.NetworkProtocolVersion("HTTP/1.1"),
		semconv.HTTPRequestBodySize(36),
		semconv.HTTPResponseStatusCode(http.StatusOK),
	}
	require.Equal(t, expectedAttributes, monitor.GetSpans().Snapshots()[0].Attributes())
}

func requireEqual[T any](t *testing.T, expected, actual T, opts ...cmp.Option) {
	t.Helper()
	if !cmp.Equal(expected, actual, opts...) {
		t.Errorf("\n mismatched. \n expected: %+v \n got: %+v \n diff: %+v",
			expected, actual,
			cmp.Diff(expected, actual, opts...),
		)
		t.FailNow()
	}
}

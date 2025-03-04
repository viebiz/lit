package instrumenthttp

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/viebiz/lit/monitoring"
	"github.com/viebiz/lit/monitoring/tracing/mocktracer"
	"github.com/viebiz/lit/testutil"
	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

func TestStartIncomingRequest(t *testing.T) {
	tp := mocktracer.Start()
	defer tp.Stop()

	type args struct {
		givenURL             string
		givenMethod          string
		givenHeaders         map[string]string
		givenBody            io.Reader
		givenStatus          int
		givenRespErr         error
		expLogBody           bool
		expSpanKind          trace.SpanKind
		expParentSpanContext trace.SpanContext
		expSpanContext       trace.SpanContext
		expAttributes        []attribute.KeyValue
		expSpanEvents        []sdktrace.Event
	}
	tcs := map[string]args{
		"GET": {
			givenURL:             "/api/v1/users",
			givenMethod:          http.MethodGet,
			givenStatus:          http.StatusOK,
			expSpanKind:          trace.SpanKindServer,
			expParentSpanContext: trace.NewSpanContext(trace.SpanContextConfig{Remote: true}),
			expSpanContext: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    testutil.NewTraceID(t, fmt.Sprintf("%032x", 1)),
				SpanID:     testutil.NewSpanID(t, fmt.Sprintf("%016x", 1)),
				TraceFlags: 01,
			}),
			expAttributes: []attribute.KeyValue{
				semconv.HTTPRequestMethodKey.String(http.MethodGet),
				semconv.ServerAddressKey.String("example.com"),
				semconv.UserAgentOriginal(""),
				semconv.URLFull("/api/v1/users"),
				semconv.NetworkPeerAddress("192.0.2.1:1234"),
				semconv.NetworkProtocolVersion("HTTP/1.1"),
				semconv.HTTPResponseStatusCode(http.StatusOK),
			},
		},
		"GET with body": {
			givenURL:             "/api/v1/users",
			givenMethod:          http.MethodGet,
			givenBody:            bytes.NewBufferString("{}"),
			givenStatus:          http.StatusOK,
			expSpanKind:          trace.SpanKindServer,
			expLogBody:           false,
			expParentSpanContext: trace.NewSpanContext(trace.SpanContextConfig{Remote: true}),
			expSpanContext: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    testutil.NewTraceID(t, fmt.Sprintf("%032x", 1)),
				SpanID:     testutil.NewSpanID(t, fmt.Sprintf("%016x", 1)),
				TraceFlags: 01,
			}),
			expAttributes: []attribute.KeyValue{
				semconv.HTTPRequestMethodKey.String(http.MethodGet),
				semconv.ServerAddressKey.String("example.com"),
				semconv.UserAgentOriginal(""),
				semconv.URLFull("/api/v1/users"),
				semconv.NetworkPeerAddress("192.0.2.1:1234"),
				semconv.NetworkProtocolVersion("HTTP/1.1"),
				semconv.HTTPRequestBodySize(2),
				semconv.HTTPResponseStatusCode(http.StatusOK),
			},
		},
		"GET with error": {
			givenURL:             "/api/v1/users",
			givenMethod:          http.MethodGet,
			givenStatus:          http.StatusBadRequest,
			givenRespErr:         errors.New("simulated bad request"),
			expSpanKind:          trace.SpanKindServer,
			expParentSpanContext: trace.NewSpanContext(trace.SpanContextConfig{Remote: true}),
			expSpanContext: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    testutil.NewTraceID(t, fmt.Sprintf("%032x", 1)),
				SpanID:     testutil.NewSpanID(t, fmt.Sprintf("%016x", 1)),
				TraceFlags: 01,
			}),
			expAttributes: []attribute.KeyValue{
				semconv.HTTPRequestMethodKey.String(http.MethodGet),
				semconv.ServerAddressKey.String("example.com"),
				semconv.UserAgentOriginal(""),
				semconv.URLFull("/api/v1/users"),
				semconv.NetworkPeerAddress("192.0.2.1:1234"),
				semconv.NetworkProtocolVersion("HTTP/1.1"),
				semconv.HTTPResponseStatusCode(http.StatusBadRequest),
			},
			expSpanEvents: []sdktrace.Event{
				{
					Name: "exception",
					Attributes: []attribute.KeyValue{
						semconv.ExceptionType("*errors.errorString"),
						semconv.ExceptionMessage("simulated bad request"),
					},
				},
			},
		},
		"GET with parent spancontext": {
			givenURL:    "/api/v1/users",
			givenMethod: http.MethodGet,
			givenHeaders: map[string]string{
				"Traceparent": "00-deadbeefcafebabefeedfacebadc0de1-abad1dea0ddba11c-01",
				"Tracestate":  "ot=abc123, azure=def456, aws=ghi789",
			},
			givenStatus: http.StatusOK,
			expParentSpanContext: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    testutil.NewTraceID(t, "deadbeefcafebabefeedfacebadc0de1"),
				SpanID:     testutil.NewSpanID(t, "abad1dea0ddba11c"),
				TraceFlags: 01,
				TraceState: testutil.NewTraceState(t, "ot=abc123, azure=def456, aws=ghi789"),
				Remote:     true,
			}),
			expSpanContext: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    testutil.NewTraceID(t, "deadbeefcafebabefeedfacebadc0de1"),
				SpanID:     testutil.NewSpanID(t, fmt.Sprintf("%016x", 1)),
				TraceFlags: 01,
				TraceState: testutil.NewTraceState(t, "ot=abc123, azure=def456, aws=ghi789"),
			}),
			expSpanKind: trace.SpanKindServer,
			expAttributes: []attribute.KeyValue{
				semconv.HTTPRequestMethodKey.String(http.MethodGet),
				semconv.ServerAddressKey.String("example.com"),
				semconv.UserAgentOriginal(""),
				semconv.URLFull("/api/v1/users"),
				semconv.NetworkPeerAddress("192.0.2.1:1234"),
				semconv.NetworkProtocolVersion("HTTP/1.1"),
				semconv.HTTPResponseStatusCode(http.StatusOK),
			},
		},
		"POST": {
			givenURL:             "/api/v1/users",
			givenMethod:          http.MethodPost,
			givenBody:            bodyFromFile(t, "medium.json"),
			givenStatus:          http.StatusOK,
			expLogBody:           true,
			expSpanKind:          trace.SpanKindServer,
			expParentSpanContext: trace.NewSpanContext(trace.SpanContextConfig{Remote: true}),
			expSpanContext: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    testutil.NewTraceID(t, fmt.Sprintf("%032x", 1)),
				SpanID:     testutil.NewSpanID(t, fmt.Sprintf("%016x", 1)),
				TraceFlags: 01,
			}),
			expAttributes: []attribute.KeyValue{
				semconv.HTTPRequestMethodKey.String(http.MethodPost),
				semconv.ServerAddressKey.String("example.com"),
				semconv.UserAgentOriginal(""),
				semconv.URLFull("/api/v1/users"),
				semconv.NetworkPeerAddress("192.0.2.1:1234"),
				semconv.NetworkProtocolVersion("HTTP/1.1"),
				semconv.HTTPRequestBodySize(9752),
				semconv.HTTPResponseStatusCode(http.StatusOK),
			},
		},
		"POST error when read body": {
			givenURL:             "/api/v1/users",
			givenMethod:          http.MethodPost,
			givenBody:            &errorReader{},
			givenStatus:          http.StatusOK,
			expLogBody:           false,
			expSpanKind:          trace.SpanKindServer,
			expParentSpanContext: trace.NewSpanContext(trace.SpanContextConfig{Remote: true}),
			expSpanContext: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    testutil.NewTraceID(t, fmt.Sprintf("%032x", 1)),
				SpanID:     testutil.NewSpanID(t, fmt.Sprintf("%016x", 1)),
				TraceFlags: 01,
			}),
			expAttributes: []attribute.KeyValue{
				semconv.HTTPRequestMethodKey.String(http.MethodPost),
				semconv.ServerAddressKey.String("example.com"),
				semconv.UserAgentOriginal(""),
				semconv.URLFull("/api/v1/users"),
				semconv.NetworkPeerAddress("192.0.2.1:1234"),
				semconv.NetworkProtocolVersion("HTTP/1.1"),
				semconv.HTTPResponseStatusCode(http.StatusOK),
			},
		},
		"POST invalid json": {
			givenURL:             "/api/v1/users",
			givenMethod:          http.MethodPost,
			givenBody:            bytes.NewBufferString("invalid json"),
			givenStatus:          http.StatusOK,
			expLogBody:           false,
			expSpanKind:          trace.SpanKindServer,
			expParentSpanContext: trace.NewSpanContext(trace.SpanContextConfig{Remote: true}),
			expSpanContext: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    testutil.NewTraceID(t, fmt.Sprintf("%032x", 1)),
				SpanID:     testutil.NewSpanID(t, fmt.Sprintf("%016x", 1)),
				TraceFlags: 01,
			}),
			expAttributes: []attribute.KeyValue{
				semconv.HTTPRequestMethodKey.String(http.MethodPost),
				semconv.ServerAddressKey.String("example.com"),
				semconv.UserAgentOriginal(""),
				semconv.URLFull("/api/v1/users"),
				semconv.NetworkPeerAddress("192.0.2.1:1234"),
				semconv.NetworkProtocolVersion("HTTP/1.1"),
				semconv.HTTPRequestBodySize(12),
				semconv.HTTPResponseStatusCode(http.StatusOK),
			},
		},
		"POST body out of limit": {
			givenURL:             "/api/v1/users",
			givenMethod:          http.MethodPost,
			givenBody:            bodyFromFile(t, "large.json"),
			givenStatus:          http.StatusOK,
			expLogBody:           false,
			expSpanKind:          trace.SpanKindServer,
			expParentSpanContext: trace.NewSpanContext(trace.SpanContextConfig{Remote: true}),
			expSpanContext: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    testutil.NewTraceID(t, fmt.Sprintf("%032x", 1)),
				SpanID:     testutil.NewSpanID(t, fmt.Sprintf("%016x", 1)),
				TraceFlags: 01,
			}),
			expAttributes: []attribute.KeyValue{
				semconv.HTTPRequestMethodKey.String(http.MethodPost),
				semconv.ServerAddressKey.String("example.com"),
				semconv.UserAgentOriginal(""),
				semconv.URLFull("/api/v1/users"),
				semconv.NetworkPeerAddress("192.0.2.1:1234"),
				semconv.NetworkProtocolVersion("HTTP/1.1"),
				semconv.HTTPRequestBodySize(10886),
				semconv.HTTPResponseStatusCode(http.StatusOK),
			},
		},
		"POST with error": {
			givenURL:             "/api/v1/users",
			givenMethod:          http.MethodPost,
			givenBody:            bytes.NewBufferString(`{"username":"the-witcher-knight"}`),
			givenStatus:          http.StatusBadRequest,
			givenRespErr:         errors.New("simulated error"),
			expLogBody:           true,
			expSpanKind:          trace.SpanKindServer,
			expParentSpanContext: trace.NewSpanContext(trace.SpanContextConfig{Remote: true}),
			expSpanContext: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    testutil.NewTraceID(t, fmt.Sprintf("%032x", 1)),
				SpanID:     testutil.NewSpanID(t, fmt.Sprintf("%016x", 1)),
				TraceFlags: 01,
			}),
			expAttributes: []attribute.KeyValue{
				semconv.HTTPRequestMethodKey.String(http.MethodPost),
				semconv.ServerAddressKey.String("example.com"),
				semconv.UserAgentOriginal(""),
				semconv.URLFull("/api/v1/users"),
				semconv.NetworkPeerAddress("192.0.2.1:1234"),
				semconv.NetworkProtocolVersion("HTTP/1.1"),
				semconv.HTTPRequestBodySize(33),
				semconv.HTTPResponseStatusCode(http.StatusBadRequest),
			},
			expSpanEvents: []sdktrace.Event{
				{
					Name: "exception",
					Attributes: []attribute.KeyValue{
						semconv.ExceptionType("*errors.errorString"),
						semconv.ExceptionMessage("simulated error"),
					},
				},
			},
		},
		"POST with parent span context": {
			givenURL:    "/api/v1/users",
			givenMethod: http.MethodPost,
			givenHeaders: map[string]string{
				"Traceparent": "00-deadbeefcafebabefeedfacebadc0de1-abad1dea0ddba11c-01",
				"Tracestate":  "ot=abc123, azure=def456, aws=ghi789",
			},
			givenBody:   bytes.NewBufferString(`{"username":"the-witcher-knight"}`),
			givenStatus: http.StatusOK,
			expLogBody:  true,
			expSpanKind: trace.SpanKindServer,
			expParentSpanContext: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    testutil.NewTraceID(t, "deadbeefcafebabefeedfacebadc0de1"),
				SpanID:     testutil.NewSpanID(t, "abad1dea0ddba11c"),
				TraceFlags: 01,
				TraceState: testutil.NewTraceState(t, "ot=abc123, azure=def456, aws=ghi789"),
				Remote:     true,
			}),
			expSpanContext: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    testutil.NewTraceID(t, "deadbeefcafebabefeedfacebadc0de1"),
				SpanID:     testutil.NewSpanID(t, fmt.Sprintf("%016x", 1)),
				TraceFlags: 01,
				TraceState: testutil.NewTraceState(t, "ot=abc123, azure=def456, aws=ghi789"),
			}),
			expAttributes: []attribute.KeyValue{
				semconv.HTTPRequestMethodKey.String(http.MethodPost),
				semconv.ServerAddressKey.String("example.com"),
				semconv.UserAgentOriginal(""),
				semconv.URLFull("/api/v1/users"),
				semconv.NetworkPeerAddress("192.0.2.1:1234"),
				semconv.NetworkProtocolVersion("HTTP/1.1"),
				semconv.HTTPRequestBodySize(33),
				semconv.HTTPResponseStatusCode(http.StatusOK),
			},
		},
	}

	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			//t.Parallel() Should not use parallel because use global variable
			tp.Reset() // Reset span container

			// Given
			r := httptest.NewRequest(tc.givenMethod, tc.givenURL, tc.givenBody)
			for k, v := range tc.givenHeaders {
				r.Header.Add(k, v)
			}

			m, err := monitoring.New(monitoring.Config{Writer: io.Discard})
			require.NoError(t, err)

			// When
			_, reqMeta, end := StartIncomingRequest(m, r)
			end(tc.givenStatus, tc.givenRespErr)

			// Then
			if tc.expLogBody {
				require.NotEmpty(t, reqMeta.BodyToLog)
			}

			spanStub := tp.GetLatestSpan()
			require.Equal(t, tc.expSpanKind, spanStub.SpanKind)
			testutil.Equal(t, tc.expParentSpanContext, spanStub.Parent)
			testutil.Equal(t, tc.expSpanContext, spanStub.SpanContext)
			testutil.Equal(t, tc.expAttributes, spanStub.Attributes, testutil.EquateComparable[[]attribute.KeyValue](attribute.KeyValue{}))
			if len(tc.expSpanEvents) > 0 {
				require.Equal(t, len(tc.expSpanEvents), len(spanStub.Events))
				for idx, expEvent := range tc.expSpanEvents {
					e := spanStub.Events[idx]
					require.Equal(t, expEvent.Name, e.Name)
					testutil.Equal(t, expEvent.Attributes, e.Attributes, testutil.EquateComparable[[]attribute.KeyValue](attribute.KeyValue{}))
				}
			}
		})
	}
}

func bodyFromFile(t *testing.T, name string) io.Reader {
	f, err := os.Open("testdata/" + name)
	require.NoError(t, err)
	defer f.Close()

	b, err := io.ReadAll(f)
	require.NoError(t, err)

	return bytes.NewReader(b)
}

// errorReader is a custom io.Reader that always returns an error
type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("mock read error")
}

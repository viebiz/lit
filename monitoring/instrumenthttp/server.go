package instrumenthttp

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/viebiz/lit/monitoring"
)

func StartIncomingRequest(m *monitoring.Monitor, r *http.Request) (context.Context, RequestMetadata, func(int, error)) {
	logTags := map[string]string{
		httpRequestMethodKey:   r.Method,
		serverAddressKey:       r.Host,
		userAgentKey:           r.UserAgent(),
		urlKey:                 r.URL.Path,
		networkPeerAddressKey:  r.RemoteAddr,
		networkProtocolVersion: r.Proto,
	}

	attrs := []attribute.KeyValue{
		semconv.HTTPRequestMethodKey.String(r.Method),
		semconv.ServerAddressKey.String(r.Host),
		semconv.UserAgentOriginal(r.UserAgent()),
		semconv.URLFull(r.URL.Path),
		semconv.NetworkPeerAddress(r.RemoteAddr),
		semconv.NetworkProtocolVersion(r.Proto),
	}

	ctx := r.Context()

	// Collect request metadata to log
	reqMeta := RequestMetadata{
		Method:   r.Method,
		Endpoint: r.URL.Path,
	}

	// Log request body
	if r.ContentLength > 0 {
		logTags[httpRequestBodySize] = strconv.FormatInt(r.ContentLength, 10)
		attrs = append(attrs, semconv.HTTPRequestBodySize(int(r.ContentLength)))
	}

	if bodyBytes := readRequestBody(m, r); len(bodyBytes) > 0 {
		reqMeta.BodyToLog = bodyBytes
	}

	// Extract trace context from request headers
	curSpanCtx := otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(r.Header))
	spanCtx := trace.SpanContextFromContext(curSpanCtx)

	// Start new span
	ctx, span := tracer.Start(trace.ContextWithRemoteSpanContext(ctx, spanCtx), httpIncomingSpanName,
		trace.WithSpanKind(trace.SpanKindServer),
		trace.WithAttributes(attrs...),
	)
	m = monitoring.InjectTracingInfo(m, span.SpanContext())

	m = m.With(logTags)
	ctx = monitoring.SetInContext(ctx, m)

	return ctx, reqMeta,
		func(status int, err error) {
			if err != nil {
				span.SetStatus(codes.Error, err.Error())
				span.RecordError(err)
			}

			span.SetAttributes(semconv.HTTPResponseStatusCode(status))
			span.End()
		}
}

type RequestMetadata struct {
	Method    string
	Endpoint  string
	BodyToLog []byte
}

func readRequestBody(m *monitoring.Monitor, r *http.Request) []byte {
	if r.ContentLength == 0 {
		return nil
	}

	if r.Method != http.MethodPost && r.Method != http.MethodPut && r.Method != http.MethodPatch {
		return nil
	}

	if r.Header.Get(requestHeaderContentType) != contextTypeJSON {
		return nil
	}

	if r.ContentLength > 10_000 {
		// Quite unlikely that request body JSON payload will be more than this. This max limit already gives ~500 lines
		// of JSON payload.
		return nil
	}

	bodyBytes, err := io.ReadAll(r.Body) // Directly read the body into a byte slice
	if err != nil {
		m.Errorf(err, "failed to read request body")
		return nil
	}

	// Restore request body so it can be read again
	r.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	if !json.Valid(bodyBytes) { // We don't care about invalid JSON for logging
		return nil
	}

	// TODO: redact request body
	return bodyBytes
}

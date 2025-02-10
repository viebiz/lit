package monitoring

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	httpIncomingSpanName = "http.incoming_request"

	shouldLogHTTPRequestBody = true
)

var (
	methodsWithRequestBodyMap = map[string]bool{
		http.MethodPost:  true,
		http.MethodPut:   true,
		http.MethodPatch: true,
	}
)

func StartIncomingRequest(logger *Logger, r *http.Request) (context.Context, RequestMetadata, func(int, error)) {
	logFields := []attribute.KeyValue{
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
	if bodyBytes := readRequestBody(logger, r); len(bodyBytes) > 0 {
		logFields = append(logFields, semconv.HTTPRequestBodySize(len(bodyBytes)))
		reqMeta.BodyToLog = bodyBytes
	}

	// Extract trace context from request headers
	curSpanCtx := getTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(r.Header))
	spanCtx := trace.SpanContextFromContext(curSpanCtx)

	bags := baggage.FromContext(curSpanCtx)
	reqMeta.ContextData = make([]string, bags.Len())
	for idx, kv := range bags.Members() {
		reqMeta.ContextData[idx] = kv.String()
	}

	// Add baggage to the context
	ctx = baggage.ContextWithBaggage(ctx, bags)
	ctx, span := getTracer().Start(trace.ContextWithRemoteSpanContext(ctx, spanCtx), httpIncomingSpanName,
		trace.WithSpanKind(trace.SpanKindServer),
		trace.WithAttributes(logFields...),
	)

	return injectTracingInfo(SetInContext(ctx, logger), span.SpanContext()),
		reqMeta,
		func(status int, err error) {
			if err != nil {
				span.SetStatus(codes.Error, err.Error())
				span.RecordError(err)
			}

			span.SetAttributes(semconv.HTTPResponseStatusCode(status))
			span.End(trace.WithTimestamp(timeNowFunc().UTC()))
		}
}

type RequestMetadata struct {
	Method      string
	Endpoint    string
	BodyToLog   []byte
	ContextData []string
}

func readRequestBody(logger *Logger, r *http.Request) []byte {
	if r.ContentLength == 0 {
		return nil
	}

	if !(shouldLogHTTPRequestBody && methodsWithRequestBodyMap[r.Method]) {
		return nil
	}

	if r.ContentLength > 10_000 {
		// Quite unlikely that request body JSON payload will be more than this. This max limit already gives ~500 lines
		// of JSON payload.
		return nil
	}

	bodyBytes := bytes.NewBuffer(nil)
	if _, err := bodyBytes.ReadFrom(r.Body); err != nil {
		logger.Errorf(err, "failed to read request body")
	}

	if !json.Valid(bodyBytes.Bytes()) { // We don't care about invalid JSON for logging
		return nil
	}

	// Rewind request body, that allows read again
	r.Body = io.NopCloser(bodyBytes)

	// TODO: redact request body
	return bodyBytes.Bytes()
}

package monitoring

import (
	"context"
	"strings"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

const (
	unaryIncomingSpanName = "grpc.unary_incoming_call"

	shouldLogUnaryRequestBody = true
)

var (
	timeNowFunc = time.Now
)

func StartUnaryIncomingCall(ctx context.Context, fullMethod string, req any) (context.Context, GRPCRequestMetadata, func(error)) {
	// Init log fields
	logFields := []attribute.KeyValue{
		semconv.RPCSystemGRPC,
	}

	if pr, ok := peer.FromContext(ctx); ok {
		logFields = append(logFields,
			semconv.NetworkPeerAddress(pr.Addr.String()),
			semconv.NetworkTransportKey.String(pr.Addr.Network()),
		)
	}

	if svc, m := extractFullMethod(fullMethod); m != "" {
		logFields = append(logFields,
			semconv.RPCService(svc),
			semconv.RPCMethod(m),
		)
	}

	reqMeta := GRPCRequestMetadata{
		ServiceMethod: fullMethod,
	}

	// Log request body
	if shouldLogUnaryRequestBody {
		reqMeta.BodyToLog = serializeProtoRequest(req)
	}

	// Extract metadata from incoming context
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.MD{}
	} else {
		md = md.Copy() // because it's not safe to modify
	}

	// Extract span context from metadata
	curSpanCtx := getTextMapPropagator().Extract(ctx, mdCarrier(md))
	spanCtx := trace.SpanContextFromContext(curSpanCtx)
	bags := baggage.FromContext(curSpanCtx)
	reqMeta.ContextData = make([]string, bags.Len())
	for idx, kv := range bags.Members() {
		reqMeta.ContextData[idx] = kv.String()
	}

	// Add baggage to the context (ensures the baggage is passed along with the context)
	ctx = baggage.ContextWithBaggage(ctx, bags)
	ctx, span := getTracer().Start(trace.ContextWithRemoteSpanContext(ctx, spanCtx), unaryIncomingSpanName,
		trace.WithSpanKind(trace.SpanKindServer),
		trace.WithAttributes(logFields...),
	)

	return injectTracingInfo(ctx, span.SpanContext()),
		reqMeta,
		func(err error) {
			if err == nil {
				span.SetStatus(codes.Ok, "")
				span.SetAttributes(semconv.RPCGRPCStatusCodeOk)
			} else {
				span.SetStatus(codes.Error, err.Error())
				errStatus := status.Convert(err)
				span.SetAttributes(semconv.RPCGRPCStatusCodeKey.Int(int(errStatus.Code())))
			}

			span.End(trace.WithTimestamp(timeNowFunc().UTC()))
		}
}

type GRPCRequestMetadata struct {
	ServiceMethod string
	BodyToLog     []byte
	ContextData   []string
}

// extractFullMethod extracts full method /weather.WeatherService/GetWeatherInfo
func extractFullMethod(fullMethod string) (string, string) {
	parts := strings.Split(fullMethod, "/")
	if len(parts) == 3 {
		return parts[1], parts[2]
	}

	return parts[0], ""
}

// serializeProtoRequest converts protobuf request to JSON bytes
// output may be unstable due to known issues: https://github.com/golang/protobuf/issues/1121
func serializeProtoRequest(req any) []byte {
	msg, ok := req.(proto.Message)
	if !ok {
		return nil // Ignore req body if it not proto message
	}

	b, err := protojson.Marshal(msg)
	if err != nil {
		return nil // Ignore if it is invalid proto message
	}

	return b
}

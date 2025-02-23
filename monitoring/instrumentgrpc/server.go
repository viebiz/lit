package instrumentgrpc

import (
	"context"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/viebiz/lit/monitoring"
)

func StartUnaryIncomingCall(ctx context.Context, m *monitoring.Monitor, fullMethod string, req any) (context.Context, RequestMetadata, func(error)) {
	// Init log fields
	logTags := map[string]string{
		rpcSystemKey: "grpc",
	}

	attrs := []attribute.KeyValue{
		semconv.RPCSystemGRPC,
	}

	if pr, ok := peer.FromContext(ctx); ok {
		logTags[networkPeerAddressKey] = pr.Addr.String()
		logTags[networkTransportKey] = pr.Addr.Network()

		attrs = append(attrs,
			semconv.NetworkPeerAddress(pr.Addr.String()),
			semconv.NetworkTransportKey.String(pr.Addr.Network()),
		)
	}

	if svc, method := extractFullMethod(fullMethod); method != "" {
		logTags[rpcServiceKey] = svc
		logTags[rpcMethodKey] = method

		attrs = append(attrs,
			semconv.RPCService(svc),
			semconv.RPCMethod(method),
		)
	}

	reqMeta := RequestMetadata{
		ServiceMethod: fullMethod,
	}

	// Log request body
	if shouldLogUnaryRequestBody {
		reqMeta.BodyToLog = serializeProtoMessage(req)
	}

	// Extract metadata from incoming context
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.MD{}
	} else {
		md = md.Copy() // because it's not safe to modify
	}

	// Extract span context from metadata
	curSpanCtx := otel.GetTextMapPropagator().Extract(ctx, mdCarrier(md))
	spanCtx := trace.SpanContextFromContext(curSpanCtx)

	ctx, span := tracer.Start(trace.ContextWithRemoteSpanContext(ctx, spanCtx), unaryIncomingSpanName,
		trace.WithSpanKind(trace.SpanKindServer),
		trace.WithAttributes(attrs...),
	)
	m = monitoring.InjectTracingInfo(m, span.SpanContext())

	m = m.With(logTags)
	ctx = monitoring.SetInContext(ctx, m)

	return ctx,
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

			span.End()
		}
}

type RequestMetadata struct {
	ServiceMethod string
	BodyToLog     []byte
}

// extractFullMethod extracts full method /weather.WeatherService/GetWeatherInfo
func extractFullMethod(fullMethod string) (string, string) {
	parts := strings.Split(fullMethod, "/")
	if len(parts) == 3 {
		return parts[1], parts[2]
	}

	return parts[0], ""
}

// serializeProtoMessage converts protobuf request to JSON bytes
// output may be unstable due to known issues: https://github.com/golang/protobuf/issues/1121
func serializeProtoMessage(req any) []byte {
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

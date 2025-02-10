package monitoring

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
)

const (
	unaryOutgoingCallSpanName = "grpc.unary_outgoing_call"
)

func StartGRPCUnaryCallSegment(ctx context.Context, svcInfo ExternalServiceInfo, fullMethod string) (context.Context, func(error)) {
	attrs := []attribute.KeyValue{
		semconv.RPCSystemGRPC,
		semconv.ServerAddress(svcInfo.Hostname + ":" + svcInfo.Port),
	}

	if svc, m := extractFullMethod(fullMethod); m != "" {
		attrs = append(attrs,
			semconv.RPCService(svc),
			semconv.RPCMethod(m),
		)
	}

	ctx, span := getTracer().Start(ctx, unaryOutgoingCallSpanName, trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(attrs...),
	)

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.MD{}
	} else {
		md = md.Copy() // we have to copy the metadata because it's not safe to modify
	}

	propagator := otel.GetTextMapPropagator()
	propagator.Inject(ctx, mdCarrier(md))
	ctx = metadata.NewOutgoingContext(ctx, md)
	// ? Baggage

	ctx = injectOutgoingTracingInfo(ctx, span.SpanContext())

	return ctx, func(err error) {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err, trace.WithStackTrace(true))
		}

		span.End()
	}
}

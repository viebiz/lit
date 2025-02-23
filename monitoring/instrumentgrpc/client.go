package instrumentgrpc

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"

	"github.com/viebiz/lit/monitoring"
)

func StartGRPCUnaryCallSegment(ctx context.Context, svcInfo monitoring.ExternalServiceInfo, fullMethod string) (context.Context, func(error)) {
	logTags := map[string]string{
		rpcSystemKey:     "grpc",
		serverAddressKey: svcInfo.Hostname + ":" + svcInfo.Port,
	}

	attrs := []attribute.KeyValue{
		semconv.RPCSystemGRPC,
		semconv.ServerAddress(svcInfo.Hostname + ":" + svcInfo.Port),
	}

	if svc, method := extractFullMethod(fullMethod); method != "" {
		logTags[rpcServiceKey] = svc
		logTags[rpcMethodKey] = method

		attrs = append(attrs,
			semconv.RPCService(svc),
			semconv.RPCMethod(method),
		)
	}

	ctx, span := tracer.Start(ctx, unaryOutgoingCallSpanName, trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(attrs...),
	)

	// Copy metadata to new context
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.MD{}
	} else {
		md = md.Copy() // we have to copy the metadata because it's not safe to modify
	}

	otel.GetTextMapPropagator().Inject(ctx, mdCarrier(md))
	ctx = metadata.NewOutgoingContext(ctx, md)
	// ? Baggage

	m := monitoring.InjectOutgoingTracingInfo(monitoring.FromContext(ctx), span.SpanContext())
	m = m.With(logTags)
	ctx = monitoring.SetInContext(ctx, m)

	return ctx, func(err error) {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err, trace.WithStackTrace(true))
		}

		span.End()
	}
}

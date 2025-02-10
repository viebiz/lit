package grpcclient

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/viebiz/lit/monitoring"
)

func unaryClientInterceptor(
	ctx context.Context,
	method string,
	req, reply any,
	clientConn *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) (err error) {
	var extSvcInfo externalServiceInfoOption
	for _, opt := range opts {
		if v, ok := opt.(externalServiceInfoOption); ok {
			extSvcInfo = v
			continue
		}
	}

	ctx, end := monitoring.StartGRPCUnaryCallSegment(ctx, extSvcInfo.info, method)
	defer func() {
		end(err)
	}()

	logRequestBody(ctx, req)

	if err = invoker(ctx, method, req, reply, clientConn, opts...); err != nil {
		return err
	}

	return nil
}

func logRequestBody(ctx context.Context, req interface{}) {
	monitoring.FromContext(ctx).
		With(monitoring.Field("grpc.request", parseProtoRequest(req))).
		Infof("grpc.outgoing_request")
}

func parseProtoRequest(req any) []byte {
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

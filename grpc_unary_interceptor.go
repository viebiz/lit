package lit

import (
	"context"
	"fmt"
	"runtime/debug"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/viebiz/lit/monitoring"
	"github.com/viebiz/lit/monitoring/instrumentgrpc"
)

const (
	shouldLogGRPCResponse = true
)

func unaryServerInterceptor(rootCtx context.Context) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (rs any, err error) {
		// Start tracing for incoming unary call request
		ctx, reqMeta, endInstrumentation := instrumentgrpc.StartUnaryIncomingCall(ctx, monitoring.FromContext(rootCtx), info.FullMethod, req)
		defer func() {
			if p := recover(); p != nil {
				rcvErr, ok := p.(error)
				if !ok {
					rcvErr = fmt.Errorf("%v", p)
				}

				monitoring.FromContext(ctx).Errorf(rcvErr, "Caught a panic: %s", debug.Stack())
				endInstrumentation(rcvErr)

				err = ErrDefaultInternal
			}
		}()

		rs, err = handler(ctx, req)

		endInstrumentation(err)

		logIncomingGRPCCall(ctx, reqMeta, rs)

		return rs, err
	}
}

func parseProtoMessage(resp any) []byte {
	msg, ok := resp.(proto.Message)
	if !ok {
		return nil // Ignore req body if it not proto message
	}

	b, err := protojson.Marshal(msg)
	if err != nil {
		return nil // Ignore if it is invalid proto message
	}

	return b
}

func logIncomingGRPCCall(ctx context.Context, reqMeta instrumentgrpc.RequestMetadata, result any) {
	//logFields := []monitoring.LogField{
	//	monitoring.Field("grpc.service_method", reqMeta.ServiceMethod),
	//}
	//
	//// BodyToLog always have `{}`
	//if len(reqMeta.BodyToLog) > 2 {
	//	logFields = append(logFields, monitoring.Field("grpc.request_body", reqMeta.BodyToLog))
	//}
	//
	//if resultToLog := parseProtoMessage(result); len(resultToLog) > 2 && shouldLogGRPCResponse {
	//	logFields = append(logFields, monitoring.Field("grpc.response_body", parseProtoMessage(result)))
	//}
	//
	//monitoring.FromContext(ctx).
	//	With(logFields...).
	//	Infof("grpc.unary_incoming_call")
}

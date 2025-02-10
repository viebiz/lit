package lit

import (
	"context"

	"google.golang.org/grpc"
)

func streamServerInterceptor(ctx context.Context) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// TODO: implement logic

		return handler(srv, ss)
	}
}

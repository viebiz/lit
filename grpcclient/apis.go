package grpcclient

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// NewUnauthenticatedConnection initializes and returns a new unauthenticated grpc clientConn for unary calls
func NewUnauthenticatedConnection(ctx context.Context, addr string) (Conn, error) {
	return initUnaryClient(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
}

package grpcclient

import (
	"context"

	pkgerrors "github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/viebiz/lit/monitoring"
)

// Conn defines a gRPC unary client connection interface.
type Conn interface {
	Invoke(ctx context.Context, method string, args any, reply any, opts ...grpc.CallOption) error

	NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error)
}

// NewUnauthenticatedConnection initializes and returns a new unauthenticated grpc clientConn for unary calls
func NewUnauthenticatedConnection(ctx context.Context, addr string) (Conn, error) {
	return initUnaryClient(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
}

func initUnaryClient(ctx context.Context, addr string, opts ...grpc.DialOption) (Conn, error) {
	svcInfo := monitoring.NewExternalServiceInfo(addr)

	conn, err := grpc.NewClient(addr,
		append(
			commonUnaryClientDialOptions(svcInfo),
			opts...,
		)...,
	)
	if err != nil {
		return nil, pkgerrors.WithStack(err)
	}

	return &clientConn{
		conn: conn,
	}, nil
}

type clientConn struct {
	conn *grpc.ClientConn
}

func (u clientConn) Invoke(ctx context.Context, method string, args any, reply any, opts ...grpc.CallOption) error {
	return u.conn.Invoke(ctx, method, args, reply, opts...)
}

func (u clientConn) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return u.conn.NewStream(ctx, desc, method, opts...)
}

func commonUnaryClientDialOptions(svcInfo monitoring.ExternalServiceInfo) []grpc.DialOption {
	return []grpc.DialOption{
		// Explicitly disabling this as according to doc: Retry support is currently disabled by default, but will be enabled by default in the future.
		grpc.WithDisableRetry(),
		grpc.WithDefaultCallOptions(
			externalServiceInfoOption{info: svcInfo}, // Pass service information for tracing.
		),
		grpc.WithChainUnaryInterceptor(unaryClientInterceptor),
	}
}

// externalServiceInfoOption to keeps the external service info in UnaryClient for purpose monitor
type externalServiceInfoOption struct {
	grpc.EmptyCallOption
	info monitoring.ExternalServiceInfo
}

package lightning

import (
	"context"
	"crypto/tls"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type GRPCOption func(option *[]grpc.ServerOption)

func WithTLSConfig(tlsConfig *tls.Config) GRPCOption {
	return func(opts *[]grpc.ServerOption) {
		*opts = append(*opts, grpc.Creds(credentials.NewTLS(tlsConfig)))
	}
}

func WithDefaultInterceptors(ctx context.Context) GRPCOption {
	return func(opts *[]grpc.ServerOption) {
		*opts = append(*opts,
			grpc.ChainUnaryInterceptor(unaryServerInterceptor(ctx)),
			//grpc.ChainStreamInterceptor(streamServerInterceptor(ctx)), // TODO: Implement later
		)
	}
}

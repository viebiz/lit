package otel

import (
	"context"
	"crypto/tls"
	"fmt"

	pkgerrors "github.com/pkg/errors"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func newGRPCExporter(
	ctx context.Context,
	addr string,
	useTLS bool,
	tlsConfig *tls.Config,
) (*otlptrace.Exporter, error) {
	creds := insecure.NewCredentials()
	if useTLS {
		creds = credentials.NewTLS(tlsConfig)
	}

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, pkgerrors.Wrap(err, "create grpc client")
	}

	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, pkgerrors.Wrap(err, "create exporter")
	}

	return exporter, nil
}

func newHTTPExporter(
	ctx context.Context,
	addr string,
	useTLS bool,
	tlsConfig *tls.Config,
) (*otlptrace.Exporter, error) {
	opts := []otlptracehttp.Option{
		otlptracehttp.WithEndpoint(addr),
		otlptracehttp.WithInsecure(),
	}

	if useTLS {
		opts = append(opts, otlptracehttp.WithTLSClientConfig(tlsConfig))
	}

	exporter, err := otlptracehttp.New(ctx, opts...)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "create exporter")
	}

	return exporter, nil
}

func createExporter(ctx context.Context, cfg config) (*otlptrace.Exporter, error) {
	switch cfg.TransportType {
	case TransportHTTP:
		return newHTTPExporter(ctx, cfg.ExporterURL, cfg.UseTLS, cfg.TLSConfig)
	case TransportGRPC:
		return newGRPCExporter(ctx, cfg.ExporterURL, cfg.UseTLS, cfg.TLSConfig)
	default:
		return nil, fmt.Errorf("unknown transport type: %s", cfg.TransportType)
	}
}

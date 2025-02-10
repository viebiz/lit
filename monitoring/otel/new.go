package otel

import (
	"context"
	"crypto/tls"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
)

// config holds the configuration for setting up OpenTelemetry components.
type config struct {
	ExporterURL   string               // URL of the OTLP exporter
	TransportType TransportType        // Transport type: HTTP or gRPC
	UseTLS        bool                 // Whether to use TLS
	TLSConfig     *tls.Config          // TLS configuration for secure transport
	ExtraAttrs    []attribute.KeyValue // Additional resource attributes
}

func defaultConfig(url string) config {
	return config{
		ExporterURL:   url,
		TransportType: TransportGRPC,
		UseTLS:        false,
	}
}

func Setup(ctx context.Context, url string, opts ...ExporterOption) error {
	cfg := defaultConfig(url)
	for _, opt := range opts {
		opt(&cfg)
	}

	if err := cfg.validate(); err != nil {
		return err
	}

	exporter, err := createExporter(ctx, cfg)
	if err != nil {
		return err
	}

	rsc, err := buildResource(cfg)
	if err != nil {
		return fmt.Errorf("build resource error: %w", err)
	}

	// Configure the trace provider
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(rsc),
		sdktrace.WithSpanProcessor(sdktrace.NewBatchSpanProcessor(
			exporter,
			// TODO: Enable to limit resources
			//		DefaultMaxQueueSize       = 2048
			//		DefaultScheduleDelay      = 5000
			//		DefaultExportTimeout      = 30000
			//		DefaultMaxExportBatchSize = 512
			//sdktrace.WithMaxExportBatchSize(),
			//sdktrace.WithMaxQueueSize(),
		)),
	)
	otel.SetTracerProvider(tracerProvider)

	// Setup propagators for trace context and baggage
	propagators := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, // Support W3C TraceContext.
		propagation.Baggage{},      // Support baggage propagation.
	)
	otel.SetTextMapPropagator(propagators)

	return nil
}

func (cfg config) validate() error {
	if cfg.ExporterURL == "" {
		return fmt.Errorf("exporter URL cannot be empty")
	}

	if !cfg.TransportType.IsValid() {
		return fmt.Errorf("invalid transport type: %q", cfg.TransportType)
	}

	if cfg.UseTLS && cfg.TLSConfig == nil {
		return fmt.Errorf("TLS is enabled but TLSConfig is not provided")
	}

	return nil
}

// buildResource constructs a Resource object with the default and extra attributes.
// The resource combines default attributes with user-provided ones.
func buildResource(cfg config) (*resource.Resource, error) {
	if len(cfg.ExtraAttrs) == 0 {
		return resource.Default(), nil
	}

	return resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			cfg.ExtraAttrs...,
		),
	)
}

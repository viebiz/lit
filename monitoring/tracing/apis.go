package tracing

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// Init setups Distributed Tracing service
func Init(ctx context.Context, cfg Config) error {
	// Use GRPC as default config
	if cfg.TransportType == "" {
		cfg.TransportType = TransportGRPC
	}

	// Validate configs
	if err := validateConfig(cfg); err != nil {
		return err
	}

	exporter, err := createExporter(ctx, cfg)
	if err != nil {
		return err
	}

	rsc := buildResource()

	// Configure the trace provider
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(rsc),
		sdktrace.WithSpanProcessor(sdktrace.NewBatchSpanProcessor(
			exporter,
			// TODO: Enable to limit resources, this is current configs
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
		// W3C TraceContext: Propagates trace IDs across services for distributed tracing.
		propagation.TraceContext{},
		// TODO: Support baggage in future
		// OpenTelemetry Baggage: Propagates contextual metadata (key-value pairs) across services.
		//propagation.Baggage{},
	)
	otel.SetTextMapPropagator(propagators)

	return nil
}

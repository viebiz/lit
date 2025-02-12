package monitoring

import (
	"context"
	"fmt"

	"github.com/viebiz/lit/monitoring/otel"
)

type MonitorConfig struct {
	Tags map[string]string

	ExporterURL string
}

// Setup initializes logging and tracing tools
func Setup(ctx context.Context, cfg MonitorConfig) (context.Context, error) {
	// Setup logger
	logger := NewLogger(WithFieldFromMap(cfg.Tags))

	// Integrate tracing with OTel
	if err := otel.Setup(ctx, cfg.ExporterURL); err != nil {
		return ctx, fmt.Errorf("setup opentelemetry %w", err)
	}

	// TODO: Add Sentry integration

	return SetInContext(ctx, logger), nil
}

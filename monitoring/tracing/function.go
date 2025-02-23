package tracing

import (
	"fmt"

	"go.opentelemetry.io/otel/sdk/resource"
)

func validateConfig(cfg Config) error {
	if cfg.ExporterURL == "" {
		return ErrMissingExporterURL
	}

	if !cfg.TransportType.IsValid() {
		return fmt.Errorf("%w: %s", ErrInvalidTransportType, cfg.TransportType)
	}

	return nil
}

// buildResource constructs a Resource object using OpenTelemetry defaults.
// It automatically includes environment-based attributes such as:
// - OTEL_SERVICE_NAME to set the service name.
// - OTEL_RESOURCE_ATTRIBUTES for additional resource attributes.
// The function returns the default OpenTelemetry resource configuration.
func buildResource() *resource.Resource {
	return resource.Default()
}

package otel

import (
	"crypto/tls"

	"go.opentelemetry.io/otel/attribute"
)

type ExporterOption func(*config)

// WithTransportType is option to specify the Exporter transport protocol
// Supported HTTP and gRPC (default)
func WithTransportType(t TransportType) ExporterOption {
	return func(cfg *config) {
		cfg.TransportType = t
	}
}

// WithTLS is option to enable TLS for the exporter
func WithTLS(tlsConfig *tls.Config) ExporterOption {
	return func(cfg *config) {
		cfg.UseTLS = true
		cfg.TLSConfig = tlsConfig
	}
}

// WithAttributes is option to add additional attributes to the resource metadata
func WithAttributes(attrs map[string]string) ExporterOption {
	return func(cfg *config) {
		resourceAttrs := make([]attribute.KeyValue, len(attrs))
		for k, v := range attrs {
			resourceAttrs = append(resourceAttrs, attribute.String(k, v))
		}

		cfg.ExtraAttrs = resourceAttrs
	}
}

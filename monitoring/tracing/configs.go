package tracing

import (
	"crypto/tls"
)

// Config holds the configuration for setting up OpenTelemetry components.
type Config struct {
	ExporterURL   string        // URL of the OTLP exporter
	TransportType TransportType // Transport type: HTTP or gRPC
	TLSConfig     *tls.Config   // TLS configuration for secure transport
}

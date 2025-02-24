package tracing

// TransportType represents OpenTelemetry exporter transport type
type TransportType string

const (
	TransportGRPC TransportType = "grpc"
	TransportHTTP TransportType = "http"
)

func (t TransportType) String() string {
	return string(t)
}

func (t TransportType) IsValid() bool {
	return t == TransportGRPC || t == TransportHTTP
}

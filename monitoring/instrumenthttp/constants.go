package instrumenthttp

import (
	"go.opentelemetry.io/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	tracerName           = "github.com/viebiz/lit/monitoring/instrumenthttp"
	httpIncomingSpanName = "http.incoming_request"

	// Attributes
	httpRequestMethodKey   = "http.request.method"
	serverAddressKey       = "server.address"
	userAgentKey           = "user_agent"
	urlKey                 = "url"
	networkPeerAddressKey  = "network.peer.address"
	networkProtocolVersion = "network.protocol.version"
	httpRequestBodySize    = "http.request.body.size"
)

var (
	tracer = otel.Tracer(tracerName, trace.WithSchemaURL(semconv.SchemaURL))
)

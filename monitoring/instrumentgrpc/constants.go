package instrumentgrpc

import (
	"go.opentelemetry.io/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	tracerName                = "github.com/viebiz/lit/monitoring/instrumentgrpc"
	unaryOutgoingCallSpanName = "grpc.unary_outgoing_call"
	unaryIncomingSpanName     = "grpc.unary_incoming_call"

	// Settings
	shouldLogUnaryRequestBody = true

	// Attributes
	rpcSystemKey          = "rpc.system"
	serverAddressKey      = "server.address"
	rpcServiceKey         = "rpc.service"
	rpcMethodKey          = "rpc.method"
	networkPeerAddressKey = "network.peer.address"
	networkTransportKey   = "network.transport"
)

var (
	tracer = otel.Tracer(tracerName, trace.WithSchemaURL(semconv.SchemaURL))
)

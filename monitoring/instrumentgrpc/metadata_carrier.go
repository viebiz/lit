package instrumentgrpc

import (
	"strings"

	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc/metadata"
)

// mdCarrier enabling TraceContext propagation via gRPC metadata
type mdCarrier metadata.MD

// Ensure mdCarrier implement TextMapCarrier interface
var _ propagation.TextMapCarrier = (*mdCarrier)(nil)

func (mdc mdCarrier) Get(key string) string {
	if vals := mdc[key]; len(vals) > 0 {
		return vals[0]
	}

	return ""
}

func (mdc mdCarrier) Set(key string, value string) {
	k := strings.ToLower(key) // as per google.golang.org/grpc/metadata/metadata.go
	mdc[k] = append(mdc[k], value)
}

func (mdc mdCarrier) Keys() []string {
	keys := make([]string, 0, len(mdc))
	for k := range mdc {
		keys = append(keys, k)
	}

	return keys
}

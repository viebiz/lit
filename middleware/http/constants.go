package http

// Exported values to allow easy customization of HTTP response headers in the future
var (
	// RequestIDHeaderName represents x-request-id key response header
	RequestIDHeaderName = "x-request-id"
)

const (
	httpRequestIDKey string = "http.request.id"
)

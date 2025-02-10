package http

import (
	"github.com/viebiz/lit"
	"github.com/viebiz/lit/monitoring"
)

// RequestIDMiddleware ensures each request has a unique Request ID.
// If the Request ID is provided in the request header, it uses that;
// otherwise, it generates a new one and injects it into the request context.
func RequestIDMiddleware() lightning.HandlerFunc {
	return func(c lightning.Context) {
		// Get request ID from header, if it not exists, generate a new one
		requestID := c.Request().Header.Get(RequestIDHeaderName)
		if requestID == "" {
			requestID = idFunc()
		}

		// Inject request ID to request context
		ctx := c.Request().Context()
		ctx = monitoring.InjectField(ctx, httpRequestIDKey, requestID)

		// Update the request context
		c.SetRequestContext(ctx)

		// Add request ID to response header
		c.Header(RequestIDHeaderName, requestID)

		// Continue handle request
		c.Next()
	}
}

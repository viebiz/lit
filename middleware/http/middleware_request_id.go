package http

import (
	"fmt"
	"hash/fnv"

	"github.com/viebiz/lit"
	"github.com/viebiz/lit/monitoring"
)

var (
	hash64Func = hash64
)

// RequestIDMiddleware ensures each request has a unique Request ID.
// If the Request ID is provided in the request header, it uses that;
// otherwise, it generates a new one and injects it into the request context.
func RequestIDMiddleware() lit.HandlerFunc {
	return func(c lit.Context) {
		// Get request ID from header, if it not exists, generate a new one
		requestID := c.Request().Header.Get(headerXRequestID)
		if requestID == "" {
			requestID = idFunc()
		}

		// Inject request ID to request context
		ctx := c.Request().Context()
		ctx = monitoring.InjectField(ctx, httpRequestIDKey, hash64Func(requestID))

		// Update the request context
		c.SetRequestContext(ctx)

		// Continue handle request
		c.Next()

		// Add request ID to response header
		c.Writer().Header().Add(headerXRequestID, requestID)
	}
}

func hash64(plain string) string {
	h := fnv.New64a()
	if _, err := h.Write([]byte(plain)); err != nil {
		return plain
	}

	return fmt.Sprintf("%x", h.Sum64())
}

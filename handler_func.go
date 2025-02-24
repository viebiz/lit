package lit

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/viebiz/lit/monitoring"
)

// HandlerFunc represents a lightning handler function
type HandlerFunc func(ctx Context)

// ErrHandlerFunc represents a lightning handler error function
type ErrHandlerFunc func(ctx Context) error

// WrapF is a helper function for wrapping http.HandlerFunc and returns a ErrHandlerFunc.
func WrapF(f http.HandlerFunc) HandlerFunc {
	return func(c Context) {
		f(c.Writer(), c.Request())
	}
}

// WrapH is a helper function for wrapping http.Handler and returns a Gin middleware.
func WrapH(h http.Handler) HandlerFunc {
	return func(c Context) {
		h.ServeHTTP(c.Writer(), c.Request())
	}
}

func handleUnexpectedError(handler ErrHandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		if err := handler(litContext{Context: c}); err != nil {
			respondJSON(ctx, c.Writer, err)

			// If error is bad request, can skip log the following error
			if werr, ok := err.(Error); ok {
				if werr.StatusCode() < http.StatusInternalServerError || werr.StatusCode() == http.StatusServiceUnavailable {
					return
				}
			}

			monitoring.FromContext(ctx).Errorf(err, "handle error")
			monitoring.NotifyErrorToInstrumentation(ctx, err)
		}
	}
}

func respondJSON(ctx context.Context, w ResponseWriter, obj any) {
	respondJSONWithHeaders(ctx, w, nil, obj)
}

func respondJSONWithHeaders(ctx context.Context, w ResponseWriter, headers map[string]string, obj any) {
	// Set HTTP headers
	w.Header().Set("Content-Type", "application/json")
	for h, v := range headers {
		w.Header().Set(h, v)
	}

	status := http.StatusOK
	var respBytes []byte
	var err error

	switch parsed := obj.(type) {
	case Error:
		if parsed.StatusCode() >= http.StatusInternalServerError && parsed.StatusCode() != http.StatusServiceUnavailable {
			status = http.StatusInternalServerError
			parsed = ErrDefaultInternal
		}
		status = parsed.StatusCode()
		respBytes, err = json.Marshal(parsed)
	case error:
		status = http.StatusInternalServerError
		respBytes, err = json.Marshal(ErrDefaultInternal)
	default:
		respBytes, err = json.Marshal(obj)
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		monitoring.FromContext(ctx).Errorf(err, "Marshal failed")
		return
	}

	// Write response
	w.WriteHeader(status)
	if _, err = w.Write(respBytes); err != nil {
		monitoring.FromContext(ctx).Errorf(err, "Write failed")
	}
}

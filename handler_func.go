package lit

import (
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

func wrapErrHandler(errHandlerFunc ErrHandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := litContext{Context: c}
		if err := errHandlerFunc(ctx); err != nil {
			ctx.AbortWithError(err)

			// Can skip log the following error if error is bad request or service unavailable
			if werr, ok := err.(Error); ok {
				if werr.StatusCode() < http.StatusInternalServerError || werr.StatusCode() == http.StatusServiceUnavailable {
					return
				}
			}

			monitoring.FromContext(ctx).Errorf(err, "got unexpected error")
			monitoring.NotifyErrorToInstrumentation(ctx, err)
		}
	}
}

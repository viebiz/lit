package lit

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/viebiz/lit/monitoring"
)

// HandlerFunc represents a lightning handler function
type HandlerFunc func(ctx Context)

// ErrHandlerFunc represents a lightning handler error function
type ErrHandlerFunc func(ctx Context) error

func handleUnexpectedError(handler ErrHandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := handler(litContext{Context: c}); err != nil {
			var expErr ExpectedError
			if errors.As(err, &expErr) {
				c.AbortWithStatusJSON(expErr.StatusCode(), expErr)
				return
			}

			// Response internal server error
			c.AbortWithStatusJSON(http.StatusInternalServerError, ErrInternalServerError)

			// Capture error
			monitoring.FromContext(c.Request.Context()).Errorf(err, "Got internal server error")
		}
	}
}

package lit

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"
	"strconv"

	"github.com/viebiz/lit/monitoring"
	"github.com/viebiz/lit/monitoring/instrumenthttp"
)

// rootMiddleware is a middleware function that handles tracing for incoming requests
// and recovers from any panics that may occur during request handling
func rootMiddleware(rootCtx context.Context) HandlerFunc {
	return func(c Context) {
		// Start tracing for the incoming request
		ctx, reqMeta, endInstrumentation := instrumenthttp.StartIncomingRequest(monitoring.FromContext(rootCtx), c.Request())
		defer func() {
			// Recover from any panic that may have occurred during request handling
			if p := recover(); p != nil {
				// Check if the panic value is an error; if not, format it as one
				err, ok := p.(error)
				if !ok {
					err = fmt.Errorf("%+v", p)
				}

				// Log the panic details and stack trace using the tracer
				// We use c.Request.Context() as the tracing context may have been modified during the request.
				monitoring.FromContext(c.Request().Context()).Errorf(err, "Caught a panic: %s", debug.Stack())

				// Abort the request with a 500 Internal Server Error response.
				c.AbortWithError(ErrDefaultInternal)
				// End the instrumentation, marking the request with a 500 status code and the error.
				endInstrumentation(http.StatusInternalServerError, err)
			}
		}()

		// Set instrument context to request context
		c.SetRequestContext(ctx)

		// Wrap response writer to inject trace information
		c.SetWriter(wrapWriter(ctx, c.Writer()))

		// Continue handle request
		c.Next()

		// End instrumentation and log
		endInstrumentation(c.Writer().Status(), nil)

		logIncomingRequest(c, reqMeta, "http.incoming_request")
	}
}

type responseRecorder struct {
	ResponseWriter

	ctx context.Context
}

func wrapWriter(ctx context.Context, w ResponseWriter) ResponseWriter {
	return &responseRecorder{ResponseWriter: w, ctx: ctx}
}

func (w *responseRecorder) Write(resp []byte) (n int, err error) {
	defer func() {
		if err != nil {
			monitoring.FromContext(w.ctx).Errorf(err, "Failed to write response")
		} else {
			monitoring.FromContext(w.ctx).Infof("Wrote %s", string(resp))
		}
	}()

	return w.ResponseWriter.Write(resp)
}

func logIncomingRequest(ctx Context, reqMeta instrumenthttp.RequestMetadata, msg string) {
	tags := map[string]string{
		"http.response.status": strconv.Itoa(ctx.Writer().Status()),
		"http.response.size":   strconv.Itoa(ctx.Writer().Size()),
	}

	if len(reqMeta.BodyToLog) > 0 {
		tags["http.request.body"] = string(reqMeta.BodyToLog)
	}

	monitoring.FromContext(ctx.Request().Context()).
		With(tags).
		Infof(msg)
}

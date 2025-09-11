package lit

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/viebiz/lit/monitoring"
	"github.com/viebiz/lit/monitoring/instrumenthttp"
)

const (
	SkipLoggingResponseBodyKey = "skip_logging_response_body"
)

// rootMiddleware is a middleware function that handles tracing for incoming requests
// and recovers from any panics that may occur during request handling
func rootMiddleware(rootCtx context.Context) HandlerFunc {
	return func(c Context) error {
		// Start tracing for the incoming request
		ctx, reqMeta, endInstrumentation := instrumenthttp.StartIncomingRequest(monitoring.FromContext(rootCtx), c.Request(), c.FullPath())
		// Recovery logic when got panic
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
				c.Error(&HTTPError{
					Status: http.StatusInternalServerError,
					Code:   http.StatusText(http.StatusInternalServerError),
					Desc:   "Internal Server Error",
				})
				c.Abort() // Stop further middleware execution

				// End the instrumentation, marking the request with a 500 status code and the error.
				endInstrumentation(http.StatusInternalServerError, err)
			}
		}()

		// Update context, set instrument context and update response writer
		c.SetRequestContext(ctx)

		c.SetWriter(wrapWriter(ctx, c.Writer(), c.Get, reqMeta.Method, reqMeta.Path))

		// Continue handle request
		c.Next()

		// End instrumentation and log
		endInstrumentation(c.Writer().Status(), nil)
		logIncomingRequest(c, reqMeta, "http.incoming_request")

		return nil
	}
}

type responseRecorder struct {
	ResponseWriter

	ctx          context.Context
	keyExtractor func(key string) (any, bool)
	method, path string
}

func wrapWriter(ctx context.Context, w ResponseWriter, keyExtractor func(key string) (any, bool), method, path string) ResponseWriter {
	return &responseRecorder{ResponseWriter: w, ctx: ctx, keyExtractor: keyExtractor, method: method, path: path}
}

func (w *responseRecorder) Write(resp []byte) (n int, err error) {
	defer func() {
		if err != nil {
			monitoring.FromContext(w.ctx).Error(err, "[incoming_request] Failed to write response",
				monitoring.StringField("http.request.method", w.method),
				monitoring.StringField("url.path", w.path))
		} else {
			logFields := []monitoring.Field{
				monitoring.StringField("http.request.method", w.method),
				monitoring.StringField("url.path", w.path),
			}
			if _, exists := w.keyExtractor(SkipLoggingResponseBodyKey); !exists {
				logFields = append(logFields, monitoring.JSONField("http.response.body", resp))
			}

			monitoring.FromContext(w.ctx).Info("[incoming_request] Wrote response", logFields...)
		}
	}()

	return w.ResponseWriter.Write(resp)
}

func logIncomingRequest(ctx Context, reqMeta instrumenthttp.RequestMetadata, msg string) {
	monitoring.FromContext(ctx.Request().Context()).Info(msg,
		monitoring.StringField("http.request.method", reqMeta.Method),
		monitoring.StringField("url.path", reqMeta.Path),
		monitoring.StringField("url.query", reqMeta.Query),
		monitoring.JSONField("http.request.body", reqMeta.BodyToLog),
		monitoring.IntField("http.response.status_code", ctx.Writer().Status()),
		monitoring.IntField("http.response.body.size", ctx.Writer().Size()))
}

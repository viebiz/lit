package lit

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/pprof"

	"github.com/gin-contrib/cors"
	pkgerrors "github.com/pkg/errors"
)

// Handler defines a http.Handler that adds default readiness, liveness, profiling, test routes and cors policy
func Handler(
	rootCtx context.Context,
	corsConf CORSConfig,
	routerFunc func(Router),
	opts ...HandlerOption,
) http.Handler {
	r, hdl := NewRouter()

	cfg := handlerConfig{}
	for _, opt := range opts {
		opt(&cfg)
	}

	// Workaround to set up CORS middleware faster
	if rtr, ok := r.(router); ok {
		rtr.ginRouter.Use(cors.New(corsConf.cfg))
	}

	r.Handle(http.MethodGet, "/_/healthz", WrapF(LivenessHandlerFunc))

	// Usage on pprof refer to : https: //pkg.go.dev/net/http/pprof
	// by default profilingDisabled is false
	if !cfg.profilingDisabled {
		const prefix = "/_/profile"
		r.Handle(http.MethodGet, prefix+"/", WrapF(pprof.Index))
		r.Handle(http.MethodGet, prefix+"/cmdline", WrapF(pprof.Cmdline))
		r.Handle(http.MethodGet, prefix+"/profile", WrapF(pprof.Profile))
		r.Handle(http.MethodPost, prefix+"/symbol", WrapF(pprof.Symbol))
		r.Handle(http.MethodGet, prefix+"/symbol", WrapF(pprof.Symbol))
		r.Handle(http.MethodGet, prefix+"/trace", WrapF(pprof.Trace))
		r.Handle(http.MethodGet, prefix+"/allocs", WrapH(pprof.Handler("allocs")))
		r.Handle(http.MethodGet, prefix+"/block", WrapH(pprof.Handler("block")))
		r.Handle(http.MethodGet, prefix+"/goroutine", WrapH(pprof.Handler("goroutine")))
		r.Handle(http.MethodGet, prefix+"/heap", WrapH(pprof.Handler("heap")))
		r.Handle(http.MethodGet, prefix+"/mutex", WrapH(pprof.Handler("mutex")))
		r.Handle(http.MethodGet, prefix+"/threadcreate", WrapH(pprof.Handler("threadcreate")))
	}

	r.Use(rootMiddleware(rootCtx))

	// This route will help in testing integrations with monitoring system.
	r.Get("/_/test-monitor", func(c Context) error {
		c.JSON(http.StatusOK, map[string]string{
			"message": "Test Invoked",
		})
		return pkgerrors.WithStack(errors.New("test error"))
	})

	// Setup application router
	routerFunc(r)

	return hdl
}

// LivenessHandlerFunc is a simple handler func to be used fo health check operations
func LivenessHandlerFunc(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	_, _ = fmt.Fprintln(w, "ok")
}

// handlerConfig is configurations of the Handler
type handlerConfig struct {
	profilingDisabled bool
}

package lightning

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	pkgerrors "github.com/pkg/errors"
)

type Server struct {
	httpServer    *http.Server
	withTLS       bool
	certFile      string
	keyFile       string
	shutdownGrace time.Duration
}

// NewHttpServer creates new http server
func NewHttpServer(ctx context.Context, addr string, setupRouteFunc func(Router), opts ...func(*Server)) Server {
	// Create new router
	r, hdl := NewRouter()
	r.Use(rootMiddleware(ctx))

	// Setup routes
	setupRouteFunc(r)

	// Setup server
	srv := Server{
		httpServer: &http.Server{
			Addr:    addr,
			Handler: hdl,
		},
	}

	// Configures server
	for _, opt := range opts {
		opt(&srv)
	}

	return srv
}

func (srv *Server) Run() error {
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	return srv.start(ctx)
}

func (srv *Server) start(ctx context.Context) error {
	startupErr := make(chan error)

	// Start server
	go func() {
		fmt.Printf("web server started; listening at %s\n", srv.httpServer.Addr)
		defer fmt.Println("web server shutdown")

		var err error
		if srv.withTLS {
			err = srv.httpServer.ListenAndServeTLS(srv.certFile, srv.keyFile)
		} else {
			err = srv.httpServer.ListenAndServe()
		}

		if err != nil {
			startupErr <- err
		}
	}()

	// Blocking main and waiting for shutdown.
	select {
	case err := <-startupErr:
		// ListenAndServe will always return a non-nil error
		if !errors.Is(err, http.ErrServerClosed) {
			return pkgerrors.Wrap(err, "http server stopped")
		}
		return nil
	case <-ctx.Done():
		return srv.stop()
	}
}

func (srv *Server) stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), srv.shutdownGrace)
	defer cancel()

	fmt.Printf("attempting to shutdown gracefully\n")
	defer fmt.Println("server shutdown successfully")

	if err := srv.httpServer.Shutdown(ctx); err != nil {
		fmt.Printf("failed to shutdown gracefully: %v, force shutdown\n", err)

		if err = srv.httpServer.Close(); err != nil {
			return pkgerrors.Wrap(err, "force shutdown")
		}
	}

	return nil
}

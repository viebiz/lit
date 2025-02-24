package lit

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHttpServer(t *testing.T) {
	tcs := []struct {
		givenAddr         string
		givenOpts         []ServerOption
		wantReadTimeout   time.Duration
		wantWriteTimeout  time.Duration
		wantShutdownGrace time.Duration
		wantPort          string
	}{
		{
			givenAddr:         ":3000",
			wantReadTimeout:   time.Minute,
			wantWriteTimeout:  time.Minute,
			wantShutdownGrace: 0,
			wantPort:          ":3000",
		},
		{
			givenAddr:         ":5000",
			givenOpts:         []ServerOption{ServerReadTimeout(time.Hour), ServerShutdownGrace(time.Second)},
			wantReadTimeout:   time.Hour,
			wantWriteTimeout:  time.Minute,
			wantShutdownGrace: time.Second,
			wantPort:          ":5000",
		},
		{
			givenAddr:         ":1604",
			givenOpts:         []ServerOption{ServerWriteTimeout(time.Hour)},
			wantReadTimeout:   time.Minute,
			wantWriteTimeout:  time.Hour,
			wantShutdownGrace: 0,
			wantPort:          ":1604",
		},
	}
	for i, tc := range tcs {
		t.Run(fmt.Sprintf("scenario: %d", i), func(t *testing.T) {
			// Given:

			// When:
			s := NewHttpServer(tc.givenAddr, emptyHandler{}, tc.givenOpts...)

			// Then:
			require.Equal(t, tc.wantReadTimeout, s.httpServer.ReadTimeout)
			require.Equal(t, tc.wantWriteTimeout, s.httpServer.WriteTimeout)
			require.Equal(t, tc.wantShutdownGrace, s.shutdownGrace)
			require.Equal(t, tc.wantPort, s.httpServer.Addr)
		})
	}
}

func TestRunWithContext(t *testing.T) {
	server := NewHttpServer("127.0.0.1:0", emptyHandler{}, ServerShutdownGrace(2*time.Second))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	go func() {
		_ = server.RunWithContext(ctx)
	}()

	time.Sleep(100 * time.Millisecond) // Ensure server starts

	cancel() // Trigger shutdown

	err := server.RunWithContext(ctx)
	assert.NoError(t, err)
}

func TestRun(t *testing.T) {
	const addr = "127.0.0.1:0"
	server := NewHttpServer(addr, emptyHandler{}, ServerShutdownGrace(2*time.Second))

	// Simulate SIGINT
	go func() {
		time.Sleep(100 * time.Millisecond)
		p, _ := os.FindProcess(os.Getpid())
		_ = p.Signal(syscall.SIGINT)
	}()

	err := server.Run()
	assert.NoError(t, err)
}

// Test stop ensures server shuts down properly
func TestStop(t *testing.T) {
	const addr = "127.0.0.1:0"
	server := NewHttpServer(addr, emptyHandler{}, ServerShutdownGrace(time.Second))

	err := server.stop()
	assert.NoError(t, err)
}

type emptyHandler struct{}

func (emptyHandler) ServeHTTP(http.ResponseWriter, *http.Request) {}

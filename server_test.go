package lit

import (
	"context"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewHttpServer(t *testing.T) {
	server := NewHttpServer(context.Background(), "127.0.0.1:8080", func(r Router) {
		// Mock route setup
	}, func(s *Server) {
		s.withTLS = false
	})

	assert.NotNil(t, server.httpServer)
	assert.Equal(t, "127.0.0.1:8080", server.httpServer.Addr)
	assert.Equal(t, defaultServerReadTimeout, server.httpServer.ReadTimeout)
	assert.Equal(t, defaultServerWriteTimeout, server.httpServer.WriteTimeout)
}

func TestRunWithContext(t *testing.T) {
	server := NewHttpServer(context.Background(), "127.0.0.1:0", func(r Router) {}, ServerShutdownGrace(2*time.Second))

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
	server := NewHttpServer(context.Background(), "127.0.0.1:0", func(r Router) {}, func(s *Server) {
		s.shutdownGrace = 2 * time.Second
	})

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
	server := NewHttpServer(context.Background(), "127.0.0.1:0", func(r Router) {}, func(s *Server) {
		s.shutdownGrace = time.Second
	})

	err := server.stop()
	assert.NoError(t, err)
}

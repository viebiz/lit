package lit

import (
	"time"
)

// ServerOption represents option for creates HTTP server
type ServerOption func(*Server)

func ServerShutdownGrace(duration time.Duration) ServerOption {
	return func(s *Server) {
		s.shutdownGrace = duration
	}
}

// ServerReadTimeout overrides the server's default account timeout with the given one.
func ServerReadTimeout(duration time.Duration) ServerOption {
	return func(s *Server) {
		s.httpServer.ReadTimeout = duration
	}
}

// ServerWriteTimeout overrides the server's default write timeout with the given one.
func ServerWriteTimeout(duration time.Duration) ServerOption {
	return func(s *Server) {
		s.httpServer.WriteTimeout = duration
	}
}

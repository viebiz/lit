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

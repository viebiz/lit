package instrumentpg

import (
	"github.com/viebiz/lit/postgres"
)

// WithInstrumentation wraps a ContextExecutor with instrumentation.
func WithInstrumentation(pool postgres.ContextExecutor) postgres.ContextExecutor {
	return instrumentedDB{ContextExecutor: pool}
}

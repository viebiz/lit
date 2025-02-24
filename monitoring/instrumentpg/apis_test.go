package instrumentpg

import (
	"testing"

	"github.com/viebiz/lit/postgres"
)

func TestWithInstrumentation(t *testing.T) {
	var pool postgres.ContextExecutor
	pool = WithInstrumentation(pool)
}

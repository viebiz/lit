package instrumentpg

import "testing"

func TestWithInstrumentation(t *testing.T) {
	t.Parallel()
	_ = WithInstrumentation(nil)
}

func TestWithInstrumentationTx(t *testing.T) {
	t.Parallel()
	_ = WithInstrumentationTx(nil)
}

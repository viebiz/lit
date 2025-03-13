package i18n

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetInContext(t *testing.T) {
	b := &localizer{}
	parentCtx := context.Background()
	ctx := SetInContext(parentCtx, b)

	// Retrieve using FromContext.
	retrieved := FromContext(ctx)
	require.Equal(t, b, retrieved)
}

func TestFromContext(t *testing.T) {
	lc := &localizer{}
	ctx := context.WithValue(context.Background(), contextKey{}, lc)

	retrieved := FromContext(ctx)
	require.Equal(t, lc, retrieved)
}

func TestFromContext_Noop(t *testing.T) {
	ctx := context.Background()
	retrieved := FromContext(ctx)

	require.NotNil(t, retrieved)
	require.Equal(t, retrieved, localizer{})
}

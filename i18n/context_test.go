package i18n

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetInContext(t *testing.T) {
	b := &bundle{}
	parentCtx := context.Background()
	ctx := SetInContext(parentCtx, b)

	// Retrieve using FromContext.
	retrieved := FromContext(ctx)
	require.Equal(t, b, retrieved)
}

func TestFromContext(t *testing.T) {
	b := &bundle{}
	ctx := context.WithValue(context.Background(), contextKey{}, b)

	retrieved := FromContext(ctx)
	require.Equal(t, b, retrieved)
}

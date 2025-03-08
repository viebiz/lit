package monitoring

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/viebiz/lit/monitoring/tracing/mocktracer"
	"go.uber.org/zap"
)

func Test_SetInContext(t *testing.T) {
	ctx := context.Background()
	m := &Monitor{logger: zap.NewNop()}

	ctx = SetInContext(ctx, m)
}

func Test_FromContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), monitorContextKey{}, &Monitor{logger: zap.NewNop()})

	m := FromContext(ctx)
	require.NotNil(t, m)
}

func Test_NewContext(t *testing.T) {
	mt := mocktracer.Start()
	defer mt.Stop()
	
	// Parent context
	type ctxKey struct{}
	ctx := context.WithValue(context.Background(), ctxKey{}, 1)

	ctx, span := tracer.Start(ctx, "test")
	defer span.End()

	// Child context
	ctx = NewContext(ctx)

	// Value of parent context should not exist
	require.Nil(t, ctx.Value(ctxKey{}))
}

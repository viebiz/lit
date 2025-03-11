package i18n

import (
	"context"
)

type contextKey struct{}

func FromContext(ctx context.Context) Bundle {
	return ctx.Value(contextKey{}).(Bundle)
}

func SetInContext(parentCtx context.Context, localeManager Bundle) context.Context {
	return context.WithValue(parentCtx, contextKey{}, localeManager)
}

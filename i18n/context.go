package i18n

import (
	"context"
)

type contextKey struct{}

func FromContext(ctx context.Context) Localizable {
	if m, ok := ctx.Value(contextKey{}).(Localizable); ok && m != nil {
		return m
	}

	return localizer{} // Noop localize
}

func SetInContext(parentCtx context.Context, lc Localizable) context.Context {
	return context.WithValue(parentCtx, contextKey{}, lc)
}

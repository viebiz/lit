package guard

import (
	"github.com/viebiz/lit"
	"github.com/viebiz/lit/iam"
	"github.com/viebiz/lit/monitoring"
)

func (guard AuthGuard) RequiredM2MScopeMiddleware(scopes ...string) lit.HandlerFunc {
	return func(c lit.Context) {
		ctx := c.Request().Context()

		// 1. Get M2M profile from request context
		profile := iam.GetM2MProfileFromContext(ctx)
		if profile.ID() == "" {
			monitoring.FromContext(ctx).Errorf(errUserProfileNotInCtx, "Missing M2M profile in context")
			c.AbortWithError(errForbidden)
			return
		}

		// 2. Check if profile has any required scopes
		if !profile.HasAnyScope(scopes...) {
			c.AbortWithError(errForbidden)
			return
		}

		// 3. Continue with the next handler
		c.Next()
	}
}

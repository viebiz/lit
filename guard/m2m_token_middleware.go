package guard

import (
	"github.com/viebiz/lit"
	"github.com/viebiz/lit/iam"
	"github.com/viebiz/lit/monitoring"
)

const (
	m2mIDKey = "m2m_id"
)

func (guard AuthGuard) AuthenticateM2MMiddleware() lit.HandlerFunc {
	return func(c lit.Context) {
		// 1. Get access token from request header
		tokenStr := getTokenString(c.Request())
		if tokenStr == "" {
			c.AbortWithError(errMissingAccessToken)
			return
		}

		// 2. Validate access token
		tk, err := guard.validator.Validate(getTokenString(c.Request()))
		if err != nil {
			responseErr(c, err)
			return
		}

		// 3. Extract M2M profile from token claims
		profile, err := iam.ExtractM2MProfileFromClaims(tk.Claims)
		if err != nil {
			responseErr(c, err)
			return
		}

		// 4. Inject user information to request context
		ctx := c.Request().Context()
		ctx = iam.SetM2MProfileInContext(ctx, profile)
		ctx = monitoring.InjectField(ctx, m2mIDKey, profile.ID())
		c.SetRequestContext(ctx)

		// 5. Continue handle request
		c.Next()
	}
}

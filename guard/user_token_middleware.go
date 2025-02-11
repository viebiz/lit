package guard

import (
	"net/http"
	"strings"

	"github.com/viebiz/lit"
	"github.com/viebiz/lit/iam"
	"github.com/viebiz/lit/monitoring"
)

const (
	headerAuthorization       = "Authorization"
	authorizationBearerPrefix = "Bearer"
	userIDKey                 = "user_id"
	roleKey                   = "roles"
)

func (guard AuthGuard) AuthenticateUserMiddleware() lit.HandlerFunc {
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

		// 3. Extract user profile from token claims
		profile, err := iam.ExtractUserProfileFromClaims(tk.Claims)
		if err != nil {
			responseErr(c, err)
			return
		}

		// 4. Inject user information to request context
		ctx := c.Request().Context()
		ctx = iam.SetUserProfileInContext(ctx, profile)
		ctx = monitoring.InjectFields(ctx, map[string]string{
			userIDKey: profile.ID(),
			roleKey:   profile.GetRoleString(),
		})
		c.SetRequestContext(ctx)

		// 5. Continue handle request
		c.Next()
	}
}

func getTokenString(r *http.Request) string {
	authHeaderParts := strings.Split(r.Header.Get(headerAuthorization), " ")
	if len(authHeaderParts) != 2 || authHeaderParts[0] != authorizationBearerPrefix {
		return ""
	}

	return authHeaderParts[1]
}

func responseErr(c lit.Context, err error) {
	switch err.Error() {
	case iam.ErrMissingRequiredClaim.Error(),
		iam.ErrTokenExpired.Error(),
		iam.ErrInvalidToken.Error():
		c.AbortWithError(unauthorizedErr(err))

	case iam.ErrActionIsNotAllowed.Error():
		c.AbortWithError(errForbidden)

	default:
		c.AbortWithError(lit.ErrInternalServerError)
		monitoring.FromContext(c.Request().Context()).Errorf(err, "Got unexpected error")
	}
}

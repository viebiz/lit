package guard

import (
	"errors"
	"net/http"

	"github.com/viebiz/lit"
)

const (
	unAuthorizedKey = "unauthorized"
	forbiddenKey    = "forbidden"
)

var (
	errUserProfileNotInCtx = errors.New("user profile not in context")
	errM2MProfileNotInCtx  = errors.New("m2m profile not in context")
	errMissingAccessToken  = lightning.HttpError{Status: http.StatusUnauthorized, Code: unAuthorizedKey, Description: "Access token is required"}
	errForbidden           = lightning.HttpError{Status: http.StatusForbidden, Code: forbiddenKey, Description: "Permission denied"}
)

func unauthorizedErr(err error) lightning.HttpError {
	return lightning.HttpError{Status: http.StatusUnauthorized, Code: unAuthorizedKey, Description: err.Error()}
}

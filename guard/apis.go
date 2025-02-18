package guard

import (
	"github.com/viebiz/lit/iam"
)

func New(validator iam.Validator, enforcer iam.Enforcer) AuthGuard {
	return AuthGuard{
		validator: validator,
		enforcer:  enforcer,
	}
}

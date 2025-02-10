package guard

import (
	"github.com/viebiz/lit/iam"
)

type AuthGuard struct {
	validator iam.Validator
	enforcer  iam.Enforcer
}

func New(validator iam.Validator, enforcer iam.Enforcer) AuthGuard {
	return AuthGuard{
		validator: validator,
		enforcer:  enforcer,
	}
}

package guard

import (
	"github.com/viebiz/lit/iam"
)

type AuthGuard struct {
	validator iam.Validator
	enforcer  iam.Enforcer
}

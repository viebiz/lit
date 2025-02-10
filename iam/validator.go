package iam

import (
	"time"

	"github.com/viebiz/lit/jwt"
)

var (
	timeNowFunc = func() time.Time {
		return time.Now().UTC()
	}
)

type Validator interface {
	Validate(tokenString string) (jwt.Token[Claims], error)
}

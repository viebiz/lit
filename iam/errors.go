package iam

import (
	"errors"
)

var (
	ErrMissingRequiredClaim = errors.New("missing required claim")

	ErrInvalidToken = errors.New("invalid token")

	ErrTokenExpired = errors.New("token expired")

	ErrActionIsNotAllowed = errors.New("action is not allowed")
)

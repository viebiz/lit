package jwt

import (
	"errors"
)

var (
	ErrInvalidKeyType = errors.New("invalid key type")

	ErrHashUnavailable = errors.New("unavailable hash function")

	ErrTokenMalformed = errors.New("malformed token")

	ErrInvalidToken = errors.New("invalid token")

	ErrSigningMethodNotSupported = errors.New("signing method not supported")

	ErrInvalidSignature = errors.New("invalid signature")

	ErrMissingRequiredClaim = errors.New("missing required claim in token")

	ErrTokenExpired = errors.New("token is expired")

	ErrTokenUsedBeforeIssued = errors.New("token is used before issued")

	ErrTokenNotValidYet = errors.New("token is not valid yet")
)

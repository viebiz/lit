package jwt

import (
	"time"
)

// RegisteredClaims represents standard JWT claims
// More info: https://datatracker.ietf.org/doc/html/rfc7519#section-4.1
type RegisteredClaims struct {
	Issuer string `json:"iss,omitempty"`

	Subject string `json:"sub,omitempty"`

	Audience ClaimStrings `json:"aud,omitempty"`

	IssuedAt *int64 `json:"iat,omitempty"`

	ExpiresAt *int64 `json:"exp,omitempty"`

	NotBefore *int64 `json:"nbf,omitempty"`

	JTI string `json:"jti,omitempty"`

	ClientID string `json:"client_id,omitempty"`
}

// Valid validates time based claims "exp, iat, nbf"
// if any of the above claims are not in the token, it will still
// be considered a valid claim.
func (c RegisteredClaims) Valid() error {
	now := timeNowFunc()

	// 1. Verify the `exp` claim, it's required claims
	if c.ExpiresAt == nil {
		return ErrMissingRequiredClaim
	}

	exp := time.Unix(*c.ExpiresAt, 0)
	if now.After(exp) {
		return ErrTokenExpired
	}

	// 2. Verify the `iat` claim
	// If it's empty, considered as a valid claim
	if c.IssuedAt != nil {
		issuedAt := time.Unix(*c.IssuedAt, 0)
		if now.Before(issuedAt) {
			return ErrTokenUsedBeforeIssued

		}
	}

	// 3. Verify the `nbf` claim
	// If it's empty, considered as a valid claim
	if c.NotBefore != nil {
		notBefore := time.Unix(*c.NotBefore, 0)
		if now.Before(notBefore) {
			return ErrTokenNotValidYet
		}
	}

	return nil
}

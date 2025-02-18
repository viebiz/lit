package jwt

import (
	"crypto"
	"io"
)

// Signer represents an interface for creating digital signatures
type Signer interface {
	Sign(rand io.Reader, digest []byte, opts crypto.SignerOpts) (signature []byte, err error)
}

// VerifyKey represents a key for verify token
type VerifyKey interface{}

// SigningMethod can be used add new methods for signing or verifying tokens. It
// takes a decoded signature as an input in the Verify function and produces a
// signature in Sign. The signature is then usually base64 encoded as part of a
// JWT.
type SigningMethod interface {
	// Verify returns nil if signature is valid
	Verify(signingString []byte, sig []byte, key VerifyKey) error

	// Sign returns signature or error
	Sign(signingString []byte, key Signer) ([]byte, error)

	// Alg returns the alg identifier for this method (example: 'RS256')
	Alg() string
}

// Claims represent any form of a JWT Claims
type Claims interface {
	Valid() error
}

package jwt

import (
	"crypto"
)

func NewHS256() HMAC {
	return HMAC{
		Name: SigningMethodNameHS256,
		Hash: crypto.SHA256,
	}
}

func NewHS384() HMAC {
	return HMAC{
		Name: SigningMethodNameHS384,
		Hash: crypto.SHA384,
	}
}

func NewHS512() HMAC {
	return HMAC{
		Name: SigningMethodNameHS512,
		Hash: crypto.SHA512,
	}
}

// NewRS256 creates a new RS256 signing method struct
func NewRS256() RSA {
	return RSA{
		Name: SigningMethodNameRS256,
		Hash: crypto.SHA256,
	}
}

// NewRS384 creates a new RS256 signing method struct
func NewRS384() RSA {
	return RSA{
		Name: SigningMethodNameRS384,
		Hash: crypto.SHA384,
	}
}

// NewRS512 creates a new RS256 signing method struct
func NewRS512() RSA {
	return RSA{
		Name: SigningMethodNameRS512,
		Hash: crypto.SHA512,
	}
}

// NewParser creates a new Parser with the default signing methods and validator.
func NewParser[T Claims](opts ...ParserOptions) Parser[T] {
	p := NewDefaultParser[T]()
	for _, opt := range opts {
		opt((*Parser[Claims])(&p))
	}

	return p
}

func NewDefaultParser[T Claims]() Parser[T] {
	return Parser[T]{
		signingMethods: map[string]SigningMethod{
			SigningMethodNameRS256: NewRS256(),
			SigningMethodNameRS384: NewRS384(),
			SigningMethodNameRS512: NewRS512(),
			SigningMethodNameHS256: NewHS256(),
			SigningMethodNameHS384: NewHS384(),
			SigningMethodNameHS512: NewHS512(),
		},
	}
}

// NewToken creates a new Token with the specified signing method and claims
func NewToken[T Claims](method SigningMethod, claims T) Token[T] {
	tk := Token[T]{
		Header: map[string]string{
			"typ": tokenHeaderType,
			"alg": method.Alg(),
		},
		Claims:        claims,
		signingMethod: method,
	}

	return tk
}

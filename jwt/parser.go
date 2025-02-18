package jwt

import (
	"crypto"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

// Parser is a generic struct for parsing and validating JWT strings.
// T represents a specific Claims type to ensure precise handling of predefined claims.
type Parser[T Claims] struct {
	signingMethods map[string]SigningMethod // Supported signing method
}

// Parse decodes the provided JWT string into a Token, verifies its signature using the provided public key.
// The public key can be determined dynamically based on the `kid` (Key ID) in the token header.
func (p Parser[T]) Parse(tokenString string, getKeyFunc func(string) (crypto.PublicKey, error)) (Token[T], error) {
	// 1. Parse JWT string to Token
	token, signingString, err := p.parseToken(tokenString)
	if err != nil {
		return Token[T]{}, err
	}

	// 2. Lookup signing method
	// 2.1. Get signing algorithm in token header
	alg, exists := token.Header["alg"]
	if !exists {
		return Token[T]{}, ErrInvalidToken
	}

	// 2.2. get signing method by name
	signingMethod, ok := p.getSigningMethod(alg)
	if !ok {
		return Token[T]{}, ErrSigningMethodNotSupported
	}

	// 3. Verify signature by signing method
	// 3.1. Get Public/Private key for verify token,
	// Optionally, use the `kid` (Key ID) to determine which key to use for verification.
	key, err := getKeyFunc(token.Header["kid"])
	if err != nil {
		return Token[T]{}, err
	}

	// 3.2. Verify token signature
	if err := signingMethod.Verify(signingString, token.Signature, key); err != nil {
		return Token[T]{}, ErrInvalidSignature
	}

	// 3.3. Validate token claims
	if err := token.Claims.Valid(); err != nil {
		return Token[T]{}, err
	}

	return token, nil
}

func (p Parser[T]) parseToken(tokenString string) (Token[T], []byte, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return Token[T]{}, nil, ErrTokenMalformed
	}

	// 1.1. Decode header segment
	headerBytes, err := decodeSegment([]byte(parts[0]))
	if err != nil {
		return Token[T]{}, nil, err
	}

	// 1.2. Decode JSON header to map
	var headerMap map[string]string
	if err := json.Unmarshal(headerBytes, &headerMap); err != nil {
		return Token[T]{}, nil, err
	}

	// 2.1. Decode payload segment
	payloadBytes, err := decodeSegment([]byte(parts[1]))
	if err != nil {
		return Token[T]{}, nil, err
	}

	// 2.2. Decode JSON Claims
	var claims T
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return Token[T]{}, nil, err
	}

	// 3.3. Decode signature
	sigBytes, err := decodeSegment([]byte(parts[2]))
	if err != nil {
		return Token[T]{}, nil, err
	}

	signingString := fmt.Sprintf("%s.%s", parts[0], parts[1])

	return Token[T]{
		Header:    headerMap,
		Claims:    claims,
		Signature: sigBytes,
	}, []byte(signingString), nil
}

func (p Parser[T]) getSigningMethod(alg string) (SigningMethod, bool) {
	method, ok := p.signingMethods[alg]
	return method, ok
}

func decodeSegment(b []byte) ([]byte, error) {
	decoded := make([]byte, base64.RawURLEncoding.DecodedLen(len(b)))
	if _, err := base64.RawURLEncoding.Decode(decoded, b); err != nil {
		return nil, err
	}

	return decoded, nil
}

package jwt

import (
	"bytes"
	"crypto"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

const (
	tokenHeaderType = "JWT"
)

// Token represents a JWT Token
type Token[T Claims] struct {
	Header    map[string]string
	Claims    T
	Signature []byte

	// signingMethod specifies the algorithm used to sign the token.
	signingMethod SigningMethod
}

// SignedString creates and returns a complete, signed JWT. The token is signed
// using the SigningMethod specified in the token
func (tk Token[T]) SignedString(key crypto.Signer) (string, error) {
	msg, err := tk.signingString()
	if err != nil {
		return "", err
	}

	// Sign the concatenated message
	sig, err := tk.signingMethod.Sign(msg, key)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s.%s", msg, encodeSegment(sig)), nil
}

func (tk Token[T]) signingString() ([]byte, error) {
	// Marshal the header map to JSON format
	headerBytes, err := json.Marshal(tk.Header)
	if err != nil {
		return nil, err
	}

	// Marshal the payload (claims) to JSON format
	payloadBytes, err := json.Marshal(tk.Claims)
	if err != nil {
		return nil, err
	}

	// Base64 encode the header and payload
	headerBase64, payloadBase64 := encodeSegment(headerBytes), encodeSegment(payloadBytes)

	return bytes.Join([][]byte{headerBase64, payloadBase64}, []byte(".")), nil
}

func encodeSegment(b []byte) []byte {
	buf := make([]byte, base64.RawURLEncoding.EncodedLen(len(b)))
	base64.RawURLEncoding.Encode(buf, b)

	return buf
}

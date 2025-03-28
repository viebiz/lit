package aes

import (
	"crypto/rand"
	"io"
)

// GenerateNonce generates nonce 12 bytes
func (i *impl) GenerateNonce() ([]byte, error) {
	nonce := make([]byte, 12) // nonce should be 12 bytes

	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return nonce, nil
}

// GenerateIV generates Initialization Vector 16 bytes
func (i *impl) GenerateIV() ([]byte, error) {
	nonce := make([]byte, 16) // nonce should be 12 bytes

	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return nonce, nil
}

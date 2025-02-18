package jwt

import (
	"crypto"
	"crypto/hmac"
	"io"
)

const (
	SigningMethodNameHS256 string = "HS256"
	SigningMethodNameHS384 string = "HS384"
	SigningMethodNameHS512 string = "HS512"
)

type HMAC struct {
	Name string
	Hash crypto.Hash
}

func (sm HMAC) Verify(signingString []byte, sig []byte, key VerifyKey) error {
	keyBytes, ok := key.(HMACPrivateKey)
	if !ok {
		return ErrInvalidKeyType
	}

	if !sm.Hash.Available() {
		return ErrHashUnavailable
	}

	hasher := hmac.New(sm.Hash.New, keyBytes)
	hasher.Write(signingString)
	if !hmac.Equal(hasher.Sum(nil), sig) {
		return ErrInvalidSignature
	}

	return nil
}

func (sm HMAC) Sign(signingString []byte, key Signer) ([]byte, error) {
	pkey, ok := key.(HMACPrivateKey)
	if !ok {
		return nil, ErrInvalidKeyType
	}

	if !sm.Hash.Available() {
		return nil, ErrHashUnavailable
	}

	return pkey.Sign(nil, signingString, sm.Hash)
}

func (sm HMAC) Alg() string {
	return sm.Name
}

// HMACPrivateKey represents private key for HMAC signing method
type HMACPrivateKey []byte

func (h HMACPrivateKey) Sign(rand io.Reader, signingString []byte, opts crypto.SignerOpts) ([]byte, error) {
	hasher := hmac.New(opts.HashFunc().New, h)
	hasher.Write(signingString)

	return hasher.Sum(nil), nil
}

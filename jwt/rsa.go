package jwt

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
)

const (
	SigningMethodNameRS256 string = "RS256"
	SigningMethodNameRS384 string = "RS384"
	SigningMethodNameRS512 string = "RS512"
)

// RSA implements the RSA family of signing methods
type RSA struct {
	Name string
	Hash crypto.Hash
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

// Sign implements token signing for the SigningMethod, that take Signer (*rsa.PrivateKey) for sign a token
func (sm RSA) Sign(signingString []byte, key Signer) ([]byte, error) {
	if _, ok := key.(*rsa.PrivateKey); !ok {
		return nil, ErrInvalidKeyType
	}

	if !sm.Hash.Available() {
		return nil, ErrHashUnavailable
	}

	hash := sm.Hash.New()
	hash.Write(signingString)

	return key.Sign(rand.Reader, hash.Sum(nil), sm.Hash)
}

// Verify verifies the signingString with signature by provided VerifyKey, that can consider as *rsa.PublicKey
func (sm RSA) Verify(signingString []byte, sig []byte, key VerifyKey) error {
	publicKey, ok := key.(*rsa.PublicKey)
	if !ok {
		return ErrInvalidKeyType
	}

	if !sm.Hash.Available() {
		return ErrHashUnavailable
	}

	hash := sm.Hash.New()
	hash.Write(signingString)

	return rsa.VerifyPKCS1v15(publicKey, sm.Hash, hash.Sum(nil), sig)
}

func (sm RSA) Alg() string {
	return sm.Name
}

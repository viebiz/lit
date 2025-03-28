package rsa

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"

	pkgerrors "github.com/pkg/errors"
)

var (
	newPrivateKeyFn = newPrivateKey
	newPublicKeyFn  = newPublicKey
)

func newPrivateKey(keyPath string) (*rsa.PrivateKey, error) {
	// Read file
	keyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, pkgerrors.WithStack(err)
	}

	// Decode
	block, _ := pem.Decode(keyBytes)
	parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	key, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		return nil, ErrNotRSAPrivateKey
	}

	return key, nil
}

func newPublicKey(keyPath string) (*rsa.PublicKey, error) {
	// Read file
	keyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, pkgerrors.WithStack(err)
	}

	// Decode
	block, _ := pem.Decode(keyBytes)
	parsedKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, pkgerrors.WithStack(err)
	}

	key, ok := parsedKey.(*rsa.PublicKey)
	if !ok {
		return nil, ErrNotRSAPublicKey
	}

	return key, nil
}

package rsa

import "errors"

var (
	ErrNotRSAPrivateKey = errors.New("not RSA private key")
	ErrNotRSAPublicKey  = errors.New("not RSA public key")
)

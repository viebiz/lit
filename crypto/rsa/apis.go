package rsa

import (
	"crypto/rsa"
)

func NewPrivateKey(keyPath string) (*rsa.PrivateKey, error) {
	return newPrivateKeyFn(keyPath)
}

func NewPublicKey(keyPath string) (*rsa.PublicKey, error) {
	return newPublicKeyFn(keyPath)
}

package rsa

import (
	"crypto/rsa"
)

type impl struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func New(privatePath, publicPath string) (Cipher, error) {
	privateKey, err := newPrivateKeyFn(privatePath)
	if err != nil {
		return nil, err
	}
	publicKey, err := newPublicKeyFn(publicPath)
	if err != nil {
		return nil, err
	}
	return &impl{privateKey: privateKey, publicKey: publicKey}, nil
}

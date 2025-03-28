package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
)

func (i *impl) Decrypt(encryptedText string) (string, error) {
	// Base64 decode
	encryptedTextBytes, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}

	// Decrypt encrypted text
	decryptedTextBytes, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, i.privateKey, encryptedTextBytes, nil)
	if err != nil {
		return "", err
	}

	return string(decryptedTextBytes), nil
}

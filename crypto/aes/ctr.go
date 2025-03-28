package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"

	pkgerrors "github.com/pkg/errors"
)

func (i *impl) CTREncrypt(plainText string, nonce []byte) (string, error) {
	block, err := aes.NewCipher(i.key)
	if err != nil {
		return "", pkgerrors.WithStack(err)
	}

	stream := cipher.NewCTR(block, nonce)

	plainTextBytes := []byte(plainText)
	ciphertext := make([]byte, len(plainTextBytes))
	stream.XORKeyStream(ciphertext, plainTextBytes)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (i *impl) CTRDecrypt(encryptedText string, nonce []byte) (string, error) {
	encryptedTextBytes, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", pkgerrors.WithStack(err)
	}

	block, err := aes.NewCipher(i.key)
	if err != nil {
		return "", pkgerrors.WithStack(err)
	}

	stream := cipher.NewCTR(block, nonce)

	decryptedTextBytes := make([]byte, len(encryptedTextBytes))
	stream.XORKeyStream(decryptedTextBytes, encryptedTextBytes)

	return string(decryptedTextBytes), nil
}

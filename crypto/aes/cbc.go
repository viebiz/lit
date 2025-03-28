package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	pkgerrors "github.com/pkg/errors"
)

func (i *impl) CBCEncrypt(plainText string, iv []byte) (string, error) {
	block, err := aes.NewCipher(i.key)
	if err != nil {
		return "", pkgerrors.WithStack(err)
	}
	plainTextBytes := applyPKCS7Padding([]byte(plainText))
	ciphertext := make([]byte, len(plainTextBytes))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, plainTextBytes)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (i *impl) CBCDecrypt(encryptedText string, iv []byte) (string, error) {
	encryptedTextBytes, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", pkgerrors.WithStack(err)
	}

	block, err := aes.NewCipher(i.key)
	if err != nil {
		return "", pkgerrors.WithStack(err)
	}

	decryptedText := make([]byte, len(encryptedTextBytes))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(decryptedText, encryptedTextBytes)

	plainTextBytes, err := stripPKCS7Padding(decryptedText)
	if err != nil {
		return "", pkgerrors.WithStack(err)
	}

	return string(plainTextBytes), nil
}

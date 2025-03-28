package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	pkgerrors "github.com/pkg/errors"
)

func (i *impl) GCMEncrypt(plainText string, nonce []byte) (string, error) {
	block, err := aes.NewCipher(i.key)
	if err != nil {
		return "", pkgerrors.WithStack(err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", pkgerrors.WithStack(err)
	}

	ciphertext := aesGCM.Seal(nil, nonce, []byte(plainText), nil)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (i *impl) GCMDecrypt(encryptedText string, nonce []byte) (string, error) {
	encryptedTextBytes, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", pkgerrors.WithStack(err)
	}

	block, err := aes.NewCipher(i.key)
	if err != nil {
		return "", pkgerrors.WithStack(err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", pkgerrors.WithStack(err)
	}

	decryptedText, err := aesGCM.Open(nil, nonce, encryptedTextBytes, nil)
	if err != nil {
		return "", pkgerrors.WithStack(err)
	}

	return string(decryptedText), nil
}

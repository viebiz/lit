package aes

import (
	"bytes"
	"crypto/aes"
)

// applyPKCS7Padding add 16 bytes
func applyPKCS7Padding(data []byte) []byte {
	padding := aes.BlockSize - len(data)%aes.BlockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// stripPKCS7Padding remove 16 bytes
func stripPKCS7Padding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 || length%aes.BlockSize != 0 {
		return nil, ErrInvalidData
	}

	padding := int(data[length-1])
	if padding > aes.BlockSize || padding == 0 {
		return nil, ErrInvalidPaddingValue
	}

	for i := 0; i < padding; i++ {
		if data[length-1-i] != byte(padding) {
			return nil, ErrIncorrectPadding
		}
	}

	return data[:length-padding], nil
}

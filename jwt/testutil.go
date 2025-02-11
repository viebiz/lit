package jwt

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/stretchr/testify/require"
)

func parsePrivateKeyFromPEM[T crypto.PrivateKey](b []byte) (T, error) {
	bl, _ := pem.Decode(b)
	if bl == nil {
		return *new(T), fmt.Errorf("failed to decode PEM block")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(bl.Bytes)
	if err != nil {
		return *new(T), err
	}

	rs, ok := privateKey.(T)
	if !ok {
		return *new(T), fmt.Errorf("key type not match")
	}

	return rs, nil
}

func readKeyForTest[T crypto.Signer](t require.TestingT, path string) T {
	b, err := os.ReadFile(path)
	require.NoError(t, err)

	key, err := parsePrivateKeyFromPEM[T](b)
	require.NoError(t, err)

	return key
}

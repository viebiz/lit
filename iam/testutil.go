package iam

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"math/big"
	"net/http"
	"os"
	"strings"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func readRSAPrivateKey(t require.TestingT, path string) *rsa.PrivateKey {
	b, err := os.ReadFile(path)
	require.NoError(t, err)

	block, _ := pem.Decode(b)
	require.True(t, block != nil)

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	require.NoError(t, err)

	rsaPrivateKey, ok := privateKey.(*rsa.PrivateKey)
	require.True(t, ok)

	return rsaPrivateKey
}

func readCertificate(t require.TestingT, path string) *x509.Certificate {
	b, err := os.ReadFile(path)
	require.NoError(t, err)

	block, _ := pem.Decode(b)
	require.True(t, block != nil)

	cert, err := x509.ParseCertificate(block.Bytes)
	require.NoError(t, err)

	return cert
}

func constructJWKSForTest(pubKey rsa.PublicKey, cert x509.Certificate) JWKSet {
	base64URLEncode := func(b []byte) string {
		encoded := base64.URLEncoding.EncodeToString(b)
		return strings.TrimRight(encoded, "=")
	}

	return JWKSet{
		Keys: []JWK{
			{
				KID: "json-web-key-01",
				Kty: "RSA",
				Use: "sig",
				N:   base64URLEncode(pubKey.N.Bytes()),
				E:   base64URLEncode(big.NewInt(int64(pubKey.E)).Bytes()),
				X5c: []string{
					base64.StdEncoding.EncodeToString(cert.Raw),
				},
				Alg: "RS256",
			},
		},
	}
}

type mockHTTPClient struct {
	mock.Mock
}

func (c *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := c.Called(req)

	if len(args) == 0 {
		panic("no return value specified for Do")
	}

	return args.Get(0).(*http.Response), args.Error(1)
}

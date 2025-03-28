package rsa

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	type arg struct {
		mockNewPrivateKeyFn func(keyPath string) (*rsa.PrivateKey, error)
		mockNewPublicKeyFn  func(keyPath string) (*rsa.PublicKey, error)

		givenPrivatePath string
		givenPublicPath  string
		expErr           error
	}
	tcs := map[string]arg{
		"err: unable to read private key": {
			mockNewPrivateKeyFn: func(keyPath string) (*rsa.PrivateKey, error) {
				return nil, errors.New("unable to read private key")
			},
			mockNewPublicKeyFn: func(keyPath string) (*rsa.PublicKey, error) {
				return nil, errors.New("unable to read public key")
			},
			expErr: errors.New("unable to read private key"),
		},
		"err: unable to read public key": {
			mockNewPrivateKeyFn: func(keyPath string) (*rsa.PrivateKey, error) {
				// Read file
				keyBytes, err := os.ReadFile(keyPath)
				if err != nil {
					return nil, err
				}

				// Decode
				block, _ := pem.Decode(keyBytes)
				parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
				if err != nil {
					return nil, err
				}
				key, ok := parsedKey.(*rsa.PrivateKey)
				if !ok {
					return nil, ErrNotRSAPrivateKey
				}
				return key, nil
			},
			mockNewPublicKeyFn: func(keyPath string) (*rsa.PublicKey, error) {
				return nil, errors.New("unable to read public key")
			},
			givenPrivatePath: "./testdata/rsa_private_key.pem",
			expErr:           errors.New("unable to read public key"),
		},
		"ok": {
			mockNewPrivateKeyFn: func(keyPath string) (*rsa.PrivateKey, error) {
				// Read file
				keyBytes, err := os.ReadFile(keyPath)
				if err != nil {
					return nil, err
				}

				// Decode
				block, _ := pem.Decode(keyBytes)
				parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
				if err != nil {
					return nil, err
				}

				key, ok := parsedKey.(*rsa.PrivateKey)
				if !ok {
					return nil, ErrNotRSAPrivateKey
				}
				return key, nil
			},
			mockNewPublicKeyFn: func(keyPath string) (*rsa.PublicKey, error) {
				// Read file
				keyBytes, err := os.ReadFile(keyPath)
				if err != nil {
					return nil, err
				}

				// Decode
				block, _ := pem.Decode(keyBytes)
				parsedKey, err := x509.ParsePKIXPublicKey(block.Bytes)
				if err != nil {
					return nil, err
				}

				key, ok := parsedKey.(*rsa.PublicKey)
				if !ok {
					return nil, ErrNotRSAPublicKey
				}
				return key, nil
			},
			givenPrivatePath: "./testdata/rsa_private_key.pem",
			givenPublicPath:  "./testdata/rsa_public_key.pem",
		},
	}
	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			// Given
			newPrivateKeyFn = tc.mockNewPrivateKeyFn
			newPublicKeyFn = tc.mockNewPublicKeyFn
			defer func() {
				newPrivateKeyFn = newPrivateKey
				newPublicKeyFn = newPublicKey
			}()

			// When
			instance, err := New(tc.givenPrivatePath, tc.givenPublicPath)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
				require.Nil(t, instance)
			} else {
				require.NoError(t, err)
				require.NotNil(t, instance)
			}
		})
	}
}

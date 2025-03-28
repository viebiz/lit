package rsa

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_newPrivateKey(t *testing.T) {
	type arg struct {
		givenPath string
		expErr    error
	}
	tcs := map[string]arg{
		"err: unable to read private key": {
			givenPath: "./testdata",
			expErr:    errors.New("read ./testdata: is a directory"),
		},
		"err: unable to parse private key": {
			givenPath: "./testdata/pkcs1_rsa_private_key.pem",
			expErr:    errors.New("x509: failed to parse private key (use ParsePKCS1PrivateKey instead for this key format)"),
		},
		"err: not RSA private key": {
			givenPath: "./testdata/ecdsa_private_key.pem",
			expErr:    ErrNotRSAPrivateKey,
		},
		"ok": {
			givenPath: "./testdata/rsa_private_key.pem",
		},
	}
	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()
			// Given

			// When
			instance, err := newPrivateKey(tc.givenPath)

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

func Test_newPublicKey(t *testing.T) {
	type arg struct {
		givenPath string
		expErr    error
	}
	tcs := map[string]arg{
		"err: unable to read public key": {
			givenPath: "./testdata",
			expErr:    errors.New("read ./testdata: is a directory"),
		},
		"err: unable to parse public key": {
			givenPath: "./testdata/pkcs1_rsa_public_key.pem",
			expErr:    errors.New("x509: failed to parse public key (use ParsePKCS1PublicKey instead for this key format)"),
		},
		"err: not RSA public key": {
			givenPath: "./testdata/ecdsa_public_key.pem",
			expErr:    ErrNotRSAPublicKey,
		},
		"ok": {
			givenPath: "./testdata/rsa_public_key.pem",
		},
	}
	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()
			// Given

			// When
			instance, err := newPublicKey(tc.givenPath)

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

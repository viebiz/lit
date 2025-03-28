package rsa

import (
	"crypto/rsa"
	"errors"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_impl_Encrypt(t *testing.T) {
	type arg struct {
		givenPublicKeyFn func() *rsa.PublicKey
		givenPublicKey   *rsa.PublicKey
		givenPlainText   string
		expEncryptedText string
		expErr           error
	}
	tcs := map[string]arg{
		"err: modulus must be > 1": {
			givenPublicKeyFn: func() *rsa.PublicKey {
				//publicKey, err := newPublicKey("./testdata/rsa_public_key.pem")
				//require.NoError(t, err)
				return &rsa.PublicKey{N: big.NewInt(0), E: 65537}
			},
			givenPlainText: "This is a plain text",
			expErr:         errors.New("modulus must be > 1"),
		},
		"ok": {
			givenPublicKeyFn: func() *rsa.PublicKey {
				publicKey, err := newPublicKey("./testdata/rsa_public_key.pem")
				require.NoError(t, err)
				return publicKey
			},
			givenPlainText: "This is a plain text",
		},
	}
	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()
			// Given

			// When
			instance := &impl{
				publicKey: tc.givenPublicKeyFn(),
			}
			result, err := instance.Encrypt(tc.givenPlainText)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, result)
			}
		})
	}
}

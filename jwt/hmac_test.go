package jwt

import (
	"crypto/ecdsa"
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSigningMethodHMAC_Sign(t *testing.T) {
	key := HMACPrivateKey("muryōkūsho")

	getSigningMethodFunc := func(alg string) HMAC {
		switch alg {
		case SigningMethodNameHS256:
			return NewHS256()
		case SigningMethodNameHS384:
			return NewHS384()
		case SigningMethodNameHS512:
			return NewHS512()
		default:
			panic("unknown algorithm")
		}
	}

	tcs := map[string]struct {
		givenString string
		key         Signer
		alg         string
		expResult   string
		expErr      error
	}{
		"success - HS256": {
			alg:         "HS256",
			key:         key,
			givenString: "developers-prefer-dark-mode-because-light-attracts-bugs",
			expResult:   "A23KkQvwS7apRRvpCEkhQsoCJHHDXhp-jd6ifUjsfjI",
		},
		"success - HS384": {
			alg:         "HS384",
			key:         key,
			givenString: "developers-prefer-dark-mode-because-light-attracts-bugs",
			expResult:   "UekFXNPzrjqRzU7TiIGaqmPv3JNbe_sEoOGeUyBp-obRZBkLUAznORVYWq6hEq_Q",
		},
		"success - HS512": {
			alg:         "HS512",
			key:         key,
			givenString: "developers-prefer-dark-mode-because-light-attracts-bugs",
			expResult:   "qve9yY52qXd2lE5fRh44Vq2aiaHqDSimodQ_msCWIuGeOt2a4w3giHt-bDPEDiU23-fX0_3wfzoLNqfY6bE2Qg",
		},
		"error - key type not supported": {
			alg:    "HS512",
			key:    &ecdsa.PrivateKey{},
			expErr: ErrInvalidKeyType,
		},
	}

	for scenario, tc := range tcs {
		t.Run(scenario, func(t *testing.T) {
			// Given
			method := getSigningMethodFunc(tc.alg)

			// When
			actualResult, err := method.Sign([]byte(tc.givenString), tc.key)

			// Then
			if tc.expErr != nil {
				require.Equal(t, tc.expErr, err)
			} else {
				require.NoError(t, err)

				expResult, err := base64.RawURLEncoding.DecodeString(tc.expResult)
				require.NoError(t, err)
				require.Equal(t, expResult, actualResult)
			}
		})
	}
}

func TestSigningMethodHMAC_Verify(t *testing.T) {
	key := HMACPrivateKey("muryōkūsho")

	getSigningMethodFunc := func(alg string) HMAC {
		switch alg {
		case SigningMethodNameHS256:
			return NewHS256()
		case SigningMethodNameHS384:
			return NewHS384()
		case SigningMethodNameHS512:
			return NewHS512()
		default:
			panic("unknown algorithm")
		}
	}

	tcs := map[string]struct {
		alg           string
		signingString string
		signature     string
		key           VerifyKey
		expErr        error
	}{
		"success - HS256": {
			alg:           "HS256",
			key:           key,
			signingString: "developers-prefer-dark-mode-because-light-attracts-bugs",
			signature:     "A23KkQvwS7apRRvpCEkhQsoCJHHDXhp-jd6ifUjsfjI",
		},
		"success - HS384": {
			alg:           "HS384",
			key:           key,
			signingString: "developers-prefer-dark-mode-because-light-attracts-bugs",
			signature:     "UekFXNPzrjqRzU7TiIGaqmPv3JNbe_sEoOGeUyBp-obRZBkLUAznORVYWq6hEq_Q",
		},
		"success - HS512": {
			alg:           "HS512",
			key:           key,
			signingString: "developers-prefer-dark-mode-because-light-attracts-bugs",
			signature:     "qve9yY52qXd2lE5fRh44Vq2aiaHqDSimodQ_msCWIuGeOt2a4w3giHt-bDPEDiU23-fX0_3wfzoLNqfY6bE2Qg",
		},
		"error - invalid key type": {
			alg:    "HS512",
			key:    &ecdsa.PublicKey{},
			expErr: ErrInvalidKeyType,
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given
			method := getSigningMethodFunc(tc.alg)

			// When
			sig, err := base64.RawURLEncoding.DecodeString(tc.signature)
			require.NoError(t, err)

			actualErr := method.Verify([]byte(tc.signingString), sig, tc.key)

			// Then
			if tc.expErr != nil {
				require.Equal(t, tc.expErr, actualErr)
			} else {
				require.NoError(t, actualErr)
			}
		})
	}
}

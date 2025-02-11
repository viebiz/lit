package jwt

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSigningMethodRSA_Sign(t *testing.T) {
	privateKeyPath := "testdata/sample_rsa_private_key"
	key := readKeyForTest[*rsa.PrivateKey](t, privateKeyPath)

	getSigningMethodFunc := func(alg string) RSA {
		switch alg {
		case SigningMethodNameRS256:
			return NewRS256()
		case SigningMethodNameRS384:
			return NewRS384()
		case SigningMethodNameRS512:
			return NewRS512()
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
		"success - RS256": {
			alg:         "RS256",
			key:         key,
			givenString: "developers-prefer-dark-mode-because-light-attracts-bugs",
			expResult:   "nUuMsCv-zqzKg8Gvv1ixnn3gkWy6cAypEGKijVtChNbwxlW3Cwp2De6SFrqrA3L6sE6ytzUfp19EyWykpKIS0ta03Zkvboh0QJc6iSyGQAr2iJ7H33d2av87Nai1PY4tXnQQhaaoO305e-YNKygYcOCHBhJ0QfYClb_Bb7ShfEMy_XKX4NpuCiiqdpSj7WZva_Knb3cDOh7Q_W2kSNaFAaBpspcjR156EHMiJlJvf-JHJ5fcWXUCXromS562CPV_fOrZSmy_9rGbbiWhuea8Uy5IqIvyj4VtJE9SLk2dWehqHyNiDT8vBPSKq5e9gUUt1_UOB7nbCOsJ26Ny9alICA",
		},
		"success - RS384": {
			alg:         "RS384",
			key:         key,
			givenString: "developers-prefer-dark-mode-because-light-attracts-bugs",
			expResult:   "k9xcZ1XSeJyWQ5aC9sU4AK7o73avUNKZjAs-9HpxH2h96zmz8b9hO79c_B14IPzmTQJ9cD8yF91bKGBFgmI-wXWdMc_-pdAGYSr4rw_TOBvyofhPVsxDg96ejZBYCxtCa2FBc1OOLF3KTIsPNO1Xx3N2lsPQfmjZHVDsXCcy6oIO2tMIkagVMYMqwX4RvuEht719Rtt4e654kpYNe1lEWas4Rnsd4yWPmQZrH_Seli6Y4WahSsp6jjL5s8gLNPznWHYoKRnMtqjkaUA4DwJgrOYd4s3DXYOeBLwlpTHAUtgZO-rjhIodaCj8VNXS268CrNmAiC4P9zsBJ1YNlmsozw",
		},
		"success - RS512": {
			alg:         "RS512",
			key:         key,
			givenString: "developers-prefer-dark-mode-because-light-attracts-bugs",
			expResult:   "PJQz1kGiDhpt6z8b9jrPu8kjZSn-R0rpPsxa0NmrmCu9MiTPip-o_MZc6zepLLmVVYo9eMVJnzqk4V0qdcFnohx0VjT_xBED3mDVZZxU0crWlk1-I5qUhnprSzW5oGD-GmB-fT29Br7GotD4FDwHfOxsnyFFlPIkl5AHPZRKKXcWT5qU5Q-7RP4b-fmA5wF9ktTG_726gWWXlz9UE2v6xgDoL49qjTwyLQcV3ZwCrim0WCk1fHIk2b5MPKcwWSu9SXnikjG5vPQLMVfIywrkHt370Zz_ea-EbCgCasYD08l8rK8xe5juF2IxVU2iJTBPt_8v8T-fr-kZvVdwBGADXw",
		},
		"error - key type not supported": {
			alg:    "RS512",
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

func TestSigningMethodRSA_Verify(t *testing.T) {
	privateKeyPath := "testdata/sample_rsa_private_key"
	key := readKeyForTest[*rsa.PrivateKey](t, privateKeyPath)

	getSigningMethodFunc := func(alg string) RSA {
		switch alg {
		case SigningMethodNameRS256:
			return NewRS256()
		case SigningMethodNameRS384:
			return NewRS384()
		case SigningMethodNameRS512:
			return NewRS512()
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
		"success - RS256": {
			alg:           "RS256",
			key:           key.Public(),
			signingString: "developers-prefer-dark-mode-because-light-attracts-bugs",
			signature:     "nUuMsCv-zqzKg8Gvv1ixnn3gkWy6cAypEGKijVtChNbwxlW3Cwp2De6SFrqrA3L6sE6ytzUfp19EyWykpKIS0ta03Zkvboh0QJc6iSyGQAr2iJ7H33d2av87Nai1PY4tXnQQhaaoO305e-YNKygYcOCHBhJ0QfYClb_Bb7ShfEMy_XKX4NpuCiiqdpSj7WZva_Knb3cDOh7Q_W2kSNaFAaBpspcjR156EHMiJlJvf-JHJ5fcWXUCXromS562CPV_fOrZSmy_9rGbbiWhuea8Uy5IqIvyj4VtJE9SLk2dWehqHyNiDT8vBPSKq5e9gUUt1_UOB7nbCOsJ26Ny9alICA",
		},
		"success - RS384": {
			alg:           "RS384",
			key:           key.Public(),
			signingString: "developers-prefer-dark-mode-because-light-attracts-bugs",
			signature:     "k9xcZ1XSeJyWQ5aC9sU4AK7o73avUNKZjAs-9HpxH2h96zmz8b9hO79c_B14IPzmTQJ9cD8yF91bKGBFgmI-wXWdMc_-pdAGYSr4rw_TOBvyofhPVsxDg96ejZBYCxtCa2FBc1OOLF3KTIsPNO1Xx3N2lsPQfmjZHVDsXCcy6oIO2tMIkagVMYMqwX4RvuEht719Rtt4e654kpYNe1lEWas4Rnsd4yWPmQZrH_Seli6Y4WahSsp6jjL5s8gLNPznWHYoKRnMtqjkaUA4DwJgrOYd4s3DXYOeBLwlpTHAUtgZO-rjhIodaCj8VNXS268CrNmAiC4P9zsBJ1YNlmsozw",
		},
		"success - RS512": {
			alg:           "RS512",
			key:           key.Public(),
			signingString: "developers-prefer-dark-mode-because-light-attracts-bugs",
			signature:     "PJQz1kGiDhpt6z8b9jrPu8kjZSn-R0rpPsxa0NmrmCu9MiTPip-o_MZc6zepLLmVVYo9eMVJnzqk4V0qdcFnohx0VjT_xBED3mDVZZxU0crWlk1-I5qUhnprSzW5oGD-GmB-fT29Br7GotD4FDwHfOxsnyFFlPIkl5AHPZRKKXcWT5qU5Q-7RP4b-fmA5wF9ktTG_726gWWXlz9UE2v6xgDoL49qjTwyLQcV3ZwCrim0WCk1fHIk2b5MPKcwWSu9SXnikjG5vPQLMVfIywrkHt370Zz_ea-EbCgCasYD08l8rK8xe5juF2IxVU2iJTBPt_8v8T-fr-kZvVdwBGADXw",
		},
		"error - invalid key type": {
			alg:    "RS256",
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

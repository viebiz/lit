package jwt

import (
	"crypto"
	"crypto/rsa"
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/viebiz/lit/ioutil"

	"github.com/stretchr/testify/require"
)

func TestParser_Parse(t *testing.T) {
	ioutil.SetResourceDir("testdata")
	keyFunc := func(string) (crypto.PublicKey, error) {
		privateKeyPath := "sample_rsa_private_key"
		key := readKeyForTest[*rsa.PrivateKey](t, privateKeyPath)

		return key.Public(), nil
	}

	iat := time.Date(2024, time.July, 24, 0, 0, 0, 0, time.UTC).Unix()
	exp := time.Date(2024, time.July, 24, 1, 0, 0, 0, time.UTC).Unix()

	defer func(origin func() time.Time) { timeNowFunc = origin }(timeNowFunc)
	timeNowFunc = func() time.Time { return time.Unix(iat+50, 0) }

	tcs := map[string]struct {
		inputTokenString string
		inputKeyFunc     func(string) (crypto.PublicKey, error)
		expToken         Token[RegisteredClaims]
		expError         error
	}{
		"success": {
			inputTokenString: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwczovL2xpbWl0bGVzcy5tdWthZ2VuLmNvbSIsImF1ZCI6WyJodHRwczovL3Jlc291cmNlLWFwaS5jb20iXSwiaWF0IjoxNzIxNzc5MjAwLCJleHAiOjE3MjE3ODI4MDB9.dnmh7SMxs3QkEJ5pUjjYeZ9CEm8x2XWNF0_YuDO46cbm3KIB5wrJBUJ5BuRVrkOAW39ZDxaamgLJunciEPT9BR7j0dGPWguZ4EZausHZt1ehn8OGZqyd6Uj16xOqkRWndsU7kaMJuKlLUFizHG305xcrh9M5PAmJe4PMZxE84SWYcj0QpCZ58zpWXA-OWTJzVSbpkbNMfl6RKsKxQ9DPDEl8wT6JorLC18Ov_MeV6KCSuwr4f15zaSLKlUv6I5n00PtD9Uw7hG9vbTDe0LHhG9WRbtVe-8Mqyz_pPV-oMMF28B1bRY48ItBKrJphSpE8QSsIGBHXdI5Es8hYYJyKOA",
			inputKeyFunc:     keyFunc,
			expToken: Token[RegisteredClaims]{
				Header: map[string]string{
					"alg": "RS256",
					"typ": "JWT",
				},
				Claims: RegisteredClaims{
					Issuer:    "https://limitless.mukagen.com",
					Audience:  ClaimStrings([]string{"https://resource-api.com"}),
					IssuedAt:  &iat,
					ExpiresAt: &exp,
				},
			},
		},
		"error - token malformed": {
			inputTokenString: "",
			inputKeyFunc:     keyFunc,

			expError: ErrTokenMalformed,
		},
		"error - missing header alg": {
			inputTokenString: "eyJ0eXAiOiJKV1QifQ.eyJpc3MiOiJodHRwczovL2xpbWl0bGVzcy5tdWthZ2VuLmNvbSIsImF1ZCI6WyJodHRwczovL3Jlc291cmNlLWFwaS5jb20iXSwiaWF0IjoxNzIxNjkyODAwLCJleHAiOjE3MjE2OTY0MDB9.i5-l4fRbuABk7Dzs9NNh3SOPYgUc9plzesFL6lXSngJSwMi6zD9J-qu6lNFAAHxbqeE5riE51vsg829A1DTldUFKJLeQKi33nBG8AN9yk_d3v4XoKiW2cHYHb5hwLu4f1Evk55uMLTSymM_ygp-FAmdAWw66aoot_wSzKkZ0pefHjIsT5n8hC_5YnYFDse90UdJ1zSG5inZo-vubehJk4fnikPNmpePzPHhzlGvYg9DoxcQAXpWxIKjNuZwhcSEwejj3p1vM3pXY-cRrZT-axmipMD90cc4UfzPzYK7RvO4eoWfBQip7h5-yNEZUjfvT9vnSn5ErNqYHG9VI87xlQw",
			inputKeyFunc:     keyFunc,
			expError:         ErrInvalidToken,
		},
		"error - signing method not supported": {
			inputTokenString: "eyJhbGciOiJVTlNVUFBPUlRFRCIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwczovL2xpbWl0bGVzcy5tdWthZ2VuLmNvbSIsImF1ZCI6WyJodHRwczovL3Jlc291cmNlLWFwaS5jb20iXSwiaWF0IjoxNzIxNjkyODAwLCJleHAiOjE3MjE2OTY0MDB9.i0-DiHJWjzCUGQMZD6niy18jmv8ACq_DZtSAlNBvLja1cpYWrgI8xFiLGJkVdYudUL4pMO817Mr9fZunDFy5kuDQq8G9iK3YWM7AbUgBIByKtPmwCmzEHB-5chIs3pCQPpoaTusFPv83jjTK72inpOtMcwhT-uadjkPXLJvaNaKKdBq6P3LLI4nUpIn_-PD8DrFL2BQOslIdPN-fy_Jg4-PCdbStQpM4Zm3XB5qgwKL-nxfbCwXVqwOHgMkh6KVQMQP8G2HZ_qkxsZpbYNs1s0ihIYNucCvG63gzGDlibGxEhnFjGme_dWjGogMsd0zRTXGrtO-L19DxtJ5lLUDsWA",
			inputKeyFunc:     keyFunc,
			expError:         ErrSigningMethodNotSupported,
		},
		"error - getKeyFunc error": {
			inputTokenString: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwczovL2xpbWl0bGVzcy5tdWthZ2VuLmNvbSIsImF1ZCI6WyJodHRwczovL3Jlc291cmNlLWFwaS5jb20iXSwiaWF0IjoxNzIxNzc5MjAwLCJleHAiOjE3MjE3ODI4MDB9.dnmh7SMxs3QkEJ5pUjjYeZ9CEm8x2XWNF0_YuDO46cbm3KIB5wrJBUJ5BuRVrkOAW39ZDxaamgLJunciEPT9BR7j0dGPWguZ4EZausHZt1ehn8OGZqyd6Uj16xOqkRWndsU7kaMJuKlLUFizHG305xcrh9M5PAmJe4PMZxE84SWYcj0QpCZ58zpWXA-OWTJzVSbpkbNMfl6RKsKxQ9DPDEl8wT6JorLC18Ov_MeV6KCSuwr4f15zaSLKlUv6I5n00PtD9Uw7hG9vbTDe0LHhG9WRbtVe-8Mqyz_pPV-oMMF28B1bRY48ItBKrJphSpE8QSsIGBHXdI5Es8hYYJyKOA",
			inputKeyFunc: func(s string) (crypto.PublicKey, error) {
				return nil, errors.New("simulated error")
			},
			expError: errors.New("simulated error"),
		},
		"error - invalid signature": {
			inputTokenString: "eyJhbGciOiJSUzI1NiIsInR5cCI6ImF0K2p3dCIsImtpZCI6Imh1YnAxOUhWZnQ4cTRYYWxrVmYtTyJ9.eyJpc3MiOiJodHRwczovL2Rldi13aXRjaGVyLnVzLmF1dGgwLmNvbS8iLCJzdWIiOiJHOG1EREVSVlZkakJYZzJST0Q1QUkxSjIwSHgzTjU3bUBjbGllbnRzIiwiYXVkIjoiaHR0cHM6Ly9saW1pdGxlc3MubXVrYWdlbi5jb20iLCJpYXQiOjE3Mjc2MjQzNzgsImV4cCI6MTcyNzcxMDc3OCwic2NvcGUiOiJyZWFkOnVzZXJzIHdyaXRlOnVzZXJzIiwianRpIjoibm5xYktTM1hGV1RTaERyUVV2WkRyNyIsImNsaWVudF9pZCI6Ikc4bURERVJWVmRqQlhnMlJPRDVBSTFKMjBIeDNONTdtIn0.BfALvjEPZ40J2h6L2fknf7wYsKGPbjqMolH7-O-HVVK-9Pj8fEuyDAETHDqIlfQaN6hZV1I8iTSNgX_OglrtfMQ93mwQ9ToSi8bwVVsyrVWGite_4MU7bjHlLZkgfqbw81uzPOrTZfdFnCTjkrXLk98IkchRXa3s_AX8s-SGjFkp_hyGh3lI-M5hPcuoCnhQY16kH4DzFmE_d4UBGdlrXwluSx4JlM8DfOm75oZs6Ts3EMuPCwBVz2hbV0zLC9ynFzj0LF5CAMsY4HWxmSsiNuCoSvzyiYKhnnkNCiTU6tfMmdNEI43JS62U8Z-CF89Ubc0Eym3cz9DFRNjPl8aZUQ",
			inputKeyFunc:     keyFunc,
			expError:         ErrInvalidSignature,
		},
		"error - token expired": {
			inputTokenString: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwczovL2xpbWl0bGVzcy5tdWthZ2VuLmNvbSIsImF1ZCI6WyJodHRwczovL3Jlc291cmNlLWFwaS5jb20iXSwiaWF0IjoxNzIxNjkyODAwLCJleHAiOjE3MjE2OTY0MDB9.nVUYXq9IWYIr-UiJF6dNWITf78THC3VnakXg4goLZ3OvV-z4WsZpARz4rflGJMzegB2by8qBVdIEJK7XmIGqbP-QT3xJK0MzynSBAcfMIhczUqrw8oeA0myhXT08PmIlI6Vc6EFHfux0j7Ju5U3JwBFOIN09twrrQrUBwy8W7quqH3ZtiVFDiQfw5tu-VtEuD-ohdm0j4TvDcST16e48X8Jo6QGCkGwNYGya_tlFhYwgB-3xY_EKENv8gXxTIRf91mO07UXjcvvtOtxjqBPWPw7PGdPIdDmyyVxzKrH9xB1kFpLieE45Iijn8EgGFx8pML6i-kNXfisO_95HvhdcrA",

			inputKeyFunc: keyFunc,
			expError:     ErrTokenExpired,
		},
	}

	for scenario, tc := range tcs {
		t.Run(scenario, func(t *testing.T) {
			// Given
			parser := NewParser[RegisteredClaims]()

			// When
			tk, err := parser.Parse(tc.inputTokenString, tc.inputKeyFunc)

			// Then
			if tc.expError != nil {
				require.EqualError(t, err, tc.expError.Error())
			} else {
				require.NoError(t, err)
				ignoreFieldsOpts := cmpopts.IgnoreFields(Token[RegisteredClaims]{}, "Signature", "signingMethod")

				if !cmp.Equal(tc.expToken, tk, ignoreFieldsOpts) {
					t.Errorf("\n result mismatched. Diff: %+v", cmp.Diff(tc.expToken, tk, ignoreFieldsOpts))
					t.FailNow()
				}
			}
		})
	}
}

func TestNewParser(t *testing.T) {
	t.Run("default parser", func(t *testing.T) {
		// When
		parser := NewParser[RegisteredClaims]()

		// Then
		require.NotNil(t, parser.signingMethods)
		require.Contains(t, parser.signingMethods, SigningMethodNameRS256)
		require.Contains(t, parser.signingMethods, SigningMethodNameHS256)
		require.Equal(t, NewRS256(), parser.signingMethods[SigningMethodNameRS256])
		require.Equal(t, NewHS256(), parser.signingMethods[SigningMethodNameHS256])
	})
}

func TestNewDefaultParser(t *testing.T) {
	t.Run("default parser", func(t *testing.T) {
		// When
		parser := NewDefaultParser[RegisteredClaims]()

		// Then
		require.NotNil(t, parser.signingMethods)
		require.Len(t, parser.signingMethods, 6) // RS256, RS384, RS512, HS256, HS384, HS512
		require.Contains(t, parser.signingMethods, SigningMethodNameRS256)
		require.Contains(t, parser.signingMethods, SigningMethodNameRS384)
		require.Contains(t, parser.signingMethods, SigningMethodNameRS512)
		require.Contains(t, parser.signingMethods, SigningMethodNameHS256)
		require.Contains(t, parser.signingMethods, SigningMethodNameHS384)
		require.Contains(t, parser.signingMethods, SigningMethodNameHS512)
	})
}

func TestWithSigningMethod(t *testing.T) {
	t.Run("add custom signing method", func(t *testing.T) {
		// Given
		customMethod := NewHS256()
		customMethod.Name = "CUSTOM"

		// When
		opt := WithSigningMethod(customMethod)
		parser := NewDefaultParser[RegisteredClaims]()
		opt((*Parser[Claims])(&parser))

		// Then
		require.Contains(t, parser.signingMethods, "CUSTOM")
		require.Equal(t, customMethod, parser.signingMethods["CUSTOM"])
	})
}

func TestWithSigningMethods(t *testing.T) {
	t.Run("override signing methods", func(t *testing.T) {
		// Given
		customMethods := map[string]SigningMethod{
			"CUSTOM1": NewHS256(),
			"CUSTOM2": NewRS256(),
		}

		// When
		opt := WithSigningMethods(customMethods)
		parser := NewDefaultParser[RegisteredClaims]()
		opt((*Parser[Claims])(&parser))

		// Then
		require.Len(t, parser.signingMethods, 2)
		require.Contains(t, parser.signingMethods, "CUSTOM1")
		require.Contains(t, parser.signingMethods, "CUSTOM2")
		require.NotContains(t, parser.signingMethods, SigningMethodNameRS256)
	})
}

func TestParser_parseToken(t *testing.T) {
	parser := NewDefaultParser[RegisteredClaims]()

	tcs := map[string]struct {
		tokenString string
		expError    error
		expHeader   map[string]string
		expClaims   RegisteredClaims
	}{
		"valid token": {
			tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0LXN1YmplY3QiLCJpc3MiOiJ0ZXN0LWlzc3VlciJ9.dGVzdC1zaWduYXR1cmU", // signature is base64 of "test-signature"
			expHeader: map[string]string{
				"alg": "HS256",
				"typ": "JWT",
			},
			expClaims: RegisteredClaims{
				Subject: "test-subject",
				Issuer:  "test-issuer",
			},
		},
		"malformed token - too few parts": {
			tokenString: "header.payload",
			expError:    ErrTokenMalformed,
		},
		"malformed token - too many parts": {
			tokenString: "header.payload.signature.extra",
			expError:    ErrTokenMalformed,
		},
		"invalid header json": {
			tokenString: "aW52YWxpZC1qc29u.payload.signature", // base64 of "invalid-json"
			expError:    nil,                                  // Will be a json.SyntaxError
		},
		"invalid payload json": {
			tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.aW52YWxpZC1qc29u.signature", // base64 of "invalid-json"
			expError:    nil,                                                               // Will be a json.SyntaxError
		},
	}

	for scenario, tc := range tcs {
		t.Run(scenario, func(t *testing.T) {
			// When
			token, signingString, err := parser.parseToken(tc.tokenString)

			// Then
			if tc.expError != nil {
				require.Error(t, err)
				require.Equal(t, tc.expError, err)
			} else if tc.expError == nil && tc.expHeader == nil {
				// Special case for invalid JSON tests - expect error but don't check type
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expHeader, token.Header)
				require.Equal(t, tc.expClaims, token.Claims)
				require.Contains(t, string(signingString), ".")
			}
		})
	}
}

func TestParser_getSigningMethod(t *testing.T) {
	parser := NewDefaultParser[RegisteredClaims]()

	tcs := map[string]struct {
		alg       string
		expMethod SigningMethod
		expFound  bool
	}{
		"existing method": {
			alg:       SigningMethodNameHS256,
			expMethod: NewHS256(),
			expFound:  true,
		},
		"non-existing method": {
			alg:      "NONEXISTENT",
			expFound: false,
		},
	}

	for scenario, tc := range tcs {
		t.Run(scenario, func(t *testing.T) {
			// When
			method, found := parser.getSigningMethod(tc.alg)

			// Then
			require.Equal(t, tc.expFound, found)
			if tc.expFound {
				require.Equal(t, tc.expMethod, method)
			}
		})
	}
}

func Test_decodeSegment(t *testing.T) {
	tcs := map[string]struct {
		input    string
		expected string
		expError bool
	}{
		"valid base64": {
			input:    "SGVsbG8gV29ybGQ", // "Hello World" in base64
			expected: "Hello World",
			expError: false,
		},
		"empty string": {
			input:    "",
			expected: "",
			expError: false,
		},
		"invalid base64": {
			input:    "invalid-base64!",
			expError: true,
		},
	}

	for scenario, tc := range tcs {
		t.Run(scenario, func(t *testing.T) {
			// When
			result, err := decodeSegment([]byte(tc.input))

			// Then
			if tc.expError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected, string(result))
			}
		})
	}
}

func BenchmarkParser_Parse(b *testing.B) {
	ioutil.SetResourceDir("testdata")
	keyFunc := func(string) (crypto.PublicKey, error) {
		privateKeyPath := "sample_rsa_private_key"
		key := readKeyForTest[*rsa.PrivateKey](b, privateKeyPath)

		return key.Public(), nil
	}

	iat := time.Date(2024, time.July, 24, 0, 0, 0, 0, time.UTC).Unix()
	exp := time.Date(2024, time.July, 24, 1, 0, 0, 0, time.UTC).Unix()

	defer func(origin func() time.Time) { timeNowFunc = origin }(timeNowFunc)
	timeNowFunc = func() time.Time { return time.Unix(iat+50, 0) }

	parser := NewParser[RegisteredClaims]()

	tcs := map[string]struct {
		inputTokenString string
		inputKeyFunc     func(string) (crypto.PublicKey, error)
		expToken         Token[RegisteredClaims]
		expError         error
	}{
		"success": {
			inputTokenString: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwczovL2xpbWl0bGVzcy5tdWthZ2VuLmNvbSIsImF1ZCI6WyJodHRwczovL3Jlc291cmNlLWFwaS5jb20iXSwiaWF0IjoxNzIxNzc5MjAwLCJleHAiOjE3MjE3ODI4MDB9.dnmh7SMxs3QkEJ5pUjjYeZ9CEm8x2XWNF0_YuDO46cbm3KIB5wrJBUJ5BuRVrkOAW39ZDxaamgLJunciEPT9BR7j0dGPWguZ4EZausHZt1ehn8OGZqyd6Uj16xOqkRWndsU7kaMJuKlLUFizHG305xcrh9M5PAmJe4PMZxE84SWYcj0QpCZ58zpWXA-OWTJzVSbpkbNMfl6RKsKxQ9DPDEl8wT6JorLC18Ov_MeV6KCSuwr4f15zaSLKlUv6I5n00PtD9Uw7hG9vbTDe0LHhG9WRbtVe-8Mqyz_pPV-oMMF28B1bRY48ItBKrJphSpE8QSsIGBHXdI5Es8hYYJyKOA",
			inputKeyFunc:     keyFunc,
			expToken: Token[RegisteredClaims]{
				Header: map[string]string{
					"alg": "RS256",
					"typ": "JWT",
				},
				Claims: RegisteredClaims{
					Issuer:    "https://limitless.mukagen.com",
					Audience:  ClaimStrings([]string{"https://resource-api.com"}),
					IssuedAt:  &iat,
					ExpiresAt: &exp,
				},
			},
		},
		"error - token malformed": {
			inputTokenString: "",
			inputKeyFunc:     keyFunc,

			expError: ErrTokenMalformed,
		},
		"error - missing header alg": {
			inputTokenString: "eyJ0eXAiOiJKV1QifQ.eyJpc3MiOiJodHRwczovL2xpbWl0bGVzcy5tdWthZ2VuLmNvbSIsImF1ZCI6WyJodHRwczovL3Jlc291cmNlLWFwaS5jb20iXSwiaWF0IjoxNzIxNjkyODAwLCJleHAiOjE3MjE2OTY0MDB9.i5-l4fRbuABk7Dzs9NNh3SOPYgUc9plzesFL6lXSngJSwMi6zD9J-qu6lNFAAHxbqeE5riE51vsg829A1DTldUFKJLeQKi33nBG8AN9yk_d3v4XoKiW2cHYHb5hwLu4f1Evk55uMLTSymM_ygp-FAmdAWw66aoot_wSzKkZ0pefHjIsT5n8hC_5YnYFDse90UdJ1zSG5inZo-vubehJk4fnikPNmpePzPHhzlGvYg9DoxcQAXpWxIKjNuZwhcSEwejj3p1vM3pXY-cRrZT-axmipMD90cc4UfzPzYK7RvO4eoWfBQip7h5-yNEZUjfvT9vnSn5ErNqYHG9VI87xlQw",
			inputKeyFunc:     keyFunc,
			expError:         ErrInvalidToken,
		},
		"error - signing method not supported": {
			inputTokenString: "eyJhbGciOiJVTlNVUFBPUlRFRCIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwczovL2xpbWl0bGVzcy5tdWthZ2VuLmNvbSIsImF1ZCI6WyJodHRwczovL3Jlc291cmNlLWFwaS5jb20iXSwiaWF0IjoxNzIxNjkyODAwLCJleHAiOjE3MjE2OTY0MDB9.i0-DiHJWjzCUGQMZD6niy18jmv8ACq_DZtSAlNBvLja1cpYWrgI8xFiLGJkVdYudUL4pMO817Mr9fZunDFy5kuDQq8G9iK3YWM7AbUgBIByKtPmwCmzEHB-5chIs3pCQPpoaTusFPv83jjTK72inpOtMcwhT-uadjkPXLJvaNaKKdBq6P3LLI4nUpIn_-PD8DrFL2BQOslIdPN-fy_Jg4-PCdbStQpM4Zm3XB5qgwKL-nxfbCwXVqwOHgMkh6KVQMQP8G2HZ_qkxsZpbYNs1s0ihIYNucCvG63gzGDlibGxEhnFjGme_dWjGogMsd0zRTXGrtO-L19DxtJ5lLUDsWA",
			inputKeyFunc:     keyFunc,
			expError:         ErrSigningMethodNotSupported,
		},
		"error - getKeyFunc error": {
			inputTokenString: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwczovL2xpbWl0bGVzcy5tdWthZ2VuLmNvbSIsImF1ZCI6WyJodHRwczovL3Jlc291cmNlLWFwaS5jb20iXSwiaWF0IjoxNzIxNzc5MjAwLCJleHAiOjE3MjE3ODI4MDB9.dnmh7SMxs3QkEJ5pUjjYeZ9CEm8x2XWNF0_YuDO46cbm3KIB5wrJBUJ5BuRVrkOAW39ZDxaamgLJunciEPT9BR7j0dGPWguZ4EZausHZt1ehn8OGZqyd6Uj16xOqkRWndsU7kaMJuKlLUFizHG305xcrh9M5PAmJe4PMZxE84SWYcj0QpCZ58zpWXA-OWTJzVSbpkbNMfl6RKsKxQ9DPDEl8wT6JorLC18Ov_MeV6KCSuwr4f15zaSLKlUv6I5n00PtD9Uw7hG9vbTDe0LHhG9WRbtVe-8Mqyz_pPV-oMMF28B1bRY48ItBKrJphSpE8QSsIGBHXdI5Es8hYYJyKOA",
			inputKeyFunc: func(s string) (crypto.PublicKey, error) {
				return nil, errors.New("simulated error")
			},
			expError: errors.New("simulated error"),
		},
		"error - invalid signature": {
			inputTokenString: "eyJhbGciOiJSUzI1NiIsInR5cCI6ImF0K2p3dCIsImtpZCI6Imh1YnAxOUhWZnQ4cTRYYWxrVmYtTyJ9.eyJpc3MiOiJodHRwczovL2Rldi13aXRjaGVyLnVzLmF1dGgwLmNvbS8iLCJzdWIiOiJHOG1EREVSVlZkakJYZzJST0Q1QUkxSjIwSHgzTjU3bUBjbGllbnRzIiwiYXVkIjoiaHR0cHM6Ly9saW1pdGxlc3MubXVrYWdlbi5jb20iLCJpYXQiOjE3Mjc2MjQzNzgsImV4cCI6MTcyNzcxMDc3OCwic2NvcGUiOiJyZWFkOnVzZXJzIHdyaXRlOnVzZXJzIiwianRpIjoibm5xYktTM1hGV1RTaERyUVV2WkRyNyIsImNsaWVudF9pZCI6Ikc4bURERVJWVmRqQlhnMlJPRDVBSTFKMjBIeDNONTdtIn0.BfALvjEPZ40J2h6L2fknf7wYsKGPbjqMolH7-O-HVVK-9Pj8fEuyDAETHDqIlfQaN6hZV1I8iTSNgX_OglrtfMQ93mwQ9ToSi8bwVVsyrVWGite_4MU7bjHlLZkgfqbw81uzPOrTZfdFnCTjkrXLk98IkchRXa3s_AX8s-SGjFkp_hyGh3lI-M5hPcuoCnhQY16kH4DzFmE_d4UBGdlrXwluSx4JlM8DfOm75oZs6Ts3EMuPCwBVz2hbV0zLC9ynFzj0LF5CAMsY4HWxmSsiNuCoSvzyiYKhnnkNCiTU6tfMmdNEI43JS62U8Z-CF89Ubc0Eym3cz9DFRNjPl8aZUQ",
			inputKeyFunc:     keyFunc,
			expError:         ErrInvalidSignature,
		},
		"error - token expired": {
			inputTokenString: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwczovL2xpbWl0bGVzcy5tdWthZ2VuLmNvbSIsImF1ZCI6WyJodHRwczovL3Jlc291cmNlLWFwaS5jb20iXSwiaWF0IjoxNzIxNjkyODAwLCJleHAiOjE3MjE2OTY0MDB9.nVUYXq9IWYIr-UiJF6dNWITf78THC3VnakXg4goLZ3OvV-z4WsZpARz4rflGJMzegB2by8qBVdIEJK7XmIGqbP-QT3xJK0MzynSBAcfMIhczUqrw8oeA0myhXT08PmIlI6Vc6EFHfux0j7Ju5U3JwBFOIN09twrrQrUBwy8W7quqH3ZtiVFDiQfw5tu-VtEuD-ohdm0j4TvDcST16e48X8Jo6QGCkGwNYGya_tlFhYwgB-3xY_EKENv8gXxTIRf91mO07UXjcvvtOtxjqBPWPw7PGdPIdDmyyVxzKrH9xB1kFpLieE45Iijn8EgGFx8pML6i-kNXfisO_95HvhdcrA",

			inputKeyFunc: keyFunc,
			expError:     ErrTokenExpired,
		},
	}
	for scenario, tc := range tcs {
		b.Run(scenario, func(b *testing.B) {
			b.Helper()
			b.ReportAllocs()
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_, _ = parser.Parse(tc.inputTokenString, tc.inputKeyFunc)
				}
			})
		})
	}
}

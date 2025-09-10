package jwt

import (
	"crypto"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/viebiz/lit/ioutil"
)

func TestToken_SignedString(t *testing.T) {
	ioutil.SetResourceDir("testdata")
	privateKeyPath := "sample_rsa_private_key"
	key := readKeyForTest[*rsa.PrivateKey](t, privateKeyPath)
	iat := time.Date(2024, time.July, 24, 0, 0, 0, 0, time.UTC).Unix()
	exp := time.Date(2024, time.July, 24, 1, 0, 0, 0, time.UTC).Unix()

	tcs := map[string]struct {
		givenPayload Claims
		method       SigningMethod
		key          crypto.Signer
		expResult    string
		expErr       error
	}{
		"success": {
			givenPayload: RegisteredClaims{
				Issuer:    "https://limitless.mukagen.com",
				Audience:  ClaimStrings([]string{"https://resource-api.com"}),
				IssuedAt:  &iat,
				ExpiresAt: &exp,
			},
			method:    NewRS256(),
			key:       key,
			expResult: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwczovL2xpbWl0bGVzcy5tdWthZ2VuLmNvbSIsImF1ZCI6WyJodHRwczovL3Jlc291cmNlLWFwaS5jb20iXSwiaWF0IjoxNzIxNzc5MjAwLCJleHAiOjE3MjE3ODI4MDB9.dnmh7SMxs3QkEJ5pUjjYeZ9CEm8x2XWNF0_YuDO46cbm3KIB5wrJBUJ5BuRVrkOAW39ZDxaamgLJunciEPT9BR7j0dGPWguZ4EZausHZt1ehn8OGZqyd6Uj16xOqkRWndsU7kaMJuKlLUFizHG305xcrh9M5PAmJe4PMZxE84SWYcj0QpCZ58zpWXA-OWTJzVSbpkbNMfl6RKsKxQ9DPDEl8wT6JorLC18Ov_MeV6KCSuwr4f15zaSLKlUv6I5n00PtD9Uw7hG9vbTDe0LHhG9WRbtVe-8Mqyz_pPV-oMMF28B1bRY48ItBKrJphSpE8QSsIGBHXdI5Es8hYYJyKOA",
		},
	}

	for scenario, tc := range tcs {
		t.Run(scenario, func(t *testing.T) {
			// Given
			tk := NewToken(tc.method, tc.givenPayload)

			// When
			rs, err := tk.SignedString(key)

			// Then
			if tc.expErr != nil {
				require.Equal(t, tc.expErr, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expResult, rs)
			}
		})
	}
}

func BenchmarkToken_signingString(b *testing.B) {
	iat := time.Date(2024, time.July, 24, 0, 0, 0, 0, time.UTC).Unix()

	tk := NewToken(NewRS256(), RegisteredClaims{
		Subject:  "limitless",
		IssuedAt: &iat,
	})

	b.Run("BenchmarkToken_SigningString", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = tk.signingString()
		}
	})
}

func TestNewToken(t *testing.T) {
	claims := RegisteredClaims{
		Subject: "test-subject",
		Issuer:  "test-issuer",
	}

	tcs := map[string]struct {
		method SigningMethod
		claims Claims
		expAlg string
		expTyp string
	}{
		"HS256": {
			method: NewHS256(),
			claims: claims,
			expAlg: "HS256",
			expTyp: "JWT",
		},
		"RS256": {
			method: NewRS256(),
			claims: claims,
			expAlg: "RS256",
			expTyp: "JWT",
		},
	}

	for scenario, tc := range tcs {
		t.Run(scenario, func(t *testing.T) {
			// When
			token := NewToken(tc.method, tc.claims)

			// Then
			require.Equal(t, tc.expAlg, token.Header["alg"])
			require.Equal(t, tc.expTyp, token.Header["typ"])
			require.Equal(t, tc.claims, token.Claims)
			require.Equal(t, tc.method, token.signingMethod)
		})
	}
}

func TestToken_signingString(t *testing.T) {
	iat := time.Date(2024, time.July, 24, 0, 0, 0, 0, time.UTC).Unix()
	claims := RegisteredClaims{
		Subject:  "test-subject",
		IssuedAt: &iat,
	}

	token := NewToken(NewHS256(), claims)

	t.Run("success", func(t *testing.T) {
		// When
		signingString, err := token.signingString()

		// Then
		require.NoError(t, err)
		require.Contains(t, string(signingString), ".")
		parts := string(signingString)
		require.Contains(t, parts, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9") // base64 of header
	})
}

func Test_encodeSegment(t *testing.T) {
	tcs := map[string]struct {
		input    []byte
		expected string
	}{
		"simple string": {
			input:    []byte("hello world"),
			expected: "aGVsbG8gd29ybGQ",
		},
		"empty": {
			input:    []byte(""),
			expected: "",
		},
		"special chars": {
			input:    []byte("test+data/with=special?chars"),
			expected: "dGVzdCtkYXRhL3dpdGg9c3BlY2lhbD9jaGFycw",
		},
	}

	for scenario, tc := range tcs {
		t.Run(scenario, func(t *testing.T) {
			// When
			result := encodeSegment(tc.input)

			// Then
			require.Equal(t, tc.expected, string(result))
		})
	}
}

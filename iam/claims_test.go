package iam

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/viebiz/lit/jwt"
)

func TestClaims_MarshalJSON(t *testing.T) {
	iat := time.Date(2024, time.July, 24, 0, 0, 0, 0, time.UTC).Unix()
	nbf := time.Date(2024, time.July, 24, 0, 0, 0, 0, time.UTC).Unix()
	exp := time.Date(2024, time.July, 24, 1, 0, 0, 0, time.UTC).Unix()

	tcs := map[string]struct {
		in        Claims
		expResult string
		expErr    error
	}{
		"success": {
			in: Claims{
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "https://limitless.mukagen.com",
					Audience:  []string{"https://resource-api.com"},
					Subject:   "mukagen|USER-ID",
					IssuedAt:  &iat,
					ExpiresAt: &exp,
					NotBefore: &nbf,
					JTI:       "JWTID",
					ClientID:  "CLIENT-UUID",
				},
				ExtraClaims: map[string]interface{}{
					"scope": "openid profile reademail",
					"https://resource.api": struct {
						PreferredContact string `json:"preferred_contact"`
					}{
						PreferredContact: "phone",
					},
				},
			},
			expResult: "{\"aud\":[\"https://resource-api.com\"],\"client_id\":\"CLIENT-UUID\",\"exp\":1721782800,\"https://resource.api\":{\"preferred_contact\":\"phone\"},\"iat\":1721779200,\"iss\":\"https://limitless.mukagen.com\",\"jti\":\"JWTID\",\"nbf\":1721779200,\"scope\":\"openid profile reademail\",\"sub\":\"mukagen|USER-ID\"}",
		},
	}

	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()
			// Given

			// When
			result, err := json.Marshal(&tc.in)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expResult, string(result))
			}
		})
	}
}

func TestClaims_UnmarshalJSON(t *testing.T) {
	iat := time.Date(2024, time.July, 24, 0, 0, 0, 0, time.UTC).Unix()
	nbf := time.Date(2024, time.July, 24, 0, 0, 0, 0, time.UTC).Unix()
	exp := time.Date(2024, time.July, 24, 1, 0, 0, 0, time.UTC).Unix()

	tcs := map[string]struct {
		in        string
		expResult Claims
		expErr    error
	}{
		"success": {
			in: "{\"aud\":\"https://resource-api.com\",\"client_id\":\"CLIENT-UUID\",\"exp\":1721782800,\"https://resource.api\":{\"preferred_contact\":\"phone\"},\"iat\":1721779200,\"iss\":\"https://limitless.mukagen.com\",\"jti\":\"JWTID\",\"nbf\":1721779200,\"scope\":\"openid profile reademail\",\"sub\":\"mukagen|USER-ID\"}",
			expResult: Claims{
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "https://limitless.mukagen.com",
					Audience:  []string{"https://resource-api.com"},
					Subject:   "mukagen|USER-ID",
					IssuedAt:  &iat,
					ExpiresAt: &exp,
					NotBefore: &nbf,
					JTI:       "JWTID",
					ClientID:  "CLIENT-UUID",
				},
				ExtraClaims: map[string]interface{}{
					"scope": "openid profile reademail",
					"https://resource.api": map[string]interface{}{
						"preferred_contact": "phone",
					},
				},
			},
		},
		"success - empty claims": {
			in:        "{}",
			expResult: Claims{},
		},
	}

	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()
			// Given

			// When
			var result Claims
			err := json.Unmarshal([]byte(tc.in), &result)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expResult, result)
			}
		})
	}
}

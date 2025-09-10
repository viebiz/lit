package jwt

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestClaimStrings_MarshalJSON(t *testing.T) {
	tcs := map[string]struct {
		input    ClaimStrings
		expected string
	}{
		"single string": {
			input:    ClaimStrings{"single-audience"},
			expected: `["single-audience"]`,
		},
		"multiple strings": {
			input:    ClaimStrings{"aud1", "aud2", "aud3"},
			expected: `["aud1","aud2","aud3"]`,
		},
		"empty slice": {
			input:    ClaimStrings{},
			expected: `[]`,
		},
	}

	for scenario, tc := range tcs {
		t.Run(scenario, func(t *testing.T) {
			// When
			result, err := json.Marshal(tc.input)

			// Then
			require.NoError(t, err)
			require.Equal(t, tc.expected, string(result))
		})
	}
}

func TestClaimStrings_UnmarshalJSON(t *testing.T) {
	tcs := map[string]struct {
		input    string
		expected ClaimStrings
		expErr   bool
	}{
		"single string": {
			input:    `"single-audience"`,
			expected: ClaimStrings{"single-audience"},
			expErr:   false,
		},
		"array of strings": {
			input:    `["aud1","aud2"]`,
			expected: ClaimStrings{"aud1", "aud2"},
			expErr:   false,
		},
		"empty string": {
			input:    `""`,
			expected: nil,
			expErr:   false,
		},
		"empty array": {
			input:    `[]`,
			expected: ClaimStrings{},
			expErr:   false,
		},
		"array of interfaces": {
			input:    `["string",123,true]`,
			expected: ClaimStrings{"string", "123", "true"},
			expErr:   false,
		},
		"invalid json": {
			input:  `invalid`,
			expErr: true,
		},
		"invalid type": {
			input:  `{"invalid":"type"}`,
			expErr: true,
		},
	}

	for scenario, tc := range tcs {
		t.Run(scenario, func(t *testing.T) {
			// Given
			var result ClaimStrings

			// When
			err := json.Unmarshal([]byte(tc.input), &result)

			// Then
			if tc.expErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestRegisteredClaims_Valid(t *testing.T) {
	now := time.Date(2024, time.July, 24, 12, 0, 0, 0, time.UTC).Unix()
	past := time.Date(2024, time.July, 24, 11, 0, 0, 0, time.UTC).Unix()
	future := time.Date(2024, time.July, 24, 13, 0, 0, 0, time.UTC).Unix()

	defer func(origin func() time.Time) { timeNowFunc = origin }(timeNowFunc)
	timeNowFunc = func() time.Time { return time.Unix(now, 0) }

	tcs := map[string]struct {
		claims   RegisteredClaims
		expError error
	}{
		"valid claims": {
			claims: RegisteredClaims{
				ExpiresAt: &future,
			},
			expError: nil,
		},
		"valid with iat": {
			claims: RegisteredClaims{
				IssuedAt:  &past,
				ExpiresAt: &future,
			},
			expError: nil,
		},
		"valid with nbf": {
			claims: RegisteredClaims{
				NotBefore: &past,
				ExpiresAt: &future,
			},
			expError: nil,
		},
		"missing exp": {
			claims: RegisteredClaims{
				Subject: "test",
			},
			expError: ErrMissingRequiredClaim,
		},
		"expired token": {
			claims: RegisteredClaims{
				ExpiresAt: &past,
			},
			expError: ErrTokenExpired,
		},
		"token used before issued": {
			claims: RegisteredClaims{
				IssuedAt:  &future,
				ExpiresAt: &future,
			},
			expError: ErrTokenUsedBeforeIssued,
		},
		"token not valid yet": {
			claims: RegisteredClaims{
				NotBefore: &future,
				ExpiresAt: &future,
			},
			expError: ErrTokenNotValidYet,
		},
	}

	for scenario, tc := range tcs {
		t.Run(scenario, func(t *testing.T) {
			// When
			err := tc.claims.Valid()

			// Then
			if tc.expError != nil {
				require.Equal(t, tc.expError, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

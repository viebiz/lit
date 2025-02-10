package iam

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/viebiz/lit/jwt"
)

func TestRFC9068Validator_Validate(t *testing.T) {
	const (
		issuer          string = "https://mukagen.com"
		audience        string = "https://limitless.mukagen.com"
		privateKeyPath  string = "testdata/sample_rsa_private_key"
		certificatePath string = "testdata/sample_rsa_certificate"
	)
	staticTimeNow := time.Date(2024, time.July, 24, 0, 0, 0, 0, time.UTC)

	tcs := map[string]struct {
		givenToken string
		expResult  jwt.Token[Claims]
		expErr     error
	}{
		"success": {
			givenToken: "eyJhbGciOiJSUzI1NiIsImtpZCI6Impzb24td2ViLWtleS0wMSIsInR5cCI6ImF0K2p3dCJ9.eyJhdWQiOlsiaHR0cHM6Ly9saW1pdGxlc3MubXVrYWdlbi5jb20iXSwiY2xpZW50X2lkIjoianVzby1zaGkiLCJleHAiOjE3MjE3ODI4MDAsImlhdCI6MTcyMTc3OTIwMCwiaXNzIjoiaHR0cHM6Ly9tdWthZ2VuLmNvbS8iLCJqdGkiOiJ0aGlzLWlzLXV1aWQiLCJzY29wZXMiOiJza2lsbHM6cmVhZCIsInN1YiI6Imp1c28tc2hpQGNsaWVudHMifQ.JsPQrWhBq7IuiEc-skI4I1f0KCvbUcq3JROLkw8dAEe6lN375TvKR0xYHdhUcjYpbL1JtNFL0Bv2Bh13JV9twOWxIp3mrRl_89L8_ThFM8GaJL5x1YSpBFqr-YkXf6HwpXzM4n7B6x4gmAVR7uDMw5zJBLTaCGiVPDIL_LjgiyZS7AQ7HNnj5D-egHDNXjQIqWBriWDALuoX2tHS1qubQYNdggm-RrUKLiRQdsIh5bohb5T0-11IIRYmcAfL9U17o6NnzMGC3BM9e4xSruJa63F3y_whmwiOtbo6EDeTPsxyBzADt26s3ToZEcYY2nV7iz3lgXbuk61ZRYdV-4yUtg",
			expResult: jwt.Token[Claims]{
				Header: map[string]string{
					"typ": "at+jwt",
					"alg": "RS256",
					"kid": "json-web-key-01",
				},
				Claims: Claims{
					RegisteredClaims: jwt.RegisteredClaims{
						Issuer:    "https://mukagen.com/",
						Subject:   "juso-shi@clients",
						Audience:  []string{"https://limitless.mukagen.com"},
						IssuedAt:  pointerTo(staticTimeNow.Unix()),
						ExpiresAt: pointerTo(staticTimeNow.Add(time.Hour).Unix()),
						ClientID:  "juso-shi",
						JTI:       "this-is-uuid",
					},
					ExtraClaims: map[string]interface{}{
						"scopes": "skills:read",
					},
				},
			},
		},
		"error - invalid signature": {
			givenToken: "eyJhbGciOiJSUzI1NiIsImtpZCI6Impzb24td2ViLWtleS0wMSIsInR5cCI6ImF0K2p3dCJ9.eyJhdWQiOlsiaHR0cHM6Ly9saW1pdGxlc3MubXVrYWdlbi5jb20iXSwiY2xpZW50X2lkIjoianVzby1zaGkiLCJleHAiOjE3MjE3ODI4MDAsImlhdCI6MTcyMTc3OTIwMCwiaXNzIjoiaHR0cHM6Ly9tdWthZ2VuLmNvbS8iLCJqdGkiOiJ0aGlzLWlzLXV1aWQiLCJzY29wZXMiOiJza2lsbHM6cmVhZCIsInN1YiI6Imp1c28tc2hpQGNsaWVudHMifQ.JsPQrWhBq7IuiEc-skI4I1f0KCvbUcq3JROLkw8dAEe6lN375TvKR0xYHdhUcjYpbL1JtNFL0Bv2Bh13JV9twOWxIp3mrRl_89L8_ThFM8GaJL5x1YSpBFqr-YkXf6HwpXzM4n7B6x4gmAVR7uDMw5zJBLTaCGiVPDIL_LjgiyZS7AQ7HNnj5D-egHDNXjQIqWBriWDALuoX2tHS1qubQYNdggm-RrUKLiRQdsIh5bohb5T0-11IIRYmcAfL9U17o6NnzMGC3BM9e4xSruJa63F3y_whmwiOtbo6EDeTPsxyBzADt26s3ToZEcYY2nV7iz3lgXbuk61ZRYdinvalid",
			expErr:     jwt.ErrInvalidSignature,
		},
		"error - token malformed": {
			givenToken: "",
			expErr:     jwt.ErrTokenMalformed,
		},
		"error - invalid token when header not have typ": {
			givenToken: "eyJhbGciOiJSUzI1NiIsImtpZCI6Impzb24td2ViLWtleS0wMSJ9.eyJhdWQiOlsiaHR0cHM6Ly9saW1pdGxlc3MubXVrYWdlbi5jb20iXSwiY2xpZW50X2lkIjoianVzby1zaGkiLCJleHAiOjE3MjE2OTY0MDAsImlhdCI6MTcyMTY5MjgwMCwiaXNzIjoiaHR0cHM6Ly9tdWthZ2VuLmNvbS8iLCJqdGkiOiJ0aGlzLWlzLXV1aWQiLCJzY29wZXMiOiJza2lsbHM6cmVhZCIsInN1YiI6Imp1c28tc2hpQGNsaWVudHMifQ.HDvoTY7-0Tffm0mIvug8aX557K4PeS03SUsv545ODF-5IZxWLN3hK8hVHceTmo2bpYN-xN5mHOfaezOqlvZtVU85Ln3JWhahtqul7sHoppEsJhgqnRnivI67ZJGdh4x5DVZm86Vp2qqUwubTGImOA-HotqKYtV9ZILv65ySf2jm0lC6z0qt2-TbKjg32ITBbS4W5HPXiMOulI5Yl6bgE_TA_RxFj_f6wjwiaVL3aycWnnoZmqdJ_PTVG1RlUVavjqH-iQZg2NfqbHV6kRgb66OU15ZHQGztr5KEwPOcCD-jO1mzUHewCdD47cDS3Wb6OXesPOY09bwyiMnrpyO-GKA",
			expErr:     ErrInvalidToken,
		},
		"error - expires token": {
			givenToken: "eyJhbGciOiJSUzI1NiIsImtpZCI6Impzb24td2ViLWtleS0wMSIsInR5cCI6ImF0K2p3dCJ9.eyJhdWQiOlsiaHR0cHM6Ly9saW1pdGxlc3MubXVrYWdlbi5jb20iXSwiY2xpZW50X2lkIjoianVzby1zaGkiLCJleHAiOjE3MjE2OTY0MDAsImlhdCI6MTcyMTY5MjgwMCwiaXNzIjoiaHR0cHM6Ly9tdWthZ2VuLmNvbS8iLCJqdGkiOiJ0aGlzLWlzLXV1aWQiLCJzY29wZXMiOiJza2lsbHM6cmVhZCIsInN1YiI6Imp1c28tc2hpQGNsaWVudHMifQ.Q5wJsbFSPS_YZ42xOqwo7qjN2eVHslwbg_zYRdjrfOIRhqHwwN-jKfF7g2qf2xqTvASxgoYtCf7tSQ4ffVF8jCxnEKYGJV_anFGELLxi0KZpGpuKBeCk2-XajfI8t0GSB4dJsRhLZ3sGdKcfhozS_johnHUTXVKhboQLturETG155uXkyZ8YAWzh_ZDSCuXcg8RPGs6eUuLeBZ0V6WztES5XkOp8ZKZA5RhazRbhZ9k5stfWql2NYZQ4n2DlXpykherYAFoxag0JAFoDXgeU5ZsQysUPUsOoDSTgmKyptQLFfXgGplIPSbgMY6_yl-zgjrkzsrCtZCRJWYJZsjpGrg",
			expErr:     ErrTokenExpired,
		},
	}

	for scenario, tc := range tcs {
		t.Run(scenario, func(t *testing.T) {
			// Given
			defer func(origin func() time.Time) { timeNowFunc = origin }(timeNowFunc)
			timeNowFunc = func() time.Time { return staticTimeNow }

			pkey := readRSAPrivateKey(t, privateKeyPath)
			cert := readCertificate(t, certificatePath)
			jwkSet := constructJWKSForTest(pkey.PublicKey, *cert)
			jwkJSON, err := json.Marshal(jwkSet)
			require.NoError(t, err)

			mockClient := new(mockHTTPClient)
			mockClient.On("Do", mock.Anything).Return(&http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(jwkJSON)),
			}, nil)

			validator, err := NewRFC9068Validator(issuer, audience, mockClient)
			require.NoError(t, err)

			// When
			tk, err := validator.Validate(tc.givenToken)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
				ignoreFieldsOpts := cmpopts.IgnoreFields(jwt.Token[Claims]{}, "Signature", "signingMethod")

				if !cmp.Equal(tc.expResult, tk, ignoreFieldsOpts) {
					t.Errorf("\n result mismatched. Diff: %+v", cmp.Diff(tc.expResult, tk, ignoreFieldsOpts))
					t.FailNow()
				}
			}
			mockClient.AssertExpectations(t)
		})
	}
}

func pointerTo[T any](v T) *T {
	return &v
}

// ==========================================
// ENABLE IT IF YOU NEED GENERATE SOME JWT FOR TEST
// ==========================================
//func Test_help(t *testing.T) {
//	const (
//		privateKeyPath = "testdata/sample_rsa_private_key"
//		certPath       = "testdata/sample_rsa_certificate"
//	)
//
//	iat := time.Date(2024, time.December, 4, 0, 0, 0, 0, time.UTC).Unix()
//	exp := time.Date(2024, time.December, 4, 0, 30, 0, 0, time.UTC).Unix()
//
//	pKey := readRSAPrivateKey(t, privateKeyPath)
//	cert := readCertificate(t, certPath)
//
//	jwks := constructJWKSForTest(pKey.PublicKey, *cert)
//	jwskStr, _ := json.Marshal(jwks)
//
//	claims := &Claims{
//		RegisteredClaims: jwt.RegisteredClaims{
//			Issuer:    "https://warhammer40k.imperium.io",
//			Subject:   "imperium|space_marine",
//			Audience:  []string{"https://space-marine.com"},
//			IssuedAt:  &iat,
//			ExpiresAt: &exp,
//			ClientID:  "space_marine",
//			JTI:       "GalacticPhantom",
//		},
//	}
//
//	tk := jwt.NewToken(jwt.NewRS256(), claims)
//	tk.Header["typ"] = "at+jwt"
//	tk.Header["kid"] = "KID-001"
//	str, err := tk.SignedString(pKey)
//	require.NoError(t, err)
//
//	fmt.Printf("Public JWKS: %s\n", string(jwskStr))
//	fmt.Printf("JWT: %s\n", str)
//
//	parser := jwt.NewDefaultParser[*Claims]()
//	tk, err = parser.Parse(str, func(s string) (crypto.PublicKey, error) {
//		x5c := jwks.Keys[0].X5c[0]
//		pemBlock, _ := pem.Decode([]byte("-----BEGIN CERTIFICATE-----\n" + x5c + "\n-----END CERTIFICATE-----"))
//		if pemBlock == nil {
//			return nil, fmt.Errorf("could not parse certificate PEM")
//		}
//
//		cr, err := x509.ParseCertificate(pemBlock.Bytes)
//		if err != nil {
//			return nil, err
//		}
//
//		return cr.PublicKey, nil
//	})
//
//	require.NoError(t, err)
//	fmt.Println(tk)
//}

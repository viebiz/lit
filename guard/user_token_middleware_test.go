package guard

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/viebiz/lit"
	"github.com/viebiz/lit/iam"
	"github.com/viebiz/lit/jwt"
)

func TestAuthenticateUserMiddleware(t *testing.T) {
	staticTimeNow := time.Date(2024, time.December, 1, 0, 0, 0, 0, time.UTC)

	type mockData struct {
		expCall  bool
		inputStr string
		outToken jwt.Token[iam.Claims]
		outErr   error
	}
	tcs := map[string]struct {
		givenToken string
		mockData   mockData
		expResult  iam.UserProfile
		expErr     error
	}{
		"success": {
			givenToken: "eyJhbGciOiJSUzI1NiIsImtpZCI6IktJRC0wMDEiLCJ0eXAiOiJhdCtqd3QifQ.eyJhdWQiOlsiaHR0cHM6Ly9zcGFjZS1tYXJpbmUuY29tIl0sImNsaWVudF9pZCI6InNwYWNlX21hcmluZSIsImV4cCI6MTczMzAxMzAwMCwiaWF0IjoxNzMzMDExMjAwLCJpc3MiOiJodHRwczovL3dhcmhhbW1lcjQway5pbXBlcml1bS5pbyIsImp0aSI6IkdhbGFjdGljUGhhbnRvbSIsInJvbGVzIjpbInByaW1hcmNoIl0sInN1YiI6ImltcGVyaXVtfHNwYWNlX21hcmluZSJ9.QVTKaGIUhRSWRDiiggHToKiy3CaNSsd4CVfUiqA8IZetTALuQaHCaQcXjeXPBGj9vUxiI0pq5oLvQHhfVxUj9hwQ5PkHqRw-GdqdaJVj8q3wjX3Ulxk6J1p8TSddWi3TUNSktMhVVny5Fu0uWoDTmm3oxiKDf66aNWMaGRlQLb4cP9YuglqC5HiR7mrYE5gG95fNHO9fWnA2Ao6FblCdOvfYZNDOo771bJJWSpJDZEpwyDn_h1jtwcLq8vNHRu9Yga-B416tWCMgV6kIUzUGvs3QD2dtX1MLuQiLB9h0_ZBqQa4r6baSjhGKUFWF6m1ioY8j8rqDjVMmPkgLP3W0QQ",
			mockData: mockData{
				expCall:  true,
				inputStr: "eyJhbGciOiJSUzI1NiIsImtpZCI6IktJRC0wMDEiLCJ0eXAiOiJhdCtqd3QifQ.eyJhdWQiOlsiaHR0cHM6Ly9zcGFjZS1tYXJpbmUuY29tIl0sImNsaWVudF9pZCI6InNwYWNlX21hcmluZSIsImV4cCI6MTczMzAxMzAwMCwiaWF0IjoxNzMzMDExMjAwLCJpc3MiOiJodHRwczovL3dhcmhhbW1lcjQway5pbXBlcml1bS5pbyIsImp0aSI6IkdhbGFjdGljUGhhbnRvbSIsInJvbGVzIjpbInByaW1hcmNoIl0sInN1YiI6ImltcGVyaXVtfHNwYWNlX21hcmluZSJ9.QVTKaGIUhRSWRDiiggHToKiy3CaNSsd4CVfUiqA8IZetTALuQaHCaQcXjeXPBGj9vUxiI0pq5oLvQHhfVxUj9hwQ5PkHqRw-GdqdaJVj8q3wjX3Ulxk6J1p8TSddWi3TUNSktMhVVny5Fu0uWoDTmm3oxiKDf66aNWMaGRlQLb4cP9YuglqC5HiR7mrYE5gG95fNHO9fWnA2Ao6FblCdOvfYZNDOo771bJJWSpJDZEpwyDn_h1jtwcLq8vNHRu9Yga-B416tWCMgV6kIUzUGvs3QD2dtX1MLuQiLB9h0_ZBqQa4r6baSjhGKUFWF6m1ioY8j8rqDjVMmPkgLP3W0QQ",
				outToken: jwt.Token[iam.Claims]{
					Header: map[string]string{
						"typ": "at+jwt",
						"alg": "RS256",
						"kid": "KID-001",
					},
					Claims: iam.Claims{
						RegisteredClaims: jwt.RegisteredClaims{
							Issuer:    "https://warhammer40k.imperium.io",
							Subject:   "imperium|space_marine",
							Audience:  []string{"https://space-marine.com"},
							IssuedAt:  pointerTo(staticTimeNow.Unix()),
							ExpiresAt: pointerTo(staticTimeNow.Add(30 * time.Minute).Unix()),
							ClientID:  "space_marine",
							JTI:       "GalacticPhantom",
						},
						ExtraClaims: map[string]interface{}{
							"roles": []string{"primarch"},
						},
					},
				},
			},
			expResult: iam.NewUserProfile("imperium|space_marine", []string{"primarch"}, nil),
		},
		"error - missing access token": {
			givenToken: "",
			expErr:     errMissingAccessToken,
		},
		"error - token expires": {
			givenToken: "eyJhbGciOiJSUzI1NiIsImtpZCI6IktJRC0wMDEiLCJ0eXAiOiJhdCtqd3QifQ.eyJhdWQiOlsiaHR0cHM6Ly9zcGFjZS1tYXJpbmUuY29tIl0sImNsaWVudF9pZCI6InNwYWNlX21hcmluZSIsImV4cCI6MTczMzI3MjIwMCwiaWF0IjoxNzMzMjcwNDAwLCJpc3MiOiJodHRwczovL3dhcmhhbW1lcjQway5pbXBlcml1bS5pbyIsImp0aSI6IkdhbGFjdGljUGhhbnRvbSIsInJvbGVzIjpbInByaW1hcmNoIl0sInN1YiI6ImltcGVyaXVtfHNwYWNlX21hcmluZSJ9.F81jB-JGpHwLrQ4Wy96rkQx4GslpbFmmn1w-eFuuOrsF1yBOE6FASdI3nxbWQFLjSTOxDUzyxhKLYF8BP7ImHHDabtWvcWBW8Gw3GNzByPODxzSqw1ceYWc0SmG3EmrsqDFJZ9Rnvy6KmPLtVZhf8z4u4dyyCkBgnVxre7L9yeCoQe5zFQayw-hlhDzXYGCKFdmhILRcM0VsDycM-YbWIP7E0GBv-8dtPoLMXBVDjz8ln7eMwCbn69NTjWmW2BmEuRF5QkgnVmmnv9Ldw9iQcgg8N9NruARNCfXvnT9NTxkhsBXLR1ACNfJo3r1XcYMF-1IODVHxwmVZ2503nIPMZw",
			mockData: mockData{
				expCall:  true,
				inputStr: "eyJhbGciOiJSUzI1NiIsImtpZCI6IktJRC0wMDEiLCJ0eXAiOiJhdCtqd3QifQ.eyJhdWQiOlsiaHR0cHM6Ly9zcGFjZS1tYXJpbmUuY29tIl0sImNsaWVudF9pZCI6InNwYWNlX21hcmluZSIsImV4cCI6MTczMzI3MjIwMCwiaWF0IjoxNzMzMjcwNDAwLCJpc3MiOiJodHRwczovL3dhcmhhbW1lcjQway5pbXBlcml1bS5pbyIsImp0aSI6IkdhbGFjdGljUGhhbnRvbSIsInJvbGVzIjpbInByaW1hcmNoIl0sInN1YiI6ImltcGVyaXVtfHNwYWNlX21hcmluZSJ9.F81jB-JGpHwLrQ4Wy96rkQx4GslpbFmmn1w-eFuuOrsF1yBOE6FASdI3nxbWQFLjSTOxDUzyxhKLYF8BP7ImHHDabtWvcWBW8Gw3GNzByPODxzSqw1ceYWc0SmG3EmrsqDFJZ9Rnvy6KmPLtVZhf8z4u4dyyCkBgnVxre7L9yeCoQe5zFQayw-hlhDzXYGCKFdmhILRcM0VsDycM-YbWIP7E0GBv-8dtPoLMXBVDjz8ln7eMwCbn69NTjWmW2BmEuRF5QkgnVmmnv9Ldw9iQcgg8N9NruARNCfXvnT9NTxkhsBXLR1ACNfJo3r1XcYMF-1IODVHxwmVZ2503nIPMZw",
				outErr:   iam.ErrTokenExpired,
			},
			expErr: unauthorizedErr(iam.ErrTokenExpired),
		},
		"error - missing role claim": {
			givenToken: "eyJhbGciOiJSUzI1NiIsImtpZCI6IktJRC0wMDEiLCJ0eXAiOiJhdCtqd3QifQ.eyJhdWQiOlsiaHR0cHM6Ly9zcGFjZS1tYXJpbmUuY29tIl0sImNsaWVudF9pZCI6InNwYWNlX21hcmluZSIsImV4cCI6MTczMzAxMzAwMCwiaWF0IjoxNzMzMDExMjAwLCJpc3MiOiJodHRwczovL3dhcmhhbW1lcjQway5pbXBlcml1bS5pbyIsImp0aSI6IkdhbGFjdGljUGhhbnRvbSIsInJvbGVzIjpbInByaW1hcmNoIl0sInN1YiI6ImltcGVyaXVtfHNwYWNlX21hcmluZSJ9.QVTKaGIUhRSWRDiiggHToKiy3CaNSsd4CVfUiqA8IZetTALuQaHCaQcXjeXPBGj9vUxiI0pq5oLvQHhfVxUj9hwQ5PkHqRw-GdqdaJVj8q3wjX3Ulxk6J1p8TSddWi3TUNSktMhVVny5Fu0uWoDTmm3oxiKDf66aNWMaGRlQLb4cP9YuglqC5HiR7mrYE5gG95fNHO9fWnA2Ao6FblCdOvfYZNDOo771bJJWSpJDZEpwyDn_h1jtwcLq8vNHRu9Yga-B416tWCMgV6kIUzUGvs3QD2dtX1MLuQiLB9h0_ZBqQa4r6baSjhGKUFWF6m1ioY8j8rqDjVMmPkgLP3W0QQ",
			mockData: mockData{
				expCall:  true,
				inputStr: "eyJhbGciOiJSUzI1NiIsImtpZCI6IktJRC0wMDEiLCJ0eXAiOiJhdCtqd3QifQ.eyJhdWQiOlsiaHR0cHM6Ly9zcGFjZS1tYXJpbmUuY29tIl0sImNsaWVudF9pZCI6InNwYWNlX21hcmluZSIsImV4cCI6MTczMzAxMzAwMCwiaWF0IjoxNzMzMDExMjAwLCJpc3MiOiJodHRwczovL3dhcmhhbW1lcjQway5pbXBlcml1bS5pbyIsImp0aSI6IkdhbGFjdGljUGhhbnRvbSIsInJvbGVzIjpbInByaW1hcmNoIl0sInN1YiI6ImltcGVyaXVtfHNwYWNlX21hcmluZSJ9.QVTKaGIUhRSWRDiiggHToKiy3CaNSsd4CVfUiqA8IZetTALuQaHCaQcXjeXPBGj9vUxiI0pq5oLvQHhfVxUj9hwQ5PkHqRw-GdqdaJVj8q3wjX3Ulxk6J1p8TSddWi3TUNSktMhVVny5Fu0uWoDTmm3oxiKDf66aNWMaGRlQLb4cP9YuglqC5HiR7mrYE5gG95fNHO9fWnA2Ao6FblCdOvfYZNDOo771bJJWSpJDZEpwyDn_h1jtwcLq8vNHRu9Yga-B416tWCMgV6kIUzUGvs3QD2dtX1MLuQiLB9h0_ZBqQa4r6baSjhGKUFWF6m1ioY8j8rqDjVMmPkgLP3W0QQ",
				outToken: jwt.Token[iam.Claims]{
					Header: map[string]string{
						"typ": "at+jwt",
						"alg": "RS256",
						"kid": "KID-001",
					},
					Claims: iam.Claims{
						RegisteredClaims: jwt.RegisteredClaims{
							Issuer:    "https://warhammer40k.imperium.io",
							Subject:   "imperium|space_marine",
							Audience:  []string{"https://space-marine.com"},
							IssuedAt:  pointerTo(staticTimeNow.Unix()),
							ExpiresAt: pointerTo(staticTimeNow.Add(time.Hour).Unix()),
							ClientID:  "space_marine",
							JTI:       "GalacticPhantom",
						},
					},
				},
			},
			expErr: unauthorizedErr(iam.ErrMissingRequiredClaim),
		},
		"error - invalid token": {
			givenToken: "invalid-token",
			mockData: mockData{
				expCall:  true,
				inputStr: "invalid-token",
				outErr:   iam.ErrInvalidToken,
			},
			expErr: unauthorizedErr(iam.ErrInvalidToken),
		},
		"error - internal server error": {
			givenToken: "invalid-token",
			mockData: mockData{
				expCall:  true,
				inputStr: "invalid-token",
				outErr:   errors.New("simulate server error"),
			},
			expErr: lightning.ErrInternalServerError,
		},
	}

	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			// Given
			request := httptest.NewRequest(http.MethodGet, "/", nil)
			request.Header.Set(headerAuthorization, fmt.Sprintf("Bearer %s", tc.givenToken))

			respRecord := httptest.NewRecorder()

			_, ctx, _ := lightning.NewRouterForTest(respRecord)
			ctx.SetRequest(request)

			mockInstance := new(iam.MockValidator)
			if tc.mockData.expCall {
				mockInstance.On("Validate", tc.mockData.inputStr).
					Return(tc.mockData.outToken, tc.mockData.outErr)
			}

			guard := New(mockInstance, nil)

			// When
			hdl := guard.AuthenticateUserMiddleware()
			hdl(ctx)
			actualProfile := iam.GetUserProfileFromContext(ctx.Request().Context())

			// Then
			if tc.expErr != nil {
				var iamErr lightning.HttpError
				if errors.As(tc.expErr, &iamErr) {
					require.Equal(t, respRecord.Code, iamErr.Status)
				} else {
					require.Equal(t, respRecord.Code, http.StatusInternalServerError)
				}

				require.Equal(t, tc.expErr.Error(), respRecord.Body.String())
			} else {
				require.Equal(t, tc.expResult, actualProfile)
			}
			mockInstance.AssertExpectations(t)
		})
	}
}

func pointerTo[T any](value T) *T {
	return &value
}

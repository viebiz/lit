package guard

import (
	"encoding/json"
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

func TestAuthenticateM2MMiddleware(t *testing.T) {
	staticTimeNow := time.Date(2024, time.July, 24, 0, 0, 0, 0, time.UTC)

	type mockData struct {
		expCall  bool
		inputStr string
		outToken jwt.Token[iam.Claims]
		outErr   error
	}
	tcs := map[string]struct {
		givenToken string
		mockData   mockData
		expResult  iam.M2MProfile
		expErr     error
	}{
		"success": {
			givenToken: "eyJhbGciOiJSUzI1NiIsImtpZCI6IktJRC0wMDEiLCJ0eXAiOiJhdCtqd3QifQ.eyJhdWQiOlsiaHR0cHM6Ly9zcGFjZS1tYXJpbmUuY29tIl0sImNsaWVudF9pZCI6InNwYWNlX21hcmluZSIsImV4cCI6MTczMzAxMzAwMCwiaWF0IjoxNzMzMDExMjAwLCJpc3MiOiJodHRwczovL3dhcmhhbW1lcjQway5pbXBlcml1bS5pbyIsImp0aSI6IkdhbGFjdGljUGhhbnRvbSIsInNjb3BlIjoic3F1YWQgYXJtb3J5IHdlYXBvbnJ5Iiwic3ViIjoiaW1wZXJpdW18c3BhY2VfbWFyaW5lIn0.YLYdl8L0kGlq9eeaLM6rYCmEKXzB7earv99TCjDJZH2aVcG8A7kyPCz_UxffLUYLVVqk3BsDCHe1pV9WdWDpiXRmNaKAzLZxqeOpa5QhGdm0X00owjx-oydc4BCTvcikZVESMSysIvViAOs2a1ADOqyrG7clsaiGd1EUguphNIs6F-Xsx3Y0hEirC7nLS1mqrHNMMxat81XYpRhfYCQGmeW73IbGWmbtL0CnRo1b7PmmRk3Jlo7QzkeDLnCaHBMAPg5jR3MpRqzyYW8ptew-Jxom0cpkJnnHT-aQASYmdaU1UmFYCbFsRvtwk3VEr6Qt_QYTSj62sobfXU90y9mv0A",
			mockData: mockData{
				expCall:  true,
				inputStr: "eyJhbGciOiJSUzI1NiIsImtpZCI6IktJRC0wMDEiLCJ0eXAiOiJhdCtqd3QifQ.eyJhdWQiOlsiaHR0cHM6Ly9zcGFjZS1tYXJpbmUuY29tIl0sImNsaWVudF9pZCI6InNwYWNlX21hcmluZSIsImV4cCI6MTczMzAxMzAwMCwiaWF0IjoxNzMzMDExMjAwLCJpc3MiOiJodHRwczovL3dhcmhhbW1lcjQway5pbXBlcml1bS5pbyIsImp0aSI6IkdhbGFjdGljUGhhbnRvbSIsInNjb3BlIjoic3F1YWQgYXJtb3J5IHdlYXBvbnJ5Iiwic3ViIjoiaW1wZXJpdW18c3BhY2VfbWFyaW5lIn0.YLYdl8L0kGlq9eeaLM6rYCmEKXzB7earv99TCjDJZH2aVcG8A7kyPCz_UxffLUYLVVqk3BsDCHe1pV9WdWDpiXRmNaKAzLZxqeOpa5QhGdm0X00owjx-oydc4BCTvcikZVESMSysIvViAOs2a1ADOqyrG7clsaiGd1EUguphNIs6F-Xsx3Y0hEirC7nLS1mqrHNMMxat81XYpRhfYCQGmeW73IbGWmbtL0CnRo1b7PmmRk3Jlo7QzkeDLnCaHBMAPg5jR3MpRqzyYW8ptew-Jxom0cpkJnnHT-aQASYmdaU1UmFYCbFsRvtwk3VEr6Qt_QYTSj62sobfXU90y9mv0A",
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
							"scope": "squad armory weaponry",
						},
					},
				},
			},
			expResult: iam.NewM2MProfile("imperium|space_marine", []string{"squad", "armory", "weaponry"}),
		},
		"error - missing access token": {
			givenToken: "",
			expErr:     errMissingAccessToken,
		},
		"error - token expires": {
			givenToken: "eyJhbGciOiJSUzI1NiIsImtpZCI6IktJRC0wMDEiLCJ0eXAiOiJhdCtqd3QifQ.eyJhdWQiOlsiaHR0cHM6Ly9zcGFjZS1tYXJpbmUuY29tIl0sImNsaWVudF9pZCI6InNwYWNlX21hcmluZSIsImV4cCI6MTczMzI3MjIwMCwiaWF0IjoxNzMzMjcwNDAwLCJpc3MiOiJodHRwczovL3dhcmhhbW1lcjQway5pbXBlcml1bS5pbyIsImp0aSI6IkdhbGFjdGljUGhhbnRvbSIsInNjb3BlIjoic3F1YWQgYXJtb3J5IHdlYXBvbnJ5Iiwic3ViIjoiaW1wZXJpdW18c3BhY2VfbWFyaW5lIn0.QVS9ecu0D4kI3wT30-zt7mfgM8XoCvcumC55PECGEgbkA5CYAwqIJluEKybdQAZvaj2HD-4gswC-wQRo89ai43J1LIL-xfcENlmmdGvqjq0mxpcKk1X3PmFy7E9KgYUuPDkPvgZE7Ib5Ly89_KroQbBlZpplcBo34HjXhUtcnCb7djIg7IcaSZ2EvFM7uxmaXbPYWBDCFrj4rc2zoQZbQ3HOfO206lqtWeEYXwIV5austly1a4-plvBUsLbo4LfFyyg9C4f-P8aTy5pgpeOW08YodkBUsyUBXnc3h0uwMx-y2qr5F95p1bCQ5RRKwp_wXtC8JMW-eriYiS0kffd8Aw",
			mockData: mockData{
				expCall:  true,
				inputStr: "eyJhbGciOiJSUzI1NiIsImtpZCI6IktJRC0wMDEiLCJ0eXAiOiJhdCtqd3QifQ.eyJhdWQiOlsiaHR0cHM6Ly9zcGFjZS1tYXJpbmUuY29tIl0sImNsaWVudF9pZCI6InNwYWNlX21hcmluZSIsImV4cCI6MTczMzI3MjIwMCwiaWF0IjoxNzMzMjcwNDAwLCJpc3MiOiJodHRwczovL3dhcmhhbW1lcjQway5pbXBlcml1bS5pbyIsImp0aSI6IkdhbGFjdGljUGhhbnRvbSIsInNjb3BlIjoic3F1YWQgYXJtb3J5IHdlYXBvbnJ5Iiwic3ViIjoiaW1wZXJpdW18c3BhY2VfbWFyaW5lIn0.QVS9ecu0D4kI3wT30-zt7mfgM8XoCvcumC55PECGEgbkA5CYAwqIJluEKybdQAZvaj2HD-4gswC-wQRo89ai43J1LIL-xfcENlmmdGvqjq0mxpcKk1X3PmFy7E9KgYUuPDkPvgZE7Ib5Ly89_KroQbBlZpplcBo34HjXhUtcnCb7djIg7IcaSZ2EvFM7uxmaXbPYWBDCFrj4rc2zoQZbQ3HOfO206lqtWeEYXwIV5austly1a4-plvBUsLbo4LfFyyg9C4f-P8aTy5pgpeOW08YodkBUsyUBXnc3h0uwMx-y2qr5F95p1bCQ5RRKwp_wXtC8JMW-eriYiS0kffd8Aw",
				outErr:   iam.ErrTokenExpired,
			},
			expErr: unauthorizedErr(iam.ErrTokenExpired),
		},
		"error - missing scope claim": {
			givenToken: "eyJhbGciOiJSUzI1NiIsImtpZCI6IktJRC0wMDEiLCJ0eXAiOiJhdCtqd3QifQ.eyJhdWQiOlsiaHR0cHM6Ly9zcGFjZS1tYXJpbmUuY29tIl0sImNsaWVudF9pZCI6InNwYWNlX21hcmluZSIsImV4cCI6MTczMzI3MjIwMCwiaWF0IjoxNzMzMjcwNDAwLCJpc3MiOiJodHRwczovL3dhcmhhbW1lcjQway5pbXBlcml1bS5pbyIsImp0aSI6IkdhbGFjdGljUGhhbnRvbSIsInN1YiI6ImltcGVyaXVtfHNwYWNlX21hcmluZSJ9.GNRUsKgxsPWSMTOoPgE2J2HB6JL342rKGBgsaZYPMlurhQEq9GcKMQXHjMSg1fbd6vnHmtbx1RtMtLLkAEeKo7G3DLgTRHyEMLYbNaYpmLNt0Y7kRqDE9MEFmJgUgpE0SV0q7EpZl7y1V9nW_V90FkQpgjMIYvBT6vSawdXaEuaNoBeReGPQej6PXja_b2iM6mSSYt6ejSgmbyX32uQY6lsHAN2pzIB01fFNidcndxhKJlmQfBKE_L6sl1xvZYccZNwrKlgiYisQxAkwSthocLyZrvk9TyU8NcqGC-dtsrHj276o1Arfvf0X0v-lqgxkZLMzJSYBbBw_MTc-_9hh5g",
			mockData: mockData{
				expCall:  true,
				inputStr: "eyJhbGciOiJSUzI1NiIsImtpZCI6IktJRC0wMDEiLCJ0eXAiOiJhdCtqd3QifQ.eyJhdWQiOlsiaHR0cHM6Ly9zcGFjZS1tYXJpbmUuY29tIl0sImNsaWVudF9pZCI6InNwYWNlX21hcmluZSIsImV4cCI6MTczMzI3MjIwMCwiaWF0IjoxNzMzMjcwNDAwLCJpc3MiOiJodHRwczovL3dhcmhhbW1lcjQway5pbXBlcml1bS5pbyIsImp0aSI6IkdhbGFjdGljUGhhbnRvbSIsInN1YiI6ImltcGVyaXVtfHNwYWNlX21hcmluZSJ9.GNRUsKgxsPWSMTOoPgE2J2HB6JL342rKGBgsaZYPMlurhQEq9GcKMQXHjMSg1fbd6vnHmtbx1RtMtLLkAEeKo7G3DLgTRHyEMLYbNaYpmLNt0Y7kRqDE9MEFmJgUgpE0SV0q7EpZl7y1V9nW_V90FkQpgjMIYvBT6vSawdXaEuaNoBeReGPQej6PXja_b2iM6mSSYt6ejSgmbyX32uQY6lsHAN2pzIB01fFNidcndxhKJlmQfBKE_L6sl1xvZYccZNwrKlgiYisQxAkwSthocLyZrvk9TyU8NcqGC-dtsrHj276o1Arfvf0X0v-lqgxkZLMzJSYBbBw_MTc-_9hh5g",
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
			expErr: lit.ErrDefaultInternal,
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

			_, ctx, _ := lit.NewRouterForTest(respRecord)
			ctx.SetRequest(request)

			mockInstance := new(iam.MockValidator)
			if tc.mockData.expCall {
				mockInstance.On("Validate", tc.mockData.inputStr).
					Return(tc.mockData.outToken, tc.mockData.outErr)
			}

			guard := New(mockInstance, nil)

			// When
			hdl := guard.AuthenticateM2MMiddleware()
			hdl(ctx)
			actualProfile := iam.GetM2MProfileFromContext(ctx.Request().Context())

			// Then
			if tc.expErr != nil {
				var iamErr lit.HttpError
				if errors.As(tc.expErr, &iamErr) {
					require.Equal(t, respRecord.Code, iamErr.Status)
				} else {
					require.Equal(t, respRecord.Code, http.StatusInternalServerError)
				}

				expResult, err := json.Marshal(tc.expErr)
				require.NoError(t, err)
				require.Equal(t, expResult, respRecord.Body.Bytes())
			} else {
				require.Equal(t, tc.expResult, actualProfile)
			}
			mockInstance.AssertExpectations(t)
		})
	}
}

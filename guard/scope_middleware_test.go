package guard

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/viebiz/lit"
	"github.com/viebiz/lit/iam"
)

func TestRequiredM2MScopeMiddleware(t *testing.T) {
	tcs := map[string]struct {
		in         []string
		m2mProfile iam.M2MProfile
		expErr     error
	}{
		"success": {
			in:         []string{"weaponry"},
			m2mProfile: iam.NewM2MProfile("imperium|ultra_marine", []string{"squad", "armory", "weaponry"}),
		},
		"error - profile not exists": {
			in:     []string{"armory"},
			expErr: errForbidden,
		},
		"error - missing required scopes": {
			in:         []string{"armory"},
			m2mProfile: iam.NewM2MProfile("imperium|dark_angel", []string{"squad", "relics", "reinforcements"}),
			expErr:     errForbidden,
		},
	}

	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			// Given
			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req = req.WithContext(iam.SetM2MProfileInContext(req.Context(), tc.m2mProfile))

			_, ctx, _ := lit.NewRouterForTest(rr)
			ctx.SetRequest(req)

			// When
			guard := New(nil, nil)
			hdl := guard.RequiredM2MScopeMiddleware(tc.in...)
			hdl(ctx)

			// Then
			if tc.expErr != nil {
				var iamErr lit.HttpError
				if errors.As(tc.expErr, &iamErr) {
					require.Equal(t, rr.Code, iamErr.Status)
				} else {
					require.Equal(t, rr.Code, http.StatusInternalServerError)
				}

				expResult, err := json.Marshal(tc.expErr)
				require.NoError(t, err)
				require.Equal(t, expResult, rr.Body.Bytes())
			}
		})
	}
}

package lit

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/viebiz/lit/i18n"
	"github.com/viebiz/lit/testutil"
)

func TestLitContext_Bind(t *testing.T) {
	type complexStruct struct {
		ID               int `form:"id" json:"id" binding:"required"`
		Equal            int `form:"equal" json:"equal" binding:"eq=100"`
		NotEqual         int `form:"notequal" json:"notequal" binding:"ne=0"`
		LessThan         int `form:"lessthan" json:"lessthan" binding:"lt=50"`
		LessThanEqual    int `form:"lessthanequal" json:"lessthanequal" binding:"lte=50"`
		GreaterThan      int `form:"greaterthan" json:"greaterthan" binding:"gt=10"`
		GreaterThanEqual int `form:"greaterthanequal" json:"greaterthanequal" binding:"gte=10"`
		Multi            int `form:"multi" json:"multi" binding:"min=10,max=20,required"`
	}

	tcs := map[string]struct {
		givenContentType string
		givenRequestBody string
		expectedErr      error
	}{
		"success json": {
			givenContentType: "application/json",
			givenRequestBody: `{
				"id":1,
				"equal":100,
				"notequal":5,
				"lessthan":30,
				"lessthanequal":40,
				"greaterthan":20,
				"greaterthanequal":10,
				"multi":15
			}`,
			expectedErr: nil,
		},
		"got error json": {
			givenContentType: "application/json",
			givenRequestBody: `{
				"equal":99,
				"notequal":0,
				"lessthan":60,
				"lessthanequal":60,
				"greaterthan":10,
				"greaterthanequal":5,
				"multi":30
			}`,
			expectedErr: ValidationError{
				"ID":               "The ID field is required",
				"Equal":            "The Equal field must be 100",
				"NotEqual":         "ne",
				"LessThan":         "lt",
				"LessThanEqual":    "lte",
				"GreaterThan":      "gt",
				"GreaterThanEqual": "gte",
				"Multi":            "The Multi field must be at most 20 but got 30",
			},
		},
		"success form": {
			givenContentType: "application/x-www-form-urlencoded",
			givenRequestBody: func() string {
				f := url.Values{}
				f.Add("id", "1")
				f.Add("equal", "100")
				f.Add("notequal", "5")
				f.Add("lessthan", "30")
				f.Add("lessthanequal", "40")
				f.Add("greaterthan", "20")
				f.Add("greaterthanequal", "10")
				f.Add("multi", "20")
				return f.Encode()
			}(),
			expectedErr: nil,
		},
		"got error form": {
			givenContentType: "application/x-www-form-urlencoded",
			givenRequestBody: func() string {
				f := url.Values{}
				f.Add("equal", "9")
				f.Add("notequal", "0")
				f.Add("lessthan", "60")
				f.Add("lessthanequal", "60")
				f.Add("greaterthan", "10")
				f.Add("greaterthanequal", "5")
				f.Add("multi", "3")
				return f.Encode()
			}(),
			expectedErr: ValidationError{
				"ID":               "The ID field is required",
				"Equal":            "The Equal field must be 100",
				"NotEqual":         "ne",
				"LessThan":         "lt",
				"LessThanEqual":    "lte",
				"GreaterThan":      "gt",
				"GreaterThanEqual": "gte",
				"Multi":            "The Multi field must be at least 10 but got 3",
			},
		},
	}

	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			// GIVEN
			w := httptest.NewRecorder()
			ctx := CreateTestContext(w)
			req := httptest.NewRequest(http.MethodPost, "/dummy", bytes.NewBufferString(tc.givenRequestBody))
			req.Header.Set("Content-Type", tc.givenContentType)
			ctx.SetRequest(req)

			var compObj complexStruct

			langBundle := i18n.Init(context.Background(), i18n.BundleConfig{
				SourcePath: "i18n/testdata",
			})
			lc := langBundle.GetLocalize("en")
			ctx.SetRequestContext(i18n.SetInContext(ctx, lc))

			// WHEN
			err := ctx.Bind(&compObj)

			// THEN
			if tc.expectedErr != nil {
				testutil.Equal(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

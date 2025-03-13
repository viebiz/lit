package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/viebiz/lit"
	"github.com/viebiz/lit/i18n"
)

func TestLocalizationMiddleware(t *testing.T) {
	tcs := map[string]struct {
		srcPath        string
		headerValue    string
		acceptLang     []string
		givenMessageID string
		expectedLang   string
		expectedStatus int
		expectedBody   string
	}{
		"Accepted language": {
			srcPath:        "testdata",
			headerValue:    "en",
			acceptLang:     []string{"en"},
			givenMessageID: "helloPerson",
			expectedLang:   "en",
			expectedStatus: http.StatusOK,
			expectedBody:   "\"Hello Loc Dang\"",
		},
		"Non accepted language": {
			srcPath:        "testdata",
			headerValue:    "vi",
			acceptLang:     []string{"en"},
			givenMessageID: "helloPerson",
			expectedLang:   "vi",
			expectedStatus: http.StatusOK,
			expectedBody:   "\"helloPerson\"",
		},
		"Message ID not exist": {
			srcPath:        "testdata",
			headerValue:    "en",
			acceptLang:     []string{"en"},
			givenMessageID: "ThisMessageIDDoesNotExist",
			expectedLang:   "en",
			expectedStatus: http.StatusOK,
			expectedBody:   "\"ThisMessageIDDoesNotExist\"",
		},
		"Missing source path": {
			srcPath:        "", // This will use default source path 'resources/i18n' and it should fail to load the message file.
			headerValue:    "en",
			acceptLang:     []string{"en"},
			givenMessageID: "helloPerson",
			expectedLang:   "en",
			expectedStatus: http.StatusOK,
			expectedBody:   "\"helloPerson\"",
		},
		"Missing header value": {
			srcPath:        "testdata",
			acceptLang:     []string{"en"},
			givenMessageID: "helloPerson",
			expectedLang:   "en",
			expectedStatus: http.StatusOK,
			expectedBody:   "\"helloPerson\"",
		},
	}

	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			// Given
			w := httptest.NewRecorder()
			route, c, hdlRequest := lit.NewRouterForTest(w)

			// Create a request and set the language header.
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			req.Header.Set(defaultLangHeader, tc.headerValue)
			c.SetRequest(req)

			// Use LocalizationMiddleware with the test's accepted languages.
			route.Use(LocalizationMiddleware(context.Background(), Config{
				SourcePath: tc.srcPath,
			}))

			route.Get("/test", func(c lit.Context) error {
				ctx := c.Request().Context()

				// The message "helloPerson" should be localized if found.
				rs := i18n.FromContext(ctx).Localize(tc.givenMessageID, map[string]interface{}{
					"Name": "Loc Dang",
				})
				c.JSON(http.StatusOK, rs)
				return nil
			})

			// When
			hdlRequest()

			// Then
			require.Equal(t, tc.expectedLang, w.Header().Get(defaultLangResponseHeader))
			require.Equal(t, tc.expectedStatus, w.Code)
			require.Equal(t, tc.expectedBody, w.Body.String())
		})
	}
}

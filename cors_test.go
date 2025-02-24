package lit

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewCORSConfig(t *testing.T) {
	type arg struct {
		givenOrigins        []string
		expAllowedOrigins   []string
		expAllowedMethods   []string
		expAllowedHeaders   []string
		expExposedHeaders   []string
		expAllowCredentials bool
		expMaxAge           time.Duration
	}
	tcs := map[string]arg{
		"no origins": {
			givenOrigins:      []string{},
			expAllowedOrigins: []string{},
			expAllowedMethods: []string{
				http.MethodGet,
				http.MethodPost,
				http.MethodPut,
				http.MethodPatch,
				http.MethodDelete,
				http.MethodHead,
				http.MethodOptions},
			expAllowedHeaders: []string{
				"Accept", "Origin", "Content-Type", "Content-Length", "Authorization",
				"traceparent", "tracestate", "baggage",
			},
			expExposedHeaders:   []string{"Link"},
			expAllowCredentials: true,
			expMaxAge:           300,
		},
		"full settings": {
			givenOrigins:      []string{"https://localhost:3000"},
			expAllowedOrigins: []string{"https://localhost:3000"},
			expAllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodOptions},
			expAllowedHeaders: []string{
				"Accept",
				"Authorization",
				"Content-Type",
				"traceparent",
				"tracestate",
				"baggage",
			},
			expExposedHeaders:   []string{"Link"},
			expAllowCredentials: false,
			expMaxAge:           300,
		},
	}
	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given && When:
			corsConfig := NewCORSConfig(tc.givenOrigins)
			corsConfig.SetAllowHeaders(tc.expAllowedHeaders...)
			corsConfig.SetAllowMethods(tc.expAllowedMethods...)
			corsConfig.SetExposeHeaders(tc.expExposedHeaders...)
			corsConfig.SetMaxAge(tc.expMaxAge)
			if !tc.expAllowCredentials {
				corsConfig.DisableCredentials()
			}

			// Then:
			require.Equal(t, tc.expAllowedOrigins, corsConfig.cfg.AllowOrigins)
			require.Equal(t, tc.expAllowedMethods, corsConfig.cfg.AllowMethods)
			require.Equal(t, tc.expAllowedHeaders, corsConfig.cfg.AllowHeaders)
			require.Equal(t, tc.expExposedHeaders, corsConfig.cfg.ExposeHeaders)
			require.Equal(t, tc.expAllowCredentials, corsConfig.cfg.AllowCredentials)
			require.Equal(t, tc.expMaxAge, corsConfig.cfg.MaxAge)
		})
	}
}

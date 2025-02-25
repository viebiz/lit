package lit

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
)

// CORSConfig holds the CORS configuration
type CORSConfig struct {
	cfg cors.Config
}

// NewCORSConfig initializes and returns a CORSConfig with predefined defaults.
//
// Default configurations:
//   - Allowed methods: GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS
//   - Allowed headers:
//   - Standard: "Accept", "Origin", "Content-Type", "Content-Length", "Authorization"
//   - OpenTelemetry: "traceparent", "tracestate", "baggage"
//   - Allow credentials: true (supports cookies and authorization headers)
//   - Max age: 300 seconds (caches preflight response for 5 minutes)
func NewCORSConfig(origins []string) CORSConfig {
	return CORSConfig{
		cfg: cors.Config{
			AllowOrigins: origins,
			AllowMethods: []string{
				http.MethodGet,
				http.MethodPost,
				http.MethodPut,
				http.MethodPatch,
				http.MethodDelete,
				http.MethodHead,
				http.MethodOptions,
			},
			AllowHeaders: []string{
				"Accept", "Origin", "Content-Type", "Content-Length", "Authorization", // Basic headers
				"traceparent", "tracestate", "baggage", // OpenTelemetry headers
			},
			ExposeHeaders:    []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300,
		},
	}
}

func (corsCfg *CORSConfig) SetAllowMethods(methods ...string) {
	corsCfg.cfg.AllowMethods = methods
}

func (corsCfg *CORSConfig) SetAllowHeaders(headers ...string) {
	corsCfg.cfg.AllowHeaders = headers
}

func (corsCfg *CORSConfig) SetExposeHeaders(headers ...string) {
	corsCfg.cfg.ExposeHeaders = headers
}

func (corsCfg *CORSConfig) DisableCredentials() {
	corsCfg.cfg.AllowCredentials = false
}

func (corsCfg *CORSConfig) SetMaxAge(maxAge time.Duration) {
	corsCfg.cfg.MaxAge = maxAge
}

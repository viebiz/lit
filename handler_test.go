package lit

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a test HTTP handler
	handler := Handler(
		context.Background(),
		NewCORSConfig([]string{"*"}), // Allow all origins for testing
		mockRouter,
	)

	// Create a test HTTP server
	server := httptest.NewServer(handler)
	defer server.Close()

	// Test cases
	tests := map[string]struct {
		method         string
		url            string
		expectedStatus int
		expectedBody   string
	}{

		"Health Check":           {method: "GET", url: "/_/healthz", expectedStatus: http.StatusOK, expectedBody: "ok\n"},
		"Custom Route":           {method: "GET", url: "/test-route", expectedStatus: http.StatusOK, expectedBody: `{"message":"Hello, World!"}`},
		"Test Monitor":           {method: "GET", url: "/_/test-monitor", expectedStatus: http.StatusOK, expectedBody: `{"message":"Test Invoked"}`},
		"Profiling Route":        {method: "GET", url: "/_/profile/", expectedStatus: http.StatusOK, expectedBody: ""},
		"Allocs Profiling Route": {method: "GET", url: "/_/profile/allocs", expectedStatus: http.StatusOK, expectedBody: ""},
		"Non-existent Route":     {method: "GET", url: "/not-found", expectedStatus: http.StatusNotFound, expectedBody: ""},
	}

	// Run test cases
	for scenario, tc := range tests {
		t.Run(scenario, func(t *testing.T) {
			//t.Parallel() // httptest.NewServer does not support parallel
			// When
			req, _ := http.NewRequest(tc.method, server.URL+tc.url, nil)
			resp, err := http.DefaultClient.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			// Then
			assert.Equal(t, tc.expectedStatus, resp.StatusCode)
		})
	}
}

func TestHandlerWithProfilingDisabled(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a test HTTP handler with profiling disabled
	handler := Handler(
		context.Background(),
		NewCORSConfig([]string{"*"}),
		mockRouter,
		HandlerWithProfilingDisabled(), // Disable profiling
	)

	// Create a test HTTP server
	server := httptest.NewServer(handler)
	defer server.Close()

	// Check that profiling routes are not accessible
	profilingRoutes := []string{
		"/_/profile/", "/_/profile/cmdline", "/_/profile/profile", "/_/profile/trace",
	}

	for _, route := range profilingRoutes {
		t.Run("Profiling Disabled: "+route, func(t *testing.T) {
			req, err := http.NewRequest("GET", server.URL+route, nil)
			require.NoError(t, err)
			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()
			assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		})
	}
}

func mockRouter(r Router) {
	r.Get("/test-route", func(c Context) error {
		c.JSON(http.StatusOK, map[string]string{"message": "Hello, World!"})
		return nil
	})
}

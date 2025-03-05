package httpclient

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestClient_AllDefaults(t *testing.T) {
	pool := NewSharedCustomPool()
	result, err := NewUnauthenticated(
		Config{URL: "https://localhost:3000/v1/something", Method: http.MethodGet, ServiceName: "svc"}, pool,
	)

	require.NoError(t, err)
	require.Equal(t, pool.Client, result.underlyingClient)
	require.Equal(t, http.MethodGet, result.method)
	require.Equal(t, "https://localhost:3000/v1/something", result.url)
	//require.Equal(t, "", result.userAgent)
	require.Equal(t, "svc", result.serviceName)
	require.Equal(t, defaultMaxRetriesOnErrOrTimeout, int(result.timeoutAndRetryOption.maxRetries))
	require.Equal(t, defaultMaxWaitInclRetries, result.timeoutAndRetryOption.maxWaitInclRetries)
	require.Equal(t, defaultRetryOnTimeout, result.timeoutAndRetryOption.onTimeout)
	require.Equal(t, defaultContentType, result.contentType)
	require.Equal(t, header{}, result.header)
}

func TestNewClient_AllOverrides(t *testing.T) {
	pool := NewSharedCustomPool(
		OverridePoolTimeoutDuration(time.Minute),
		OverridePoolMaxIdleConns(1),
		OverridePoolMaxIdleConnsPerHost(1),
		OverridePoolMaxConnsPerHost(1),
	)
	result, err := NewUnauthenticated(
		Config{URL: "https://localhost:1604/api/v1/users", Method: http.MethodGet, ServiceName: "svc"}, pool,
		OverrideTimeoutAndRetryOption(
			10, 0, 45*time.Second, true, nil),
		OverrideContentType("text/plain"),
		OverrideBaseRequestHeaders(map[string]string{"k": "v"}),
	)

	require.NoError(t, err)
	require.Equal(t, pool.Client, result.underlyingClient)
	require.Equal(t, http.MethodGet, result.method)
	require.Equal(t, "https://localhost:1604/api/v1/users", result.url)
	//require.Equal(t, "", result.userAgent)
	require.Equal(t, "svc", result.serviceName)
	require.Equal(t, 10, int(result.timeoutAndRetryOption.maxRetries))
	require.Equal(t, 45*time.Second, result.timeoutAndRetryOption.maxWaitInclRetries)
	require.Equal(t, true, result.timeoutAndRetryOption.onTimeout)
	require.Equal(t, "text/plain", result.contentType)
	require.Equal(t, header{values: map[string]string{"k": "v"}}, result.header)
}

func TestWithBasicAuth(t *testing.T) {
	// Given:
	ctx := context.Background()
	payload := []byte(`{"id": 1, "name": "titus", "chapter": "Ultramarine"}`)
	query := url.Values{}
	query.Add("key1", "value1")
	query.Add("key2", "value2")
	vars := map[string]string{
		"var1": "one",
	}
	contentType := "application/json"
	apiKeyName := "X-API-KEY"
	apiKeyValue := "this-is-a-secret-key"
	headers := map[string]string{
		"X-Default":   "123456789",
		"X-Overriden": "123456789",
		"X-API-KEY":   apiKeyValue,
	}
	rheaders := map[string]string{
		"X-Overriden": "OVERRIDEN",
	}

	// Mock:
	var called bool
	mockSvr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Assert request payload
		rbody, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		require.Equal(t, string(payload), string(rbody))
		// Assert path variable substitution
		require.False(t, strings.Contains(r.URL.Path, ":var1"))
		require.True(t, strings.HasSuffix(r.URL.Path, "/one"))
		// Assert query params
		require.Equal(t, query.Encode(), r.URL.Query().Encode())
		// Assert request headers
		require.Equal(t, contentType, r.Header.Get("Content-Type"))
		require.Equal(t, rheaders["X-Overriden"], r.Header.Get("X-Overriden"))
		require.Equal(t, headers["X-Default"], r.Header.Get("X-Default"))
		require.Equal(t, headers["X-API-KEY"], r.Header.Get("X-API-KEY"))

		called = true

		w.Header().Set("key", "value")
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("response body"))
	}))
	sURL := mockSvr.URL + "/:var1"

	c, err := NewWithAPIKey(
		Config{URL: sURL, Method: http.MethodPost, ServiceName: "svc"},
		NewSharedCustomPool(),
		APIKeyConfig{Key: apiKeyName, Value: apiKeyValue},
		OverrideBaseRequestHeaders(headers),
		OverrideContentType(contentType),
	)
	require.NoError(t, err)

	// When:
	resp, err := c.Send(ctx, Payload{
		Body:        payload,
		QueryParams: query,
		PathVars:    vars,
		Header:      rheaders,
	})

	// Then:
	require.NoError(t, err)
	require.True(t, called)
	require.Equal(t, http.StatusAccepted, resp.Status)
	require.Equal(t, "response body", string(resp.Body))
}

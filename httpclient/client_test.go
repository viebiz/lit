package httpclient

import (
	"net/http"
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

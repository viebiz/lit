package httpclient

import (
	"net/http"
	"time"
)

// PoolOption alters behaviour of the http.Client
type PoolOption func(c *http.Client, t *http.Transport)

// OverridePoolTimeoutDuration overrides the timeout for each try
func OverridePoolTimeoutDuration(timeout time.Duration) PoolOption {
	return func(c *http.Client, _ *http.Transport) {
		c.Timeout = timeout
	}
}

// OverridePoolMaxIdleConns overrides the max idle conns
func OverridePoolMaxIdleConns(n int) PoolOption {
	return func(_ *http.Client, t *http.Transport) {
		t.MaxIdleConns = n
	}
}

// OverridePoolMaxConnsPerHost overrides the max conns per host
func OverridePoolMaxConnsPerHost(n int) PoolOption {
	return func(_ *http.Client, t *http.Transport) {
		t.MaxConnsPerHost = n
	}
}

// OverridePoolMaxIdleConnsPerHost overrides the max idle conns per host
func OverridePoolMaxIdleConnsPerHost(n int) PoolOption {
	return func(_ *http.Client, t *http.Transport) {
		t.MaxIdleConnsPerHost = n
	}
}

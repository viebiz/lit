package httpclient

import (
	"net/http"
	"time"
)

type PoolOption func(*http.Client)

func WithTimeout(timeout time.Duration) PoolOption {
	return func(c *http.Client) {
		c.Timeout = timeout
	}
}

func WithTransport(transport http.RoundTripper) PoolOption {
	return func(c *http.Client) {
		c.Transport = transport
	}
}

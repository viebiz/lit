package httpclient

import (
	"net/http"
)

// SharedCustomPool is a custom wrapper around http.Client that sets the client timeout to zero so that it can be
// controlled by Client instead.
type SharedCustomPool struct {
	*http.Client
}

func createDefaultTransport() *http.Transport {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConnsPerHost = defaultMaxIdleConnsPerHost
	return t
}

func newUnderlyingClient(t *http.Transport, opts ...PoolOption) *http.Client {
	c := &http.Client{
		Timeout: defaultTimeoutPerTry,
	}
	for _, opt := range opts {
		opt(c, t)
	}
	c.Transport = t
	return c
}

func newUnderlyingCustomClient(t *http.Transport, opts ...PoolOption) *http.Client {
	c := newUnderlyingClient(t, opts...)
	c.Timeout = 0
	return c
}

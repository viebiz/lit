package httpclient

import (
	"net/http"
)

const (
	defaultMaxIdleConns        = 100
	defaultMaxConnsPerHost     = 100
	defaultMaxIdleConnsPerHost = 100
)

func createDefaultTransport() *http.Transport {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConns = defaultMaxIdleConns
	transport.MaxConnsPerHost = defaultMaxConnsPerHost
	transport.MaxIdleConnsPerHost = defaultMaxIdleConnsPerHost

	return transport
}

package httpclient

import (
	"net/http"
)

// NewSharedPool returns a new http.Client instance with customizable options, ensures
// the efficient reuse of resources by pooling HTTP connections, thereby reducing
// overhead and improving performance across multiple clients
//
// Example:
//
//	func main() {
//		sharedPool := NewSharedPool()
//
//		thirdPartyServiceClient1 := thirdPartyService.NewClient(sharedPool)
//		thirdPartyServiceClient2 := thirdPartyService.NewClient(sharedPool)
//	}
//
//	// Inside third party service
//	func (srv thirdPartyService) Send(ctx context.Context) error {
//		srv.sharedPool.Do(ctx,... ) // Implement your logic
//	}
//
// Refer https://www.loginradius.com/blog/engineering/tune-the-go-http-client-for-high-performance/
// NewSharedPool returns a new http.Client instance based on the arguments
func NewSharedPool(opts ...PoolOption) *http.Client {
	return newUnderlyingClient(createDefaultTransport(), opts...)
}

// NewSharedCustomPool returns a new custom http.Client instance with custom retry and timeout options based on the arguments
func NewSharedCustomPool(opts ...PoolOption) *SharedCustomPool {
	return &SharedCustomPool{newUnderlyingCustomClient(createDefaultTransport(), opts...)}
}

// Config holds the base config for Client
type Config struct {
	ServiceName, // The name of the service we call. Will be used with CallName to form the label for logging.
	URL, // The URL we need to call
	Method string // The HTTP Method to be used
}

// NewUnauthenticated returns a new Client instance based on the arguments without any authentication
func NewUnauthenticated(cfg Config, pool *SharedCustomPool, opts ...ClientOption) (*Client, error) {
	return newClient(pool.Client, cfg.URL, cfg.Method, cfg.ServiceName, opts...)
}

// APIKeyConfig holds the config for APIKey auth
type APIKeyConfig struct {
	Key, Value string
}

// NewWithAPIKey creates and returns a new Client instance with API key auth config
func NewWithAPIKey(cfg Config, pool *SharedCustomPool, apiKeyCfg APIKeyConfig, opts ...ClientOption) (*Client, error) {
	c, err := newClient(pool.Client, cfg.URL, cfg.Method, cfg.ServiceName, opts...)
	if err != nil {
		return nil, err
	}

	if c.header.values == nil {
		c.header.values = map[string]string{}
	}
	c.header.values[apiKeyCfg.Key] = apiKeyCfg.Value
	return c, nil
}

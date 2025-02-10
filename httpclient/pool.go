package httpclient

import (
	"net/http"
	"time"
)

const (
	defaultTimeoutPerTry = 10 * time.Second
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
func NewSharedPool(opts ...PoolOption) *http.Client {
	cl := &http.Client{
		Transport: createDefaultTransport(),
		Timeout:   defaultTimeoutPerTry,
	}

	for _, opt := range opts {
		opt(cl)
	}

	return cl
}

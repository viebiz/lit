package httpclient

import (
	"time"
)

const (
	defaultTimeoutPerTry = 10 * time.Second

	// We do not set MaxConnsPerHost by default, as monitoring and evaluating connection saturation can be challenging.
	// However, consumers can override this setting if needed.
	// Reference: https://www.loginradius.com/blog/async/tune-the-go-http-client-for-high-performance/
	defaultMaxIdleConnsPerHost = 100

	defaultRetryOnTimeout           = false
	defaultMaxRetriesOnErrOrTimeout = 0
	defaultMaxWaitInclRetries       = 15 * time.Second
	defaultContentType              = "application/json"
)

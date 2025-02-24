package httpclient

import (
	"context"
	"time"

	"github.com/cenkalti/backoff/v4"
	pkgerrors "github.com/pkg/errors"
)

// timeoutAndRetryOption holds timeout & retry info for client. This is optional. If not provided, it will
// pick up the default config.
type timeoutAndRetryOption struct {
	// Max num of retries. Setting to <= 0 means no retry
	// Default: 0
	maxRetries uint64
	// Max execution wait time per try.
	// Default: 15s
	maxWaitPerTry time.Duration
	// Max execution wait time, regardless of retries.
	// Default: 15s
	maxWaitInclRetries time.Duration
	// Set false to exclude retry on timeout errors.
	// Good for non-idempotent resources (i.e. push notifications)
	// Default: false
	onTimeout bool
	// Retry on certain http status code
	// Default: empty
	onStatusCodes map[int]bool
}

// IsValid checks if the config is valid or not
func (to timeoutAndRetryOption) IsValid() error {
	if to.maxWaitPerTry > to.maxWaitInclRetries {
		return pkgerrors.Wrap(ErrTimeoutAndRetryOptionInvalid, "maxWaitPerTry > maxWaitInclRetries")
	}
	if to.maxWaitPerTry < 0 {
		return pkgerrors.Wrap(ErrTimeoutAndRetryOptionInvalid, "maxWaitPerTry should not be less than zero")
	}
	if to.maxWaitInclRetries < 0 {
		return pkgerrors.Wrap(ErrTimeoutAndRetryOptionInvalid, "maxWaitInclRetries should not be less than zero")
	}
	if to.onTimeout && to.maxRetries == 0 {
		return pkgerrors.Wrap(ErrTimeoutAndRetryOptionInvalid, "maxRetries should not be zero when retry onTimeout is true")
	}
	if len(to.onStatusCodes) > 0 && to.maxRetries == 0 {
		return pkgerrors.Wrap(ErrTimeoutAndRetryOptionInvalid, "maxRetries should not be zero when retry onStatusCode not empty")
	}
	return nil
}

/*
OverrideTimeoutAndRetryOption method overrides default retry configs if those configs are not zero values

		e.g.
			maxRetries: 3, retry max 3 times after request failed.
			maxWaitPerTry: 10 seconds,
						   request failed
								|
	 							2 seconds backoff retry internal wait time
								|
								retry with max 10 seconds wait time
								|
								2 seconds backoff retry internal wait time
								|
								retry with max 10 seconds wait time
								|
								... loop until get request response || reach maxRetries || exceed maxWaitInclRetries.
			maxWaitInclRetries: expected >= 2 seconds + maxRetries * (maxWaitPerTry + 2 seconds), 38 seconds.
								if set to like 30 seconds, will return overflow err at that point before reach maxRetries.
			onTimeout: set to false if no need retry on timeout
			onStatusCodes: [404], retry when get 404 http status code from response
*/
func OverrideTimeoutAndRetryOption(maxRetries uint64, maxWaitPerTry, maxWaitInclRetries time.Duration, onTimeout bool, onStatusCodes []int) ClientOption {
	return func(c *Client) {
		c.timeoutAndRetryOption.maxRetries = maxRetries
		c.timeoutAndRetryOption.maxWaitPerTry = maxWaitPerTry
		c.timeoutAndRetryOption.maxWaitInclRetries = maxWaitInclRetries
		c.timeoutAndRetryOption.onTimeout = onTimeout

		for _, sc := range onStatusCodes {
			c.timeoutAndRetryOption.onStatusCodes[sc] = true
		}
	}
}

// returns the exponential backoff configuration based on Azure best practices
// Reference: https://docs.microsoft.com/en-us/azure/postgresql/concepts-connectivity#handling-transient-errors
func execWithRetry(
	ctx context.Context,
	maxRetries uint64,
	maxWaitInclRetries time.Duration,
	f func() error,
) error {
	b := backoff.NewExponentialBackOff()
	// 1. Wait for 2 seconds before your first retry. (for simplicity, we're just using backoff.InitialInterval to simulate)
	b.InitialInterval = 2 * time.Second
	b.RandomizationFactor = 0
	// 2. For each following retry, the increase the wait exponentially, up to 60 seconds.
	b.MaxElapsedTime = maxWaitInclRetries

	return backoff.Retry(
		f,
		backoff.WithContext(
			backoff.WithMaxRetries(b, maxRetries), // 3. Set a max number of retries at which point your application considers the operation failed.
			ctx,
		),
	)
}

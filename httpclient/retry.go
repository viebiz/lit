package httpclient

import (
	"time"
)

type timeoutAndRetryOption struct {
	maxRetries int

	maxWaitPerTry time.Duration

	maxWaitInclRetries time.Duration

	//onTimeout bool
	//
	//onStatusCodes map[int]bool
}

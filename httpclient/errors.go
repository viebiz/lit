package httpclient

import (
	"errors"
)

var (
	ErrTimeout                      = errors.New("timeout")
	ErrOverflowMaxWait              = errors.New("overflow max wait")
	ErrOperationContextCanceled     = errors.New("operation context canceled")
	ErrTimeoutAndRetryOptionInvalid = errors.New("retry config invalid")
	ErrMissingURL                   = errors.New("url is missing")
	ErrMissingMethod                = errors.New("method is missing")
	ErrMissingServiceName           = errors.New("missing service name")
)

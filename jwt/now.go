package jwt

import (
	"time"
)

var (
	// timeNowFunc is a helper function for testing token claim validation by mocking the current time.
	timeNowFunc = func() time.Time {
		return time.Now().UTC()
	}
)

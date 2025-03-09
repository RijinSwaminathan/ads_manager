package platform

import "errors"

var (
	// ErrRateLimitExceeded is returned when the API rate limit is exceeded
	ErrRateLimitExceeded = errors.New("rate limit exceeded")
)

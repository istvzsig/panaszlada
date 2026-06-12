package internal

import (
	"time"

	"github.com/istvzsig/retryx/retry"
)

func DefaultHTTPPolicy() retry.Policy {
	return retry.Policy{
		MaxAttempts:   3,
		BaseDelay:     100 * time.Millisecond,
		MaxDelay:      2 * time.Second,
		Jitter:        0.2,
		TimeoutPerTry: 3 * time.Second,

		RetryOn: func(err error) bool {
			return retry.RetryHTTPError(err)
		},
	}
}

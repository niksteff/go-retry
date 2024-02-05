package backoff

import (
	"math/rand"
	"time"
)

func NewLinearBackoffFunc(pause time.Duration, maxRetries int) BackoffFunc {
	var retried int

	jitter := func() time.Duration {
		return time.Duration(rand.Intn(500)) * time.Millisecond
	}

	return func() (time.Duration, bool) {
		// if we reached max retries, we return false to indicate we are done
		if retried+1 > maxRetries {
			return pause, false
		}

		retried++
		return pause + jitter(), true
	}
}

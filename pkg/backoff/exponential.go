package backoff

import (
	"math/rand"
	"time"
)

func NewExponentialBackoffFunc(startPause time.Duration, maxPause time.Duration) BackoffFunc {
	var retried int

	jitter := func() time.Duration {
		return time.Duration(rand.Intn(500)) * time.Millisecond
	}

	return func() time.Duration {
		var current time.Duration
		if retried == 0 {
			current = startPause + jitter()
		} else {
			// calculate the exponential duration by multiplying the base duration
			// with 2 to the power of the current done count
			current = (startPause * time.Duration(2<<(retried-1))) + jitter()
		}

		if current >= maxPause {
			retried++
			return maxPause + jitter()
		}

		retried++
		return current
	}
}

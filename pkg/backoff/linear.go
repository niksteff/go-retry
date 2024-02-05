package backoff

import (
	"math/rand"
	"time"
)

func NewLinearBackoffFunc(pause time.Duration) BackoffFunc {
	var retried int

	jitter := func() time.Duration {
		return time.Duration(rand.Intn(500)) * time.Millisecond
	}

	return func() time.Duration {
		retried++
		return pause + jitter()
	}
}

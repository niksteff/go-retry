package backoff

import "time"

type Backoff interface {
	Backoff() time.Duration
}

type BackoffFunc func() time.Duration

func (f BackoffFunc) Backoff() time.Duration {
	return f()
}

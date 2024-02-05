package backoff

import "time"

type Backoff interface {
	Backoff() (time.Duration, bool)
}

type BackoffFunc func() (time.Duration, bool)

func (f BackoffFunc) Backoff() (time.Duration, bool) {
	return f()
}

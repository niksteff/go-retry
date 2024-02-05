package retry

import (
	"fmt"
	"time"
)

type Backoff interface {
	Backoff() (time.Duration, bool)
}

type RetryError struct {
	Err error
}

func (e RetryError) Error() string {
	return fmt.Sprintf("failed after max retries: %s", e.Err.Error())
}

type Retry interface {
	Try() error
}

type RetryableFunc func() error

func (f RetryableFunc) Try() error {
	return f()
}

func Do(f Retry, fb Backoff) error {
	for {
		err := f.Try()
		if err != nil {
			d, next := fb.Backoff()
			if !next {
				return RetryError{Err: err}
			}

			time.Sleep(d)
			continue
		}

		return nil
	}
}

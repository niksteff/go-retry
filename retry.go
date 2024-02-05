package retry

import (
	"fmt"
	"time"
)

type Backoff interface {
	Backoff() time.Duration
}

type RetryError struct {
	Err error // TODO: error collection from each try?
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

func Do(f Retry, tries int, fb Backoff) error {
	for {
		err := f.Try()
		if err != nil {
			// we tried, subtract one try
			tries -= 1
			if tries <= 0 {
				// no tries left, return the last error
				return RetryError{Err: err}
			}

			d := fb.Backoff()
			time.Sleep(d)
			continue
		}

		return nil
	}
}

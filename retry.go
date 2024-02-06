package retry

import (
	"context"
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

func Do(ctx context.Context, f Retry, tries int, fb Backoff) error {
	for {
		select {
		case <-ctx.Done():
			return context.Canceled
		default:
			err := f.Try()
			if err != nil {
				// we tried, subtract one try
				tries--
				if tries <= 0 {
					// no tries left, return the last error
					return RetryError{Err: err}
				}

				time.Sleep(fb.Backoff())
				continue
			}
			
			return nil
		}

	}
}

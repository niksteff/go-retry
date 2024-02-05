package retry_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/niksteff/go-retry"
	"github.com/niksteff/go-retry/pkg/backoff"
)

var (
	ErrTest = fmt.Errorf("test error")
)

func TestRetry(t *testing.T) {
	testData := []struct {
		r     retry.Retry
		tries int
		b     backoff.Backoff
		err   error
	}{
		{
			r:     retry.RetryableFunc(func() error { return nil }),
			tries: 1,
			b:     backoff.BackoffFunc(func() time.Duration { return 0 }),
			err:   nil,
		},
		{
			r:     retry.RetryableFunc(func() error { return ErrTest }),
			tries: 1,
			b:     backoff.BackoffFunc(func() time.Duration { return 100 * time.Millisecond }),
			err:   ErrTest,
		},
	}

	for idx, d := range testData {
		err := retry.Do(d.r, d.tries, d.b)
		if err != nil {
			var rErr retry.RetryError
			if !errors.As(err, &rErr) {
				t.Errorf("test %d failed, expected error of type %T but got %T: %v", idx, rErr, err, err)
			}
		}
	}
}

func TestRetryCounter(t *testing.T) {
	expectedRetries := 3

	r := simulatePassingRetry(t, expectedRetries)

	var madeRetries int
	b := func() backoff.BackoffFunc {

		return func() time.Duration {
			madeRetries++
			t.Logf("backoff %d", madeRetries)

			if madeRetries < 3 {
				return 100 * time.Millisecond
			}

			return 0
		}
	}()

	err := retry.Do(r, expectedRetries, b)
	if err != nil {
		var rErr retry.RetryError
		if !errors.As(err, &rErr) {
			t.Errorf("expected error of type %T but got %T: %v", rErr, err, err)
		}

		return
	}

	if madeRetries != expectedRetries-1 {
		t.Errorf("expected %d retries but got %d", expectedRetries-1, madeRetries)
	}
}

func simulatePassingRetry(t *testing.T, maxRetries int) retry.RetryableFunc {
	var retries int

	return func() error {
		retries++
		t.Logf("retry %d", retries)

		if retries < maxRetries {
			return ErrTest
		}

		return nil
	}
}

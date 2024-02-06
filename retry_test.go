package retry_test

import (
	"context"
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
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		err := retry.Do(ctx, d.r, d.tries, d.b)
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

	var madeRetries int
	f := func(t *testing.T, maxRetries int) retry.RetryableFunc {
		return func() error {
			t.Logf("retry %d", madeRetries+1)
			madeRetries++
			if madeRetries < maxRetries {
				return ErrTest
			}

			return nil
		}
	}(t, expectedRetries)

	b := func() backoff.BackoffFunc {
		return func() time.Duration {
			return 100 * time.Millisecond
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := retry.Do(ctx, f, 5, b)
	if err != nil {
		var rErr retry.RetryError
		if !errors.As(err, &rErr) {
			t.Errorf("expected error of type %T but got %T: %v", rErr, err, err)
		}
	}

	if madeRetries != expectedRetries {
		t.Errorf("expected %d retries but got %d", expectedRetries, madeRetries)
	}
}

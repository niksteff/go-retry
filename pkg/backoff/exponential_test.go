package backoff_test

import (
	"testing"
	"time"

	"github.com/niksteff/go-retry/pkg/backoff"
)

func TestNewExponentialBackoffFunc(t *testing.T) {
	testData := []struct {
		base            time.Duration
		max             time.Duration
		maxRetries      int
		expectedRetries int
	}{
		{base: 100 * time.Millisecond, max: 1 * time.Second, maxRetries: 0, expectedRetries: 0},
		{base: 100 * time.Millisecond, max: 1 * time.Second, maxRetries: 1, expectedRetries: 1},
		{base: 100 * time.Millisecond, max: 1 * time.Second, maxRetries: 10, expectedRetries: 10},
	}

	for idx, d := range testData {
		t.Logf("test case %d", idx)
		f := backoff.NewExponentialBackoffFunc(d.base, d.max, d.maxRetries)

		var count int
		for {
			d, next := f()
			if !next {
				break
			}

			t.Logf("retrying after %s", d)
			count++
		}

		if count != d.expectedRetries {
			t.Fatalf("expected %d retries, got %d", d.expectedRetries, count)
		}
	}
}

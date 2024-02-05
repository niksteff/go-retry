package backoff

import (
	"testing"
	"time"
)

func TestLinearBackoff(t *testing.T) {
	testData := []struct {
		base            time.Duration
		maxRetries      int
		expectedRetries int
	}{
		{base: 100 * time.Millisecond, maxRetries: 0, expectedRetries: 0},
		{base: 100 * time.Millisecond, maxRetries: 1, expectedRetries: 1},
		{base: 100 * time.Millisecond, maxRetries: 10, expectedRetries: 10},
	}

	for idx, d := range testData {
		t.Logf("test case %d", idx)
		f := NewLinearBackoffFunc(d.base, d.maxRetries)

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

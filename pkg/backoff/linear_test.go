package backoff

import (
	"testing"
	"time"
)

func TestLinearBackoff(t *testing.T) {
	testData := []struct {
		waitTime time.Duration
		count    int
	}{
		{waitTime: 100 * time.Millisecond, count: 0},
		{waitTime: 100 * time.Millisecond, count: 1},
		{waitTime: 100 * time.Millisecond, count: 10},
	}

	for idx, d := range testData {
		f := NewLinearBackoffFunc(d.waitTime)

		var count int
		for {
			if count >= d.count {
				break
			}

			waitTime := f()
			count += 1
			if waitTime < d.waitTime {
				t.Errorf("test %d: expected wait time %s, got %s", idx, d.waitTime, waitTime)
				break
			}

			t.Logf("retrying after %s", waitTime)
		}
	}
}

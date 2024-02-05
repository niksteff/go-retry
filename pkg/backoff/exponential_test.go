package backoff_test

import (
	"testing"
	"time"

	"github.com/niksteff/go-retry/pkg/backoff"
)

func TestNewExponentialBackoffFunc(t *testing.T) {
	testData := []struct {
		base             time.Duration
		max              time.Duration
		tries            int
		expectedWaitTime time.Duration
	}{
		{base: 100 * time.Millisecond, max: 1 * time.Second, tries: 0, expectedWaitTime: 100 * time.Millisecond},
		{base: 100 * time.Millisecond, max: 1 * time.Second, tries: 1, expectedWaitTime: 100 * time.Millisecond},
		{base: 100 * time.Millisecond, max: 1 * time.Second, tries: 10, expectedWaitTime: 1 * time.Second},
	}

	for idx, d := range testData {
		t.Logf("test case %d", idx)
		f := backoff.NewExponentialBackoffFunc(d.base, d.max)

		var count int
		for {
			waitTime := f()

			t.Logf("retrying after %s", waitTime)

			count++
			if waitTime < (d.expectedWaitTime * time.Duration(count)) {
				break
			}
		}
	}
}

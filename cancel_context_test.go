package concurrent

import (
	"context"
	"testing"
	"time"
)

func TestCancelContext(t *testing.T) {
	s := NewCancelContext(context.Background())
	check(s, func() { t.Fail() }, func() {})
	s.Cancel()
	check(s, func() {}, func() { t.Fail() })
	s.Reset(context.Background())
	check(s, func() { t.Fail() }, func() {})
	s.Cancel()
	check(s, func() {}, func() { t.Fail() })
}

func check(s *CancelContext, done func(), timeout func()) {
	select {
	case <-s.Done():
		done()
	case <-time.After(10 * time.Millisecond):
		timeout()
	}
}

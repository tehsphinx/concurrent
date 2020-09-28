package concurrent

import (
	"context"
	"sync"
	"time"
)

var _ context.Context = &CancelContext{}

// NewCancelContext creates a new cancel context from given context
// and wraps it in the concurrency save struct CancelContext.
func NewCancelContext(ctx context.Context) *CancelContext {
	s := &CancelContext{}
	s.Reset(ctx)
	return s
}

// CancelContext implements a context with cancelation that can be used from multiple goroutines.
type CancelContext struct {
	ctxCancel context.Context
	cancel    context.CancelFunc
	m         sync.RWMutex
}

// Cancel cancels the context.
func (s *CancelContext) Cancel() {
	s.m.RLock()
	s.cancel()
	s.m.RUnlock()
}

// Done returns the Done channel of the cancel context.
func (s *CancelContext) Done() <-chan struct{} {
	s.m.RLock()
	defer s.m.RUnlock()

	return s.ctxCancel.Done()
}

// Reset creates a new cancel context from the given context, thereby resetting a cancelled context.
func (s *CancelContext) Reset(ctx context.Context) {
	s.m.Lock()
	s.ctxCancel, s.cancel = context.WithCancel(ctx)
	s.m.Unlock()
}

// Deadline implements context.Context
func (s *CancelContext) Deadline() (deadline time.Time, ok bool) {
	s.m.RLock()
	defer s.m.RUnlock()

	return s.ctxCancel.Deadline()
}

// Err implements context.Context
func (s *CancelContext) Err() error {
	s.m.RLock()
	defer s.m.RUnlock()

	return s.ctxCancel.Err()
}

// Value implements context.Context
func (s *CancelContext) Value(key interface{}) interface{} {
	s.m.RLock()
	defer s.m.RUnlock()

	return s.ctxCancel.Value(key)
}

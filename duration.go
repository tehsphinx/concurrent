package concurrent

import (
	"sync"
	"time"
)

// NewDuration creates a new concurrent duration
func NewDuration() *Duration {
	return &Duration{}
}

// Duration implements a cuncurrent duration
type Duration struct {
	d      time.Duration
	dMutex sync.RWMutex
}

// Set sets the duration to given value
func (s *Duration) Set(i time.Duration) {
	s.dMutex.Lock()
	s.d = i
	s.dMutex.Unlock()
}

// Get gets the duration value
func (s *Duration) Get() time.Duration {
	s.dMutex.RLock()
	defer s.dMutex.RUnlock()

	return s.d
}

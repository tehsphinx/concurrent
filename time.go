package concurrent

import (
	"encoding/json"
	"sync"
	"time"
)

// NewTime creates a new concurrent time
func NewTime() *Time {
	return &Time{}
}

// Time is a concurrent time object
type Time struct {
	time  time.Time
	mutex sync.RWMutex
}

// Now sets time to current time
func (s *Time) Now() *Time {
	s.Set(time.Now())
	return s
}

// Set sets the time
func (s *Time) Set(t time.Time) *Time {
	s.mutex.Lock()
	s.time = t
	s.mutex.Unlock()
	return s
}

// Get returns the time
func (s *Time) Get() time.Time {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.time
}

// Since implements time.Since
func (s *Time) Since() time.Duration {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return time.Since(s.time)
}

// String returns the time formatted as string
func (s *Time) String() string {
	return s.Get().String()
}

// MarshalJSON implements save marshalling
func (s *Time) MarshalJSON() ([]byte, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return json.Marshal(s.time)
}

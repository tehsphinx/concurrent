package concurrent

import (
	"sync"
)

// NewSliceString creates a new concurrent slice of strings
func NewSliceString() *SliceString {
	return &SliceString{}
}

// SliceString implements a cuncurrent string
type SliceString struct {
	slice []string
	mutex sync.RWMutex
}

// Set sets the string to given value
func (s *SliceString) Set(slice []string) {
	s.mutex.Lock()
	s.slice = slice
	s.mutex.Unlock()
}

// Get gets the string value
func (s *SliceString) Get() []string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.slice
}

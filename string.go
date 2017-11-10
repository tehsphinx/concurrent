package concurrent

import (
	"sync"
)

// NewString creates a new concurrent string
func NewString() *String {
	return &String{}
}

// String implements a cuncurrent string
type String struct {
	str      string
	strMutex sync.RWMutex
}

// Set sets the string to given value
func (s *String) Set(str string) {
	s.strMutex.Lock()
	s.str = str
	s.strMutex.Unlock()
}

// Get gets the string value
func (s *String) Get() string {
	s.strMutex.RLock()
	defer s.strMutex.RUnlock()

	return s.str
}

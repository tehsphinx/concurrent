package concurrent

import (
	"sync"
)

// NewByte creates a new concurrent byte
func NewByte() *Byte {
	return &Byte{}
}

// Byte implements a cuncurrent byte
type Byte struct {
	byte  byte
	mutex sync.RWMutex
}

// Set sets the byte to given value
func (s *Byte) Set(b byte) {
	s.mutex.Lock()
	s.byte = b
	s.mutex.Unlock()
}

// Get gets the byte value
func (s *Byte) Get() byte {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.byte
}

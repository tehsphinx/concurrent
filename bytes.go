package concurrent

import (
	"sync"
)

// NewBytes creates a new concurrent slice of bytes
func NewBytes() *Bytes {
	return &Bytes{}
}

// Bytes implements a cuncurrent slice of bytes
type Bytes struct {
	byte      []byte
	byteMutex sync.RWMutex
	//chSet    chan []byte
	//chAppend chan []byte
	//chGet    chan chan []byte
}

// Set sets the byte to given value
func (s *Bytes) Set(b []byte) {
	s.byteMutex.Lock()
	s.byte = b
	s.byteMutex.Unlock()
}

// Get gets the byte value
func (s *Bytes) Get() []byte {
	s.byteMutex.RLock()
	defer s.byteMutex.RUnlock()

	return s.byte
}

// Append appends the slice of bytes value to byte
func (s *Bytes) Append(m []byte) {
	s.byteMutex.Lock()
	defer s.byteMutex.Unlock()

	s.byte = append(s.byte, m...)
}

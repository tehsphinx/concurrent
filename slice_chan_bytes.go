package concurrent

import (
	"sync"
)

// NewSliceChanBytes creates a new concurrent slice of bytes
func NewSliceChanBytes() *SliceChanBytes {
	return &SliceChanBytes{
		slice: []chan []byte{},
	}
}

// SliceChanBytes implements a cuncurrent slice of bytes
type SliceChanBytes struct {
	slice []chan []byte
	mutex sync.RWMutex
}

// Add appends a channel to the slice
func (s *SliceChanBytes) Add(ch chan []byte) {
	s.mutex.Lock()
	s.slice = append(s.slice, ch)
	s.mutex.Unlock()
}

// Remove removes a channel from the slice and closes it
func (s *SliceChanBytes) Remove(ch chan []byte) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	found := false
	for i, c := range s.slice {
		if c == ch {
			s.slice = append(s.slice[:i], s.slice[i+1:]...)
			close(ch)
			found = true
		}
	}
	return found
}

// RemoveAll removes alls channels and closes them
func (s *SliceChanBytes) RemoveAll() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, ch := range s.slice {
		close(ch)
	}
	s.slice = []chan []byte{}
}

// Send sends on all channels
func (s *SliceChanBytes) Send(msg []byte) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, ch := range s.slice {
		ch <- msg
	}
}

// Len returns the count of the channels
func (s *SliceChanBytes) Len() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return len(s.slice)
}

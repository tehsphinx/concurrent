package concurrent

import (
	"sync"
)

// NewSliceChanString creates a new concurrent slice of bytes
func NewSliceChanString() *SliceChanString {
	return &SliceChanString{
		slice: []chan string{},
	}
}

// SliceChanString implements a cuncurrent slice of bytes
type SliceChanString struct {
	slice []chan string
	mutex sync.RWMutex
}

// Add appends a channel to the slice
func (s *SliceChanString) Add(ch chan string) {
	s.mutex.Lock()
	s.slice = append(s.slice, ch)
	s.mutex.Unlock()
}

// Remove removes a channel from the slice and closes it
func (s *SliceChanString) Remove(ch chan string) bool {
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
func (s *SliceChanString) RemoveAll() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, ch := range s.slice {
		close(ch)
	}
	s.slice = []chan string{}
}

// Send sends on all channels
func (s *SliceChanString) Send(msg string) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, ch := range s.slice {
		ch <- msg
	}
}

// Len returns the count of the channels
func (s *SliceChanString) Len() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return len(s.slice)
}

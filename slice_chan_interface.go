package concurrent

import (
	"sync"
)

// NewSliceChanInterface creates a new concurrent slice of chan interface{}
func NewSliceChanInterface() *SliceChanInterface {
	return &SliceChanInterface{
		slice: []chan interface{}{},
	}
}

// SliceChanInterface implements a cuncurrent slice of chan interface{}
type SliceChanInterface struct {
	slice []chan interface{}
	mutex sync.RWMutex
}

// Add appends a channel to the slice
func (s *SliceChanInterface) Add(ch chan interface{}) {
	s.mutex.Lock()
	s.slice = append(s.slice, ch)
	s.mutex.Unlock()
}

// Remove removes a channel from the slice and closes it
func (s *SliceChanInterface) Remove(ch chan interface{}) bool {
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
func (s *SliceChanInterface) RemoveAll() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, ch := range s.slice {
		close(ch)
	}
	s.slice = []chan interface{}{}
}

// Send sends on all channels
func (s *SliceChanInterface) Send(msg interface{}) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, ch := range s.slice {
		ch <- msg
	}
}

// Len returns the count of the channels
func (s *SliceChanInterface) Len() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return len(s.slice)
}

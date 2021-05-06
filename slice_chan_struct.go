package concurrent

import (
	"sync"
)

// NewSliceChanStruct creates a new concurrent slice of chan string
func NewSliceChanStruct() *SliceChanStruct {
	return &SliceChanStruct{
		slice: []chan struct{}{},
	}
}

// SliceChanStruct implements a cuncurrent slice of chan string
type SliceChanStruct struct {
	slice []chan struct{}
	mutex sync.RWMutex
}

// Add appends a channel to the slice
func (s *SliceChanStruct) Add(ch chan struct{}) {
	s.mutex.Lock()
	s.slice = append(s.slice, ch)
	s.mutex.Unlock()
}

// Remove removes a channel from the slice and closes it
func (s *SliceChanStruct) Remove(ch chan struct{}) bool {
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
func (s *SliceChanStruct) RemoveAll() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, ch := range s.slice {
		close(ch)
	}
	s.slice = []chan struct{}{}
}

// Send sends on all channels
func (s *SliceChanStruct) Send() {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, ch := range s.slice {
		ch <- struct{}{}
	}
}

// SendNonBlocking sends on all channels. If a channel is blocking it will skip that channel.
func (s *SliceChanStruct) SendNonBlocking() {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, ch := range s.slice {
		select {
		case ch <- struct{}{}:
		default:
		}
	}
}

// Len returns the count of the channels
func (s *SliceChanStruct) Len() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return len(s.slice)
}

package concurrent

import (
	"sync"
)

var eventChannelLen = 5

// NewBool creates a new concurrent bool
func NewBool() *Bool {
	return &Bool{
		events: []chan bool{},
	}
}

// Bool implements a cuncurrent bool
type Bool struct {
	bool      bool
	boolMutex sync.RWMutex
	events    []chan bool
}

// Set sets the bool to given value
func (s *Bool) Set(b bool) {
	s.boolMutex.Lock()
	defer s.boolMutex.Unlock()

	send := s.bool != b
	s.bool = b

	if send {
		for _, e := range s.events {
			if len(e) < eventChannelLen {
				e <- b
			}
		}
	}
}

// Get gets the bool value
func (s *Bool) Get() bool {
	s.boolMutex.RLock()
	defer s.boolMutex.RUnlock()

	return s.bool
}

// GetStatusChannel returns a channel to listen to for status changes
func (s *Bool) GetStatusChannel() <-chan bool {
	ch := make(chan bool, eventChannelLen)
	s.boolMutex.Lock()
	s.events = append(s.events, ch)
	s.boolMutex.Unlock()
	return ch
}

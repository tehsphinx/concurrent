package concurrent

import (
	"sync"
)

var eventChannelLen = 5

// NewBoolWithStatus creates a new concurrent bool with status channel
func NewBoolWithStatus() *BoolWithStatus {
	return &BoolWithStatus{
		events: []chan bool{},
	}
}

// BoolWithStatus implements a cuncurrent bool
type BoolWithStatus struct {
	bool      bool
	boolMutex sync.RWMutex
	events    []chan bool
}

// Set sets the bool to given value and returns if it changed or not. This
// can be used in race cases where a value could change inbetween the 'if'
// and setting the new value, making sure two routines executing simultaneously
// will not both get the same result checking the bool.
func (s *BoolWithStatus) Set(b bool) (ok bool) {
	s.boolMutex.Lock()
	defer s.boolMutex.Unlock()

	ok = s.bool != b
	if ok {
		s.bool = b
		for _, e := range s.events {
			if len(e) < eventChannelLen {
				e <- b
			}
		}
	}
	return ok
}

// Get gets the bool value
func (s *BoolWithStatus) Get() bool {
	s.boolMutex.RLock()
	defer s.boolMutex.RUnlock()

	return s.bool
}

// GetStatusChannel returns a channel to listen to for status changes
func (s *BoolWithStatus) GetStatusChannel() <-chan bool {
	ch := make(chan bool, eventChannelLen)
	s.boolMutex.Lock()
	s.events = append(s.events, ch)
	s.boolMutex.Unlock()
	return ch
}

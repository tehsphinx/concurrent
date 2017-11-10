package concurrent

import (
	"sync"
)

// NewInt creates a new concurrent int
func NewInt() *Int {
	return &Int{}
}

// Int implements a cuncurrent int
type Int struct {
	int      int
	intMutex sync.RWMutex
}

// Set sets the int to given value
func (s *Int) Set(i int) {
	s.intMutex.Lock()
	s.int = i
	s.intMutex.Unlock()
}

// Get gets the int value
func (s *Int) Get() int {
	s.intMutex.RLock()
	defer s.intMutex.RUnlock()

	return s.int
}

// Decrease decreases the integer
func (s *Int) Decrease() int {
	s.intMutex.Lock()
	defer s.intMutex.Unlock()

	s.int--
	return s.int
}

// Increase increases the integer
func (s *Int) Increase() int {
	s.intMutex.Lock()
	defer s.intMutex.Unlock()

	s.int++
	return s.int
}

// NewIntEvent creates a new integer event
func NewIntEvent() *IntEvent {
	return &IntEvent{
		chEvent: make(chan int, 1),
	}
}

// IntEvent is a event-like object that can be listened to
// and fired. The specialty is that it does not block if
// no one is listening. Currently only supports one listener.
type IntEvent struct {
	chEvent chan int
}

// ListenerChannel returns the channel to listen on
func (s *IntEvent) ListenerChannel() <-chan int {
	return s.chEvent
}

// Fire fires the event
func (s *IntEvent) Fire(val int) {
	if len(s.chEvent) != 0 {
		<-s.chEvent
	}
	s.chEvent <- val
}

// Close closes the event by closing the channel
func (s *IntEvent) Close() {
	close(s.chEvent)
}

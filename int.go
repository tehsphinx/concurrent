package concurrent

import (
	"encoding/json"
	"sync/atomic"
)

// NewInt creates a new concurrent int
func NewInt() *Int {
	return new(Int)
}

// Int implements a cuncurrent int
type Int int64

// Set sets the int to given value
func (s *Int) Set(i int) {
	atomic.StoreInt64((*int64)(s), int64(i))
}

// Get gets the int value
func (s *Int) Get() int {
	return int(atomic.LoadInt64((*int64)(s)))
}

// Decrease decreases the integer
func (s *Int) Decrease() int {
	return int(atomic.AddInt64((*int64)(s), -1))
}

// Increase increases the integer
func (s *Int) Increase() int {
	return int(atomic.AddInt64((*int64)(s), 1))
}

// Swap sets the int to given value and returns the old one
func (s *Int) Swap(i int) int {
	return int(atomic.SwapInt64((*int64)(s), int64(i)))
}

// MarshalJSON adds json marshalling to the concurrent bool
func (s *Int) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Get())
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

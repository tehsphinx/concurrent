package concurrent

import "log"

// NewInt creates a new concurrent int
func NewInt() *Int {
	s := &Int{
		chSet:      make(chan int, 5),
		chGet:      make(chan chan int, 5),
		chIncrease: make(chan chan int, 5),
		chDecrease: make(chan chan int, 5),
	}
	return s.run()
}

// Int implements a cuncurrent int
type Int struct {
	int        int
	chSet      chan int
	chGet      chan chan int
	chIncrease chan chan int
	chDecrease chan chan int
}

// Set sets the int to given value
func (s *Int) Set(i int) {
	s.chSet <- i
}

// Get gets the int value
func (s *Int) Get() int {
	ch := make(chan int)
	s.chGet <- ch
	return <-ch
}

// Decrease decreases the integer
func (s *Int) Decrease() int {
	ch := make(chan int)
	s.chDecrease <- ch
	return <-ch
}

// Increase increases the integer
func (s *Int) Increase() int {
	ch := make(chan int)
	s.chIncrease <- ch
	return <-ch
}

func (s *Int) run() *Int {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("concurrent.(Int).run panicked:", r)
			}

			s.run()
		}()

		for {
			select {
			case s.int = <-s.chSet:
			case ch := <-s.chGet:
				ch <- s.int
			case ch := <-s.chIncrease:
				s.int++
				ch <- s.int
			case ch := <-s.chDecrease:
				s.int--
				ch <- s.int
			}
		}
	}()

	return s
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

package concurrent

import "log"

// NewBool creates a new concurrent bool
func NewBool() *Bool {
	s := &Bool{
		chSet:      make(chan bool),
		chGet:      make(chan chan bool, 5),
		events:     []chan bool{},
		chAddEvent: make(chan chan bool),
	}
	return s.run()
}

// Bool implements a cuncurrent bool
type Bool struct {
	bool       bool
	chSet      chan bool
	chGet      chan chan bool
	events     []chan bool
	chAddEvent chan chan bool
}

// Set sets the bool to given value
func (s *Bool) Set(b bool) {
	s.chSet <- b
}

// Get gets the bool value
func (s *Bool) Get() bool {
	ch := make(chan bool)
	s.chGet <- ch
	return <-ch
}

func (s *Bool) GetStatusChannel() <-chan bool {
	ch := make(chan bool, 5)
	s.chAddEvent <- ch
	return ch
}

func (s *Bool) run() *Bool {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("concurrent.(Bool).run panicked:", r)
			}

			s.run()
		}()

		for {
			select {
			case b := <-s.chSet:
				send := s.bool != b
				s.bool = b

				if send {
					for _, e := range s.events {
						e <- b
					}
				}
			case ch := <-s.chGet:
				ch <- s.bool
			case ch := <-s.chAddEvent:
				s.events = append(s.events, ch)
			}
		}
	}()

	return s
}

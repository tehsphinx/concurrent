package concurrent

import "log"

// NewBool creates a new concurrent bool
func NewBool() *Bool {
	s := &Bool{
		chSet: make(chan bool, 5),
		chGet: make(chan chan bool, 5),
	}
	return s.run()
}

// Bool implements a cuncurrent bool
type Bool struct {
	bool  bool
	chSet chan bool
	chGet chan chan bool
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
			case s.bool = <-s.chSet:
			case ch := <-s.chGet:
				ch <- s.bool
			}
		}
	}()

	return s
}

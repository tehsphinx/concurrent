package concurrent

import "log"

// NewString creates a new concurrent string
func NewString() *String {
	s := &String{
		chSet: make(chan string, 5),
		chGet: make(chan chan string, 5),
	}
	return s.run()
}

// String implements a cuncurrent string
type String struct {
	str   string
	chSet chan string
	chGet chan chan string
}

// Set sets the string to given value
func (s *String) Set(str string) {
	s.chSet <- str
}

// Get gets the string value
func (s *String) Get() string {
	ch := make(chan string)
	s.chGet <- ch
	return <-ch
}

func (s *String) run() *String {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("concurrent.(String).run panicked:", r)
			}

			s.run()
		}()

		for {
			select {
			case s.str = <-s.chSet:
			case ch := <-s.chGet:
				ch <- s.str
			}
		}
	}()

	return s
}

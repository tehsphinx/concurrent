package concurrent

import "log"

// NewByte creates a new concurrent byte
func NewByte() *Byte {
	s := &Byte{
		chSet: make(chan byte, 5),
		chGet: make(chan chan byte, 5),
	}
	return s.run()
}

// Byte implements a cuncurrent byte
type Byte struct {
	byte  byte
	chSet chan byte
	chGet chan chan byte
}

// Set sets the byte to given value
func (s *Byte) Set(i byte) {
	s.chSet <- i
}

// Get gets the byte value
func (s *Byte) Get() byte {
	ch := make(chan byte)
	s.chGet <- ch
	return <-ch
}

func (s *Byte) run() *Byte {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("concurrent.(Byte).run panicked:", r)
			}

			s.run()
		}()

		for {
			select {
			case s.byte = <-s.chSet:
			case ch := <-s.chGet:
				ch <- s.byte
			}
		}
	}()

	return s
}

package concurrent

import (
	"log"
)

// NewByte creates a new concurrent slice of bytes
func NewByte() *Byte {
	s := &Byte{
		chSet:    make(chan []byte, 5),
		chAppend: make(chan []byte, 5),
		chGet:    make(chan chan []byte, 5),
	}
	return s.run()
}

// Byte implements a cuncurrent slice of bytes
type Byte struct {
	byte     []byte
	chSet    chan []byte
	chAppend chan []byte
	chGet    chan chan []byte
}

// Set sets the byte to given value
func (s *Byte) Set(alphabet []byte) {
	s.chSet <- alphabet
}

// Get gets the byte value
func (s *Byte) Get() []byte {
	ch := make(chan []byte)
	s.chGet <- ch
	return <-ch
}

// Append appends the slice of bytes value to byte
func (s *Byte) Append(m []byte) {
	s.chAppend <- m
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
			case addbyte := <-s.chAppend:
				s.byte = append(s.byte, addbyte...)
			case s.byte = <-s.chSet:
			case ch := <-s.chGet:
				ch <- s.byte
			}
		}
	}()

	return s
}

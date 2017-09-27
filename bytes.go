package concurrent

import (
	"log"
)

// NewBytes creates a new concurrent slice of bytes
func NewBytes() *Bytes {
	s := &Bytes{
		chSet:    make(chan []byte, 5),
		chAppend: make(chan []byte, 5),
		chGet:    make(chan chan []byte, 5),
	}
	return s.run()
}

// Bytes implements a cuncurrent slice of bytes
type Bytes struct {
	byte     []byte
	chSet    chan []byte
	chAppend chan []byte
	chGet    chan chan []byte
}

// Set sets the byte to given value
func (s *Bytes) Set(alphabet []byte) {
	s.chSet <- alphabet
}

// Get gets the byte value
func (s *Bytes) Get() []byte {
	ch := make(chan []byte)
	s.chGet <- ch
	return <-ch
}

// Append appends the slice of bytes value to byte
func (s *Bytes) Append(m []byte) {
	s.chAppend <- m
}

func (s *Bytes) run() *Bytes {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("concurrent.(Bytes).run panicked:", r)
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

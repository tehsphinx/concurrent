package concurrent

import (
	"log"
)

type actionSliceString int

const (
	actionSetSliceString actionSliceString = iota
	actionGetSliceString
)

type cmdSliceString struct {
	action  actionSliceString
	slice   []string
	chSlice chan []string
}

// NewSliceString creates a new concurrent slice of strings
func NewSliceString() *SliceString {
	s := &SliceString{
		chCmd: make(chan cmdSliceString, 5),
	}
	return s.run()
}

// SliceString implements a cuncurrent string
type SliceString struct {
	slice []string
	chCmd chan cmdSliceString
}

// Set sets the string to given value
func (s *SliceString) Set(slice []string) {
	s.chCmd <- cmdSliceString{
		action: actionSetSliceString,
		slice:  slice,
	}
}

// Get gets the string value
func (s *SliceString) Get() []string {
	ch := make(chan []string)
	s.chCmd <- cmdSliceString{
		action:  actionGetSliceString,
		chSlice: ch,
	}
	return <-ch
}

func (s *SliceString) run() *SliceString {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("concurrent.(SliceString).run panicked:", r)
			}

			s.run()
		}()

		for c := range s.chCmd {
			switch c.action {
			case actionSetSliceString:
				s.slice = c.slice
			case actionGetSliceString:
				c.chSlice <- s.slice
			}
		}
	}()

	return s
}

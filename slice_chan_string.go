package concurrent

import (
	"log"
)

type actionSliceChanString int

const (
	actionAddChanString actionSliceChanString = iota
	actionRemoveChanString
	actionRemoveAllChanString
	actionSendChanString
	actionLenChanString
)

type cmdSliceChanString struct {
	action actionSliceChanString
	ch     chan string
	chOk   chan bool
	chLen  chan int
	msg    string
}

// NewSliceChanString creates a new concurrent slice of bytes
func NewSliceChanString() *SliceChanString {
	s := &SliceChanString{
		slice: []chan string{},
		chCmd: make(chan cmdSliceChanString, 5),
	}
	return s.run()
}

// SliceChanString implements a cuncurrent slice of bytes
type SliceChanString struct {
	slice []chan string
	chCmd chan cmdSliceChanString
}

// Add appends a channel to the slice
func (s *SliceChanString) Add(ch chan string) {
	s.chCmd <- cmdSliceChanString{
		action: actionAddChanString,
		ch:     ch,
	}
}

// Remove removes a channel from the slice and closes it
func (s *SliceChanString) Remove(ch chan string) bool {
	chOk := make(chan bool)
	s.chCmd <- cmdSliceChanString{
		action: actionRemoveChanString,
		ch:     ch,
		chOk:   chOk,
	}
	return <-chOk
}

// RemoveAll removes alls channels and closes them
func (s *SliceChanString) RemoveAll() {
	s.chCmd <- cmdSliceChanString{
		action: actionRemoveAllChanString,
	}
}

// Send sends on all channels
func (s *SliceChanString) Send(msg string) {
	s.chCmd <- cmdSliceChanString{
		action: actionSendChanString,
		msg:    msg,
	}
}

// Len returns the count of the channels
func (s *SliceChanString) Len() int {
	chLen := make(chan int)
	s.chCmd <- cmdSliceChanString{
		action: actionLenChanString,
		chLen:  chLen,
	}
	return <-chLen
}

func (s *SliceChanString) run() *SliceChanString {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("concurrent.(SliceChanString).run panicked:", r)
			}

			s.run()
		}()

		for c := range s.chCmd {
			switch c.action {
			case actionAddChanString:
				s.slice = append(s.slice, c.ch)
			case actionRemoveChanString:
				found := false
				for i, ch := range s.slice {
					if ch == c.ch {
						s.slice = append(s.slice[:i], s.slice[i+1:]...)
						close(ch)
						found = true
					}
				}
				c.chOk <- found
			case actionRemoveAllChanString:
				for _, ch := range s.slice {
					close(ch)
				}
				s.slice = []chan string{}
			case actionSendChanString:
				for _, ch := range s.slice {
					ch <- c.msg
				}
			case actionLenChanString:
				c.chLen <- len(s.slice)
			}
		}
	}()

	return s
}

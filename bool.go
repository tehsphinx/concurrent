package concurrent

import (
	"sync"
	"sync/atomic"
)

// NewBool creates a new concurrent bool
func NewBool() *Bool {
	return &Bool{
		events: []chan bool{},
	}
}

// Bool implements a cuncurrent bool
type Bool struct {
	bool      int32
	boolMutex sync.RWMutex
	events    []chan bool
}

// Set sets the bool to given value and returns if it changed or not. This
// can be used in race cases where a value could change inbetween the 'if'
// and setting the new value, making sure two routines executing simultaneously
// will not both get the same result checking the bool.
func (s *Bool) Set(b bool) (ok bool) {
	var o, n int32
	if !b {
		o = 1
	} else {
		n = 1
	}
	return atomic.CompareAndSwapInt32((*int32)(&s.bool), o, n)
}

// Get gets the bool value
func (s *Bool) Get() bool {
	return atomic.LoadInt32((*int32)(&s.bool)) == 1
}

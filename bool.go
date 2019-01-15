package concurrent

import (
	"encoding/json"
	"sync/atomic"
)

// NewBool creates a new concurrent bool
func NewBool() *Bool {
	return new(Bool)
}

// Bool implements a concurrent bool
type Bool int32

// Set sets the bool to given value and returns if it changed or not. This
// can be used in race cases where a value could change inbetween the 'if'
// and setting the new value, making sure two routines executing simultaneously
// will not both get the same result checking the bool.
func (s *Bool) Set(b bool) (ok bool) {
	var o int32
	if !b {
		o = 1
	}
	return atomic.CompareAndSwapInt32((*int32)(s), o, 1-o)
}

// Get gets the bool value
func (s *Bool) Get() bool {
	return atomic.LoadInt32((*int32)(s)) == 1
}

// MarshalJSON adds json marshalling to the concurrent bool
func (s *Bool) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Get())
}

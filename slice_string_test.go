package concurrent

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSliceString(t *testing.T) {
	s := NewSliceString()

	slice := []string{"a", "b"}
	s.Set(slice)
	assert.Equal(t, slice, s.Get())

	slice = []string{"c"}
	s.Set(slice)
	assert.Equal(t, slice, s.Get())

	for i := 0; i < 100; i++ {
		go func() {
			for {
				s.Get()
			}
		}()
	}

	for i := 0; i < 10; i++ {
		go func(i string) {
			for {
				s.Set([]string{"string" + i})
			}
		}(strconv.Itoa(i))
	}

	time.Sleep(5 * time.Second)
}

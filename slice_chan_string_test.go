package concurrent

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSliceChanString(t *testing.T) {
	s := NewSliceChanString()

	ch := make(chan string, 5)
	s.Add(ch)
	assert.Equal(t, 1, s.Len())

	ch2 := make(chan string, 5)
	s.Add(ch2)
	assert.Equal(t, 2, s.Len())

	ch3 := make(chan string, 5)
	s.Add(ch3)
	assert.Equal(t, 3, s.Len())

	ch4 := make(chan string, 5)
	s.Add(ch4)
	assert.Equal(t, 4, s.Len())

	str := "bla"
	s.Send(str)
	// TODO: remove?
	time.Sleep(200 * time.Millisecond)
	assert.Equal(t, len(ch), 1)
	assert.Equal(t, len(ch2), 1)
	assert.Equal(t, len(ch3), 1)
	assert.Equal(t, len(ch3), 1)
	assert.Equal(t, len(ch4), 1)
	assert.Equal(t, str, <-ch)
	assert.Equal(t, str, <-ch2)
	assert.Equal(t, str, <-ch3)
	assert.Equal(t, str, <-ch4)

	assert.True(t, s.Remove(ch2))
	assert.Equal(t, 3, s.Len())

	assert.True(t, s.Remove(ch4))
	assert.Equal(t, 2, s.Len())

	s.RemoveAll()
	assert.Equal(t, 0, s.Len())
}

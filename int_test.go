package concurrent

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInt(t *testing.T) {
	s := NewInt()

	for i := 0; i < 10; i++ {
		s.Set(i)
		k := s.Get()
		assert.Equal(t, i, k)
	}
}

func TestInt_Increase(t *testing.T) {
	s := NewInt()

	len := 10
	for i := 0; i < len; i++ {
		s.Increase()
	}
	k := s.Get()
	assert.Equal(t, len, k)
}

func TestInt_Decrease(t *testing.T) {
	s := NewInt()

	len := 10
	for i := 0; i < len; i++ {
		s.Decrease()
	}
	k := s.Get()
	assert.Equal(t, -len, k)
}
func TestInt_Concurrent(t *testing.T) {
	s := NewInt()

	for i := 0; i < 100; i++ {
		go func() {
			for {
				s.Get()
			}
		}()
	}

	for i := 0; i < 10; i++ {
		go func(i int) {
			for {
				s.Set(i)
			}
		}(i)
	}

	time.Sleep(5 * time.Second)
}

func BenchmarkInt_Get(b *testing.B) {
	s := NewInt()
	wg := &sync.WaitGroup{}
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func(wg *sync.WaitGroup) {
			for i := 0; i < b.N/10; i++ {
				s.Get()
			}
			wg.Done()
		}(wg)
	}
	wg.Wait()
}

func BenchmarkInt_Set(b *testing.B) {
	s := NewInt()
	wg := &sync.WaitGroup{}
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func(wg *sync.WaitGroup) {
			for i := 0; i < b.N/10; i++ {
				s.Set(i)
			}
			wg.Done()
		}(wg)
	}
	wg.Wait()
}

func TestIntEvent(t *testing.T) {
	event := NewIntEvent()
	ch := event.ListenerChannel()

	go func(ch <-chan int) {
		i := 0
		for val := range ch {
			assert.Equal(t, i, val)
			i++
		}
	}(ch)

	for i := 0; i < 10; i++ {
		event.Fire(i)
		// wait a bit so we do not overwrite values
		time.Sleep(10 * time.Millisecond)
	}
	event.Close()

	// no listener this time. Should not block!
	event = NewIntEvent()
	ch = event.ListenerChannel()
	for i := 0; i < 10; i++ {
		event.Fire(i)
	}

	// should get last last value
	assert.Equal(t, 9, <-ch)
}

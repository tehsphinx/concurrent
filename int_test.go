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

	l := 10
	for i := 0; i < l; i++ {
		s.Increase()
	}
	k := s.Get()
	assert.Equal(t, l, k)
}

func TestInt_Decrease(t *testing.T) {
	s := NewInt()

	l := 10
	for i := 0; i < l; i++ {
		s.Decrease()
	}
	k := s.Get()
	assert.Equal(t, -l, k)
}

func TestInt_Concurrent(t *testing.T) {
	s := NewInt()

	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		for i := 0; i < 1<<30; i++ {
			s.Get()
		}
		wg.Done()
	}()

	for i := 0; i < 2; i++ {
		go func(n int) {
			for i := 0; i < n; i++ {
				s.Set(i)
			}
			wg.Done()
		}(i << 30)
	}
	wg.Wait()
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

package concurrent

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDuration(t *testing.T) {
	s := NewDuration()

	for i := 0; i < 10; i++ {
		s.Set(time.Duration(i))
		k := s.Get()
		assert.Equal(t, time.Duration(i), k)
	}
}

func TestDuration_Concurrent(t *testing.T) {
	s := NewDuration()

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
				s.Set(time.Duration(i))
			}
		}(i)
	}

	time.Sleep(5 * time.Second)
}

func BenchmarkDuration_Get(b *testing.B) {
	s := NewDuration()
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

func BenchmarkDuration_Set(b *testing.B) {
	s := NewDuration()
	wg := &sync.WaitGroup{}
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func(wg *sync.WaitGroup) {
			for i := 0; i < b.N/10; i++ {
				s.Set(time.Duration(i))
			}
			wg.Done()
		}(wg)
	}
	wg.Wait()
}

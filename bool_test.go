package concurrent

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBool(t *testing.T) {
	s := NewBool()

	s.Set(false)
	assert.Equal(t, false, s.Get())
	s.Set(true)
	assert.Equal(t, true, s.Get())
	s.Set(false)
	assert.Equal(t, false, s.Get())
	s.Set(false)
	assert.Equal(t, false, s.Get())
	s.Set(true)
	assert.Equal(t, true, s.Get())
}

func TestBoolConcurrent(t *testing.T) {
	s := NewBool()

	for i := 0; i < 100; i++ {
		go func() {
			for {
				s.Get()
			}
		}()
	}

	for i := 0; i < 10; i++ {
		go func(i bool) {
			for {
				s.Set(i)
			}
		}(i%2 == 0)
	}

	time.Sleep(5 * time.Second)
}

func BenchmarkBool_Get(b *testing.B) {
	s := NewBool()
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

func BenchmarkBool_Set(b *testing.B) {
	s := NewBool()
	wg := &sync.WaitGroup{}
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func(wg *sync.WaitGroup) {
			for i := 0; i < b.N/10; i++ {
				s.Set(i%2 == 0)
			}
			wg.Done()
		}(wg)
	}
	wg.Wait()
}

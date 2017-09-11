package concurrent

import (
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestString(t *testing.T) {
	s := NewString()

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
				s.Set("string" + i)
			}
		}(strconv.Itoa(i))
	}

	time.Sleep(5 * time.Second)
}

func BenchmarkString_Get(b *testing.B) {
	s := NewString()
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

func BenchmarkString_Set(b *testing.B) {
	s := NewString()
	wg := &sync.WaitGroup{}
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func(wg *sync.WaitGroup) {
			for i := 0; i < b.N/10; i++ {
				s.Set("somestring")
			}
			wg.Done()
		}(wg)
	}
	wg.Wait()
}

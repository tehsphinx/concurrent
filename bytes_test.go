package concurrent

import (
	"sync"
	"testing"
	"time"
)

func TestBytes(t *testing.T) {
	s := NewBytes()
	alphabet := []byte("abcdefghijklmnopqrstuvwxyz")
	settovar := []byte("zzzzzzzzzz")

	for i := 0; i < 10; i++ {
		go func() {
			for {
				s.Get()
			}
		}()
	}

	for j := 1; j < len(alphabet); j++ {
		go func(alphabetsl []byte) {
			for {
				s.Set(alphabetsl)
			}
		}(alphabet[:j])
	}

	for i := 0; i < len(settovar); i++ {
		go func(ab []byte) {
			for {
				s.Append(ab)
			}
		}(settovar[:i])
	}

	time.Sleep(5 * time.Second)
}

func BenchmarkBytes_Set(b *testing.B) {
	s := NewBytes()
	wg := &sync.WaitGroup{}
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func(wg *sync.WaitGroup) {
			for i := 0; i < b.N/10; i++ {
				s.Set([]byte("zz"))
			}
			wg.Done()
		}(wg)
	}
	wg.Wait()
}

func BenchmarkBytes_Get(b *testing.B) {
	s := NewBytes()
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

func BenchmarkBytes_Append(b *testing.B) {
	s := NewBytes()
	wg := &sync.WaitGroup{}
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func(wg *sync.WaitGroup) {
			for i := 0; i < b.N/10; i++ {
				s.Append([]byte("abc"))
			}
			wg.Done()
		}(wg)
	}
	wg.Wait()
}

package concurrent

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestByte(t *testing.T) {
	s := NewByte()

	s.Set('f')
	assert.Equal(t, byte('f'), s.Get())
	s.Set('4')
	assert.Equal(t, byte('4'), s.Get())
	s.Set(4)
	assert.Equal(t, byte(4), s.Get())
	s.Set('\n')
	assert.Equal(t, byte('\n'), s.Get())
	s.Set('g')
	assert.Equal(t, byte('g'), s.Get())
}

func TestByteConcurrent(t *testing.T) {
	s := NewByte()

	for i := 0; i < 100; i++ {
		go func() {
			for {
				s.Get()
			}
		}()
	}

	for i := 0; i < 10; i++ {
		go func(i byte) {
			for {
				s.Set(i)
			}
		}(byte(i))
	}

	time.Sleep(5 * time.Second)
}

func BenchmarkByte_Get(b *testing.B) {
	s := NewByte()
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

func BenchmarkByte_Set(b *testing.B) {
	s := NewByte()
	wg := &sync.WaitGroup{}
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func(wg *sync.WaitGroup) {
			for i := 0; i < b.N/10; i++ {
				s.Set(byte(i))
			}
			wg.Done()
		}(wg)
	}
	wg.Wait()
}

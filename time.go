package concurrent

import (
	"log"
	"time"
)

type actionTime int

const (
	actionSet actionTime = iota
	actionGet
	actionSince
)

type cmdTime struct {
	action     actionTime
	time       time.Time
	chTime     chan<- time.Time
	chDuration chan<- time.Duration
}

func NewTime() *Time {
	t := &Time{
		chCmd: make(chan cmdTime, 5),
	}
	return t.run()
}

type Time struct {
	time  time.Time
	chCmd chan cmdTime
}

func (s *Time) run() *Time {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("concurrent.(Time).run panicked:", r)
			}

			s.run()
		}()

		for c := range s.chCmd {
			switch c.action {
			case actionSet:
				s.time = c.time
			case actionGet:
				c.chTime <- s.time
			case actionSince:
				c.chDuration <- time.Since(s.time)
			}
		}
	}()
	return s
}

func (s *Time) Now() *Time {
	s.Set(time.Now())
	return s
}

func (s *Time) Set(t time.Time) *Time {
	s.chCmd <- cmdTime{
		action: actionSet,
		time:   t,
	}
	return s
}

func (s *Time) Get() time.Time {
	ch := make(chan time.Time)
	s.chCmd <- cmdTime{
		action: actionGet,
		chTime: ch,
	}
	return <-ch
}

func (s *Time) Since() time.Duration {
	ch := make(chan time.Duration)
	s.chCmd <- cmdTime{
		action:     actionSince,
		chDuration: ch,
	}
	return <-ch
}

func (s *Time) String() string {
	return s.Get().String()
}

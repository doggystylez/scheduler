package scheduler

import (
	"time"
)

type (
	Event interface {
		// event should implement Fire to be executed when timer expires
		Fire()
		// event should implement GetLabel to get a unique, reproducible identifier
		GetLabel() string
	}

	Schedule map[string]struct{}
)

func New() Schedule {
	return make(Schedule)
}

func (s Schedule) Add(e Event, t time.Time) {
	timer := time.NewTimer(time.Until(t))
	go func() {
		<-timer.C
		e.Fire()
		s.Remove(e)
	}()
	s[e.GetLabel()] = struct{}{}
}

func (s Schedule) Exists(e Event) bool {
	if _, ok := s[e.GetLabel()]; ok {
		return true
	} else {
		return false
	}
}

func (s Schedule) Remove(e Event) {
	delete(s, e.GetLabel())
}

package te

import (
	"errors"
	"sync"
	"time"
)

// Event represents a recurring calendar event.
type Event struct {
	Name string
	Time time.Time
	Expr Expression
}

// Schedule represents a set of named temporal expressions
// that emit recurring calendar events.
type Schedule struct {
	mu     sync.RWMutex
	exprs  map[string]Expression
	events chan Event
	quit   map[string]chan struct{}
	done   chan struct{}
}

// NewSchedule returns a new Schedule that will send
// recurring events on its channel.
func NewSchedule() *Schedule {
	return &Schedule{
		exprs:  make(map[string]Expression),
		events: make(chan Event),
		quit:   make(map[string]chan struct{}),
		done:   make(chan struct{}),
	}
}

// Get returns the expression associated with name.
func (s *Schedule) Get(name string) Expression {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.exprs == nil {
		return nil
	}
	expr, ok := s.exprs[name]
	if !ok {
		return nil
	}
	return expr
}

// Set sets the expression associated with name to expr.
// It replaces any existing expression.
func (s *Schedule) Set(name string, expr Expression) {
	if expr == nil {
		return
	}
	s.mu.Lock()
	s.del(name)
	s.exprs[name] = expr
	if !s.isClosed() {
		s.quit[name] = make(chan struct{})
		go s.watch(name, expr, s.quit[name])
	}
	s.mu.Unlock()
}

func (s *Schedule) watch(name string, expr Expression, quit chan struct{}) {
	for {
		now := time.Now()
		t := expr.Next(now)
		d := t.Sub(now)
		if d <= 0 {
			<-quit
			return
		}
		select {
		case <-time.After(d):
			s.events <- Event{Name: name, Time: t, Expr: expr}
		case <-quit:
			return
		}
	}
}

// Del deletes the expression associated with name.
func (s *Schedule) Del(name string) {
	s.mu.Lock()
	s.del(name)
	delete(s.exprs, name)
	s.mu.Unlock()
}

func (s *Schedule) del(name string) {
	quit, ok := s.quit[name]
	if !ok {
		return
	}
	quit <- struct{}{}
	delete(s.quit, name)
}

func (s *Schedule) Events() <-chan Event {
	return s.events
}

// Close stops all watchers and renders the schedule incapable
// of emitting further events.
//
// Close implements the io.Closer interface.
func (s *Schedule) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.isClosed() {
		return errors.New("te: schedule is already closed")
	}
	for name := range s.quit {
		s.del(name)
	}
	close(s.events)
	return nil
}

func (s *Schedule) isClosed() bool {
	select {
	case <-s.events:
		return true
	default:
		return false
	}
}

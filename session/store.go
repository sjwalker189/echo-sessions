package session

import (
	"errors"
	"sync"
)

type Store[T any] interface {
	Get(id string) (Session[T], error)
	Set(id string, sess Session[T]) error
	Del(id string) error

	// TODO: Touch should be throttled/and minimize network requests
	// i.e. calculate a frequency to update the session in the instance
	// that only the expires time will change
	// Touch(id string) error
	// RegenerateID(id string) (string, error)
	// Clear() error
}

type MemorySessionStore[T any] struct {
	mu    sync.Mutex
	store map[string]Session[T]
}

func NewMemorySessionStore[T any]() *MemorySessionStore[T] {
	return &MemorySessionStore[T]{
		store: make(map[string]Session[T]),
	}
}

func (s *MemorySessionStore[T]) Get(id string) (Session[T], error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	value, ok := s.store[id]
	if !ok {
		return NewSession[T](), errors.New("not found")
	}

	return value, nil
}

func (s *MemorySessionStore[T]) Set(id string, sess Session[T]) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store[id] = sess
	return nil
}

func (s *MemorySessionStore[T]) Del(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.store, id)
	return nil
}

func (s *MemorySessionStore[T]) Clear() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store = make(map[string]Session[T])
	return nil
}

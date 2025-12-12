package session

import (
	"errors"
	"sync"
)

type Store interface {
	Get(id string) (Session, error)
	Set(id string, sess Session) error
	Del(id string) error

	// TODO: Touch should be throttled/and minimize network requests
	// Touch(id string) error
	// RegenerateID(id string) (string, error)
	// Clear() error
}

type MemorySessionStore struct {
	mu    sync.Mutex
	store map[string]Session
}

func NewMemorySessionStore() *MemorySessionStore {
	return &MemorySessionStore{
		store: make(map[string]Session),
	}
}

func (s *MemorySessionStore) Get(id string) (Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	value, ok := s.store[id]
	if !ok {
		return New(), errors.New("not found")
	}

	return value, nil
}

func (s *MemorySessionStore) Set(id string, sess Session) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store[id] = sess
	return nil
}

func (s *MemorySessionStore) Del(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.store, id)
	return nil
}

func (s *MemorySessionStore) Clear() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store = make(map[string]Session)
	return nil
}

package lib

import "sync"

type Store struct {
	sync.RWMutex
	messages []string
}

func NewStore() *Store {
	return &Store{
		messages: make([]string, 0, 1024),
	}
}

func (s *Store) FindAll() []string {
	defer s.RUnlock()
	s.RLock()
	return s.messages
}

func (s *Store) Save(message string) {
	defer s.Unlock()
	s.Lock()
	s.messages = append(s.messages, message)
}

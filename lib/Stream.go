package lib

import "sync"

type Subscribers map[int]chan<- string

type Stream struct {
	sync.Mutex
	nextId      int
	subscribers Subscribers
}

func NewStream() *Stream {
	return &Stream{
		nextId:      0,
		subscribers: make(Subscribers),
	}
}

func (s *Stream) Subscribe() (<-chan string, int) {
	defer s.Unlock()
	s.Lock()
	s.nextId++
	id := s.nextId
	ch := make(chan string, 1)
	s.subscribers[id] = ch
	return ch, id
}

func (s *Stream) Unsubscribe(id int) {
	defer s.Unlock()
	s.Lock()
	delete(s.subscribers, id)
}

func (s *Stream) Fire(message string) {
	defer s.Unlock()
	s.Lock()
	for _, ch := range s.subscribers {
		ch <- message
	}
}

func (s *Stream) Shutdown() {
	s.Fire("")
}

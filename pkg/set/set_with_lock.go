package set

import "sync"

type withLock struct {
	container map[string]bool
	lock      sync.RWMutex
}

func NewSetWithLock() Set {
	return &withLock{container: make(map[string]bool)}
}

func (s *withLock) Add(values ...string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	for _, value := range values {
		s.container[value] = true
	}
}

func (s *withLock) Remove(value string) bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	_, exists := s.container[value]
	if exists {
		delete(s.container, value)
	}
	return exists
}

func (s *withLock) Has(value string) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	_, c := s.container[value]
	return c
}

func (s *withLock) Size() int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return len(s.container)
}

func (s *withLock) Clear() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.container = make(map[string]bool)
}

func (s *withLock) Items() []string {
	s.lock.RLock()
	defer s.lock.RUnlock()
	items := make([]string, 0)
	for i := range s.container {
		items = append(items, i)
	}
	return items
}

package set

type withoutLock struct {
	container map[string]bool
}

func NewSetWithoutLock() Set {
	return &withoutLock{container: make(map[string]bool)}
}

func (s *withoutLock) Add(values ...string) {
	for _, value := range values {
		s.container[value] = true
	}
}

func (s *withoutLock) Remove(value string) bool {
	_, exists := s.container[value]
	if exists {
		delete(s.container, value)
	}
	return exists
}

func (s *withoutLock) Has(value string) bool {
	_, c := s.container[value]
	return c
}

func (s *withoutLock) Size() int {
	return len(s.container)
}

func (s *withoutLock) Clear() {
	s.container = make(map[string]bool)
}

func (s *withoutLock) Items() []string {
	items := make([]string, 0)
	for i := range s.container {
		items = append(items, i)
	}
	return items
}

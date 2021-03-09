package pipe

import "context"

type InMemoryStorage struct{ Pipe }

func (s *InMemoryStorage) Load(context.Context) (Pipe, error) {
	return s.Pipe, nil
}

func (s *InMemoryStorage) Save(_ context.Context, pipe Pipe) error {
	s.Pipe = pipe

	return nil
}

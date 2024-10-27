package memory

import (
	"context"

	"github.com/Alma-media/elsa/model"
)

type Storage struct{ model.Pipe }

func (s *Storage) Load(context.Context) (model.Pipe, error) {
	return s.Pipe, nil
}

func (s *Storage) Save(_ context.Context, pipe model.Pipe) error {
	s.Pipe = pipe

	return nil
}

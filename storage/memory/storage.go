package memory

import (
	"context"

	"github.com/Alma-media/elsa/pipe"
)

type Storage struct{ pipe.Pipe }

func (s *Storage) Load(context.Context) (pipe.Pipe, error) {
	return s.Pipe, nil
}

func (s *Storage) Save(_ context.Context, pipe pipe.Pipe) error {
	s.Pipe = pipe

	return nil
}

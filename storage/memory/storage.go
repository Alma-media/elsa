package memory

import (
	"context"

	"github.com/Alma-media/elsa/flow"
)

type Storage struct{ flow.Pipe }

func (s *Storage) Load(context.Context) (flow.Pipe, error) {
	return s.Pipe, nil
}

func (s *Storage) Save(_ context.Context, pipe flow.Pipe) error {
	s.Pipe = pipe

	return nil
}

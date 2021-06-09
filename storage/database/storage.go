package database

import (
	"context"
	"database/sql"

	"github.com/Alma-media/elsa/pipe"
)

type PipeManager interface {
	Drop(tx *sql.Tx) error
	Load(tx *sql.Tx, pipe *pipe.Pipe) error
	Save(tx *sql.Tx, element pipe.Element) error
}

type Storage struct {
	db      *sql.DB
	manager PipeManager
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		db: db,
	}
}

func (storage *Storage) Load(ctx context.Context) (pipe.Pipe, error) {
	tx, err := storage.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	var pipe pipe.Pipe

	if err := storage.manager.Load(tx, &pipe); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return pipe, nil
}

func (storage *Storage) Save(ctx context.Context, pipe pipe.Pipe) error {
	tx, err := storage.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if err := storage.manager.Drop(tx); err != nil {
		return err
	}

	for _, element := range pipe {
		if err := storage.manager.Save(tx, element); err != nil {
			return err
		}
	}

	return tx.Commit()
}
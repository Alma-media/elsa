package sqlite

import (
	"context"
	"database/sql"
	"reflect"
	"testing"

	"github.com/Alma-media/elsa/pipe"
	"github.com/Alma-media/elsa/storage/database"
)

var (
	_ database.PipeManager = new(PipeManager)
)

func setup(t *testing.T) (*sql.Tx, func() error) {
	db, err := sql.Open("sqlite3", "file::memory:")
	if err != nil {
		t.Fatalf("unable to establish database connection: %s", err)
	}

	if err := Init(context.Background(), db); err != nil {
		t.Fatalf("database migration failure: %s", err)
	}

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("unable to start the transaction: %s", err)
	}

	return tx, func() error {
		if err := tx.Rollback(); err != nil {
			return err
		}

		return db.Close()
	}
}

func TestPipeManagerLoad(t *testing.T) {
	tx, release := setup(t)
	defer release()

	t.Run("load empty pipe", func(t *testing.T) {
		var pipe pipe.Pipe

		if err := new(PipeManager).Load(tx, &pipe); err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if len(pipe) != 0 {
			t.Error("pipe was expected to be empty")
		}
	})

	t.Run("load non-empty pipe", func(t *testing.T) {
		if _, err := tx.Exec(insertQuery, "foo", "bar", "count;reverse"); err != nil {
			t.Fatalf("failed to insert test data: %s", err)
		}

		if _, err := tx.Exec(insertQuery, "bar", "baz", nil); err != nil {
			t.Fatalf("failed to insert test data: %s", err)
		}

		var (
			actual   pipe.Pipe
			expected = pipe.Pipe{
				{
					BaseElement: pipe.BaseElement{
						Input:  "bar",
						Output: "baz",
					},
				},
				{
					BaseElement: pipe.BaseElement{
						Input:  "foo",
						Output: "bar",
						Pipe: []string{
							"count",
							"reverse",
						},
					},
					Processors: []pipe.Processor{
						pipe.Count,
						pipe.Reverse,
					},
				},
			}
		)

		if err := new(PipeManager).Load(tx, &actual); err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if len(actual) != len(expected) {
			t.Errorf("pipe was expected to contain %d elements", len(expected))
		}

		for index := range actual {
			if !reflect.DeepEqual(actual[index].BaseElement, expected[index].BaseElement) {
				t.Errorf(
					"the output \n%#v\ndoes not match expected\n%#v",
					actual[index].BaseElement,
					expected[index].BaseElement,
				)
			}
		}
	})
}

func TestPipeManagerSave(t *testing.T) {
	tx, release := setup(t)
	defer release()

	t.Run("save new routes", func(t *testing.T) {
		element := pipe.Element{
			BaseElement: pipe.BaseElement{
				Input:  "foo",
				Output: "bar",
				Pipe: []string{
					"print",
					"reverse",
				},
			},
		}

		if err := new(PipeManager).Save(tx, element); err != nil {
			t.Errorf("unexpected error: %s", err)
		}
	})
}

func TestPipeManagerDrop(t *testing.T) {
	tx, release := setup(t)
	defer release()

	t.Run("drop previous routes", func(t *testing.T) {
		if _, err := tx.Exec(insertQuery, "foo", "bar", "count;reverse"); err != nil {
			t.Fatalf("failed to insert test data: %s", err)
		}

		if _, err := tx.Exec(insertQuery, "bar", "baz", nil); err != nil {
			t.Fatalf("failed to insert test data: %s", err)
		}

		if err := new(PipeManager).Drop(tx); err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if err := tx.QueryRow(selectQuery).Scan(new(string), new(string), new(string)); err != sql.ErrNoRows {
			t.Errorf("unexpected error: %s", err)
		}
	})
}

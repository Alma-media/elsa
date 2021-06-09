package sqlite

import (
	"database/sql"
	"strings"

	"github.com/Alma-media/elsa/pipe"
	_ "github.com/mattn/go-sqlite3"
)

var (
	deleteQuery = `DELETE from route;`
	selectQuery = `SELECT input, output, pipe FROM route ORDER BY input, output;`
	insertQuery = `INSERT INTO route (input, output, pipe) VALUES(?, ?, ?);`
)

type PipeManager struct{}

func (PipeManager) Load(tx *sql.Tx, recv *pipe.Pipe) error {
	rows, err := tx.Query(selectQuery)
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var (
			functions sql.NullString
			element   pipe.Element
		)

		if err := rows.Scan(&element.Input, &element.Output, &functions); err != nil {
			return err
		}

		if functions.Valid {
			element.Pipe = strings.Split(functions.String, ";")
		}

		if err := element.Resolve(); err != nil {
			return err
		}

		*recv = append(*recv, element)
	}

	return nil
}

func (PipeManager) Save(tx *sql.Tx, element pipe.Element) error {
	_, err := tx.Exec(insertQuery, element.Input, element.Output, strings.Join(element.Pipe, ";"))

	return err
}

func (PipeManager) Drop(tx *sql.Tx) error {
	_, err := tx.Exec(deleteQuery)

	return err
}

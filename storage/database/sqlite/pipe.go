package sqlite

import (
	"database/sql"

	"github.com/Alma-media/elsa/flow"
	_ "github.com/mattn/go-sqlite3"
)

var (
	deleteQuery = `DELETE from route;`
	selectQuery = `SELECT input, output FROM route ORDER BY input, output;`
	insertQuery = `INSERT INTO route (input, output) VALUES(?, ?);`
)

type PipeManager struct{}

func (PipeManager) Load(tx *sql.Tx, recv *flow.Pipe) error {
	rows, err := tx.Query(selectQuery)
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var element flow.Element

		if err := rows.Scan(&element.Input, &element.Output); err != nil {
			return err
		}

		*recv = append(*recv, element)
	}

	return nil
}

func (PipeManager) Save(tx *sql.Tx, element flow.Element) error {
	_, err := tx.Exec(insertQuery, element.Input, element.Output)

	return err
}

func (PipeManager) Drop(tx *sql.Tx) error {
	_, err := tx.Exec(deleteQuery)

	return err
}

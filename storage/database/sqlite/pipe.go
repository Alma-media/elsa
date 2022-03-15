package sqlite

import (
	"database/sql"
	"encoding/json"

	"github.com/Alma-media/elsa/flow"
	_ "github.com/mattn/go-sqlite3"
)

var (
	deleteQuery = `DELETE from route;`
	selectQuery = `SELECT input, output, item FROM route ORDER BY input, output;`
	insertQuery = `INSERT INTO route (input, output, item) VALUES(?, ?, ?);`
)

type PipeManager struct{}

func (PipeManager) Load(tx *sql.Tx, recv *flow.Pipe) error {
	rows, err := tx.Query(selectQuery)
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var (
			item flow.Route
			data []byte
		)

		if err := rows.Scan(
			&item.Input.Path, &item.Output.Path, &data,
		); err != nil {
			return err
		}

		if err := json.Unmarshal(data, &item); err != nil {
			return err
		}

		*recv = append(*recv, item)
	}

	return rows.Err()
}

func (PipeManager) Save(tx *sql.Tx, item flow.Route) error {
	data, err := json.Marshal(item)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		insertQuery,
		item.Input.Path,
		item.Output.Path,
		data,
	)

	return err
}

func (PipeManager) Drop(tx *sql.Tx) error {
	_, err := tx.Exec(deleteQuery)

	return err
}

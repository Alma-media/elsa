package sqlite

import (
	"database/sql"
	"encoding/json"

	"github.com/Alma-media/elsa/model"
	_ "github.com/mattn/go-sqlite3"
)

var (
	deleteQuery = `DELETE from route;`
	selectQuery = `SELECT alias, input, output, options FROM route ORDER BY input, output;`
	insertQuery = `INSERT INTO route (alias, input, output, options) VALUES(?, ?, ?, ?);`
)

type PipeManager struct{}

func (PipeManager) Load(tx *sql.Tx, recv *model.Pipe) error {
	rows, err := tx.Query(selectQuery)
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var (
			element model.Element
			data    []byte
		)

		if err := rows.Scan(&element.Alias, &element.Input, &element.Output, &data); err != nil {
			return err
		}

		if err := json.Unmarshal(data, &element.Options); err != nil {
			return err
		}

		*recv = append(*recv, element)
	}

	return rows.Err()
}

func (PipeManager) Save(tx *sql.Tx, element model.Element) error {
	data, err := json.Marshal(element.Options)
	if err != nil {
		return err
	}

	_, err = tx.Exec(insertQuery, element.Alias, element.Input, element.Output, data)

	return err
}

func (PipeManager) Drop(tx *sql.Tx) error {
	_, err := tx.Exec(deleteQuery)

	return err
}

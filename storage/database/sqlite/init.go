package sqlite

import (
	"context"
	"database/sql"
)

var createTable = "CREATE TABLE IF NOT EXISTS \"route\" (" +
	"`input` text NOT NULL PRIMARY KEY," +
	"`output` text NOT NULL," +
	"`pipe` text," +
	"PRIMARY KEY (input, output)" +
	");"

func Init(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, createTable)

	return err
}

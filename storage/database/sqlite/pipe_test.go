package sqlite

import "github.com/Alma-media/elsa/storage/database"

var (
	_ database.PipeManager = new(PipeManager)
)

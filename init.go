package main

import (
	"fmt"
	"os"

	"github.com/remeh/mehtadata/db"
)

// Creates the SQLite schema.
func InitSchema(flags Flags) (bool, error) {
	// try to read in the env var if
	// everything is available
	// ----------------------

	// mandatory
	schema := os.Getenv("SCHEMA")

	if !StringsHasContent(schema) {
		fmt.Println(`Missing parameter.
Mandatory:
	SCHEMA      : schema.sql to init the database
		`)
	}

	return db.InitSchema(flags.DestSqlite, schema)
}

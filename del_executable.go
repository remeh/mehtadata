package main

import (
	"fmt"
	"os"

	"github.com/remeh/mehtadata/db"
)

func DelExecutable(flags Flags) (bool, error) {
	// mandatory
	filepath := os.Getenv("FILEPATH")
	platformName := os.Getenv("PLATFORM_NAME")

	ok := false

	if StringsHasContent(platformName, filepath) {
		ok = true
	}

	if !ok {
		fmt.Println(`Can't delete an executable.
Mandatory infos:
	FILEPATH      : filepath to the executable to delete
	PLATFORM_NAME : name of the platform containing this executable
		`)
		return false, fmt.Errorf("Missing fields.")
	}

	return db.DeleteExecutable(flags.DestSqlite, platformName, filepath)
}

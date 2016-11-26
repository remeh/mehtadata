package main

import (
	"fmt"
	"os"

	"github.com/remeh/mehtadata/db"
)

func DelPlatform(flags Flags) (bool, error) {
	// mandatory
	platformName := os.Getenv("PLATFORM_NAME")

	ok := false

	if StringsHasContent(platformName) {
		ok = true
	}

	if !ok {
		fmt.Println(`Can't delete an executable.
Mandatory infos:
	PLATFORM_NAME : name of the platform containing this executable
		`)
		return false, fmt.Errorf("Missing fields.")
	}

	return db.DeletePlatform(flags.DestSqlite, platformName)
}

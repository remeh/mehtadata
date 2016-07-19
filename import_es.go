package main

import (
	"db"
	"fmt"
	"log"
	"os"
	"strconv"

	"scraper"
)

// TODO(remy): could hit TheGamesDB for a real
// list of supported platforms.
func ImportES(flags Flags) error {
	var err error
	var platform int

	// parse parameters
	// ----------------------

	p := os.Getenv("PLATFORM_ID")
	gamelist := os.Getenv("GAMELIST")

	if !StringsHasContent(p, gamelist) {
		fmt.Println(`Missing fields to import the gamelist.xml file.
Mandatory infos:
	PLATFORM_ID : id of the platform into which will be inserted every entry of the gamelist.xml
	GAMELIST    : absolute path to the gamelist.xml file to import
		`)
		return fmt.Errorf("Missing fields.")
	}

	if platform, err = strconv.Atoi(p); err != nil {
		return fmt.Errorf("Wrong format for the platform ID")
	}

	// import
	// ----------------------

	file, err := os.Open(gamelist) // XXX(remy): needed ?
	if err != nil {
		return err
	}
	file.Close()

	gamesinfo, err := scraper.Decode(gamelist)
	if err != nil {
		return err
	}

	db.WriteDatabase(flags.DestSqlite, platform, &gamesinfo)

	log.Printf("%s imported.\n", gamelist)

	return nil
}

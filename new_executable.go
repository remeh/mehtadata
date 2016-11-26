package main

import (
	"fmt"
	"os"

	"github.com/remeh/mehtadata/db"
	"github.com/remeh/mehtadata/model"
)

func NewExecutable(flags Flags) (int64, bool, error) {
	// mandatory
	name := os.Getenv("NAME")
	filepath := os.Getenv("FILEPATH")
	platformName := os.Getenv("PLATFORM_NAME")

	// not mandatory
	description := os.Getenv("DESCRIPTION")
	genres := os.Getenv("GENRES")
	publisher := os.Getenv("PUBLISHER")
	developer := os.Getenv("DEVELOPER")
	releaseDate := os.Getenv("RELEASE_DATE")
	players := os.Getenv("PLAYERS")
	rating := os.Getenv("RATING")

	ok := false

	if StringsHasContent(platformName, name, filepath) {
		ok = true
	}

	if !ok {
		fmt.Println(`Can't create a new executable.
Mandatory infos:
	NAME          : name of the executable to create
	FILEPATH      : filepath to the executable to start
	PLATFORM_NAME : name of the platform containing this executable
Not mandatory:
	DESCRIPTION   : description of the executable
	GENRES        : genres of the executable
	PUBLISHER     : publisher of the executable
	DEVELOPER     : developer of the executable
	RELEASE_DATE  : release date of the executable
	PLAYERS       : players of the executable
	RATING        : rating of the executable
		`)
		return -1, false, fmt.Errorf("Missing fields.")
	}

	exec := model.Executable{
		Name:        name,
		Filepath:    filepath,
		Description: description,
		Genres:      genres,
		Publisher:   publisher,
		Developer:   developer,
		ReleaseDate: releaseDate,
		Players:     players,
		Rating:      rating,
	}

	return db.CreateExecutable(flags.DestSqlite, platformName, exec)
}

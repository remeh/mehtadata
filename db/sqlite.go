package db

import (
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/remeh/mehtadata/model"

	_ "github.com/mattn/go-sqlite3"
)

func InitSchema(database, filename string) (bool, error) {

	// database
	// ----------------------

	var err error

	db, err := sql.Open("sqlite3", database)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// read the schema
	// ----------------------

	var data []byte

	if data, err = ioutil.ReadFile(filename); err != nil {
		return false, err
	}

	// remove all line return

	data = bytes.Replace(data, []byte("\n"), []byte(" "), -1)
	data = bytes.Replace(data, []byte("\r"), []byte(" "), -1) // windows user having edited the file

	queries := strings.Split(string(data), ";")

	// prepares the transaction
	// ----------------------

	tx, err := db.Begin()
	if err != nil {
		return false, err
	}

	// execute each line
	// ----------------------

	for _, query := range queries {
		if _, err := tx.Exec(query); err != nil {
			return false, err
		}
	}

	// commit
	// ----------------------

	err = tx.Commit()

	return err == nil, err
}

// CreatePlatform creates an empty platform in the given sqlite database.
// Returns the ID of the newly created platform.
func CreatePlatform(database string, platform model.Platform) (int64, error) {

	// database
	// ----------------------

	db, err := sql.Open("sqlite3", database)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// prepares the transaction
	// ----------------------

	tx, err := db.Begin()
	if err != nil {
		return -1, err
	}

	execStmt, err := tx.Prepare(`insert into "platform" ("name", "command", "icon", "background", "type", "discover_dir", "discover_ext") values(?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return -1, err
	}

	// executes the query
	// ----------------------

	result, err := execStmt.Exec(platform.Name, platform.Command, platform.Icon, platform.Background, platform.Type, platform.DiscoverDir, platform.DiscoverExts)
	tx.Commit()

	if err != nil {
		return -1, fmt.Errorf("[err] Can't create a new platform %s in the DB: %s", platform.Name, err.Error())
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("[err] Can't create a new platform %s in the DB (id retrieving): %s", platform.Name, err.Error())
	}

	return id, nil
}

// writeDatabase writes the result of the scraping into the given database.
func WriteDatabase(database string, platform int, gamesInfo *model.Gamesinfo) {
	db, err := sql.Open("sqlite3", database)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Prepares the transaction
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	execStmt, err := tx.Prepare("insert into executable (display_name, filepath, platform_id, description, genres, developer, publisher, release_date, players, rating) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	execResStmt, err := tx.Prepare("insert into executable_resource (executable_id, type, filepath) values(?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	for idx := range gamesInfo.Games {
		gameInfo := gamesInfo.Games[idx]
		writeGame(db, platform, execStmt, execResStmt, gameInfo)
	}

	// Ends the transaction
	tx.Commit()
}

// writeGame writes one game in the DB.
func writeGame(db *sql.DB, platform int, execStmt *sql.Stmt, execResStmt *sql.Stmt, gameInfo model.Gameinfo) {
	// Entry in executable
	rating := fmt.Sprintf("%2.1f", gameInfo.Rating)
	result, err := execStmt.Exec(gameInfo.Title, gameInfo.Filepath, platform, gameInfo.Description, gameInfo.Genres, gameInfo.Developer, gameInfo.Publisher, gameInfo.ReleaseDate, gameInfo.Players, rating)
	if err != nil {
		log.Printf("[err] Can't write the info of %s in the DB: %s", gameInfo.Title, err.Error())
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("[err] Can't retrieve the last insert id for %s : %s", gameInfo.Title, err.Error())
		return
	}

	// Entries in executable_resource
	// Cover
	if len(gameInfo.CoverPath) != 0 {
		result, err = execResStmt.Exec(id, "cover", gameInfo.CoverPath)
		if err != nil {
			log.Printf("[err] Can't write the cover of %s in the DB: %s", gameInfo.Title, err.Error())
		}
	}
	// Logos
	if len(gameInfo.LogoPaths) != 0 {
		for idx := range gameInfo.LogoPaths {
			logoPath := gameInfo.LogoPaths[idx]
			result, err = execResStmt.Exec(id, "logo", logoPath)
			if err != nil {
				log.Printf("[err] Can't write a logo of %s in the DB: %s", gameInfo.Title, err.Error())
			}
		}
	}
	// Screenshots
	if len(gameInfo.ScreenshotPaths) != 0 {
		for idx := range gameInfo.ScreenshotPaths {
			screenshotPath := gameInfo.ScreenshotPaths[idx]
			result, err = execResStmt.Exec(id, "screenshot", screenshotPath)
			if err != nil {
				log.Printf("[err] Can't write a screenshot of %s in the DB: %s", gameInfo.Title, err.Error())
			}
		}
	}
	// Fanarts
	if len(gameInfo.FanartPaths) != 0 {
		for idx := range gameInfo.FanartPaths {
			fanartPath := gameInfo.FanartPaths[idx]
			result, err = execResStmt.Exec(id, "fanart", fanartPath)
			if err != nil {
				log.Printf("[err] Can't write a fanart of %s in the DB: %s", gameInfo.Title, err.Error())
			}
		}
	}
}

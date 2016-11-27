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
func CreatePlatform(database string, platform model.Platform) (int64, bool, error) {

	// database
	// ----------------------

	db, err := sql.Open("sqlite3", database)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// ensure that no platform with the same name
	// is already existing
	// ----------------------

	var i int
	if err := db.QueryRow(`
		SELECT count(*)
		FROM "platform"
		WHERE
			name = ?
	`, platform.Name).Scan(&i); err != nil {
		return -1, false, err
	}

	exists := false
	if i >= 1 {
		exists = true
	}

	// executes the query
	// ----------------------

	var result sql.Result

	if !exists {
		result, err = db.Exec(`insert into "platform" ("name", "command", "icon", "background", "type", "discover_dir", "discover_ext") values(?, ?, ?, ?, ?, ?, ?)`, platform.Name, platform.Command, platform.Icon, platform.Background, platform.Type, platform.DiscoverDir, platform.DiscoverExts)
	} else {
		result, err = db.Exec(`
			UPDATE "platform"
			SET
				"command" = ?,
				"icon" = ?,
				"background" = ?,
				"type" = ?,
				"discover_dir" = ?,
				"discover_ext" = ?
			WHERE
				"name" = ?
		`, platform.Command, platform.Icon, platform.Background, platform.Type, platform.DiscoverDir, platform.DiscoverExts, platform.Name)
	}

	if err != nil {
		return -1, false, fmt.Errorf("[err] Can't create a new platform %s in the DB: %s", platform.Name, err.Error())
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, false, fmt.Errorf("[err] Can't create a new platform %s in the DB (id retrieving): %s", platform.Name, err.Error())
	}

	return id, exists, nil
}

// DeletePlatform deletes the given platform by its name and
// all its executables.
func DeletePlatform(database, platformName string) (bool, error) {

	// database
	// ----------------------

	var err error

	db, err := sql.Open("sqlite3", database)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// ensure that the platform with the given name exists
	// and retrieve its ID.
	// ----------------------

	var platformId int
	if err := db.QueryRow(`
		SELECT "id"
		FROM "platform"
		WHERE
			name = ?
	`, platformName).Scan(&platformId); err == sql.ErrNoRows {
		return false, fmt.Errorf("Unknown platform")
	} else if err != nil {
		return false, err
	}

	// delete all executables resources
	// ----------------------

	var result sql.Result

	if result, err = db.Exec(`
		DELETE FROM "executable_resource"
		WHERE "executable_id" IN (
			SELECT "id" FROM "executable"
			WHERE
				"platform_id" = ?
		)
	`, platformId); err != nil {
		return false, err
	}

	// delete all its executables
	// ----------------------

	if result, err = db.Exec(`
		DELETE FROM "executable"
		WHERE
			"platform_id" = ?
	`, platformId); err != nil {
		return false, err
	}

	// delete the platform
	// ----------------------

	if result, err = db.Exec(`
		DELETE FROM "platform"
		WHERE
			"id" = ?
	`, platformId); err != nil {
		return false, err
	}

	// checks the deletion
	// ----------------------

	c, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("while checking for delection: %v", err.Error())
	}

	return c == 1, nil
}

func CreateResource(database, resource, filepath, platformName, typ string) (bool, error) {
	// database
	// ----------------------

	var err error

	db, err := sql.Open("sqlite3", database)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// gets the executable ID
	// ----------------------

	var executableId int
	if err = db.QueryRow(`
		SELECT "executable"."id" FROM "executable"
		JOIN "platform"
			ON "platform"."name" = ?
			AND "executable"."platform_id" = "platform"."id"
		WHERE
			"executable"."filepath" = ?
		LIMIT 1
	`, platformName, filepath).Scan(&executableId); err == sql.ErrNoRows {
		return false, fmt.Errorf("Unknown executable")
	} else if err != nil {
		return false, err
	}

	// insert the resource
	// ----------------------

	var result sql.Result

	if result, err = db.Exec(`
		INSERT INTO "executable_resource"
		("executable_id", "type", "filepath")
		VALUES
		(?, ?, ?)
	`, executableId, typ, resource); err != nil {
		return false, err
	}

	// checks the insertion
	// ----------------------

	c, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("while checking for insertion: %v", err.Error())
	}

	return c >= 1, nil
}

// DeleteExecutable deletes the executable from the database.
func DeleteExecutable(database, platformName, filepath string) (bool, error) {

	// database
	// ----------------------

	var err error

	db, err := sql.Open("sqlite3", database)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// ensure that the platform with the given name exists
	// and retrieve its ID.
	// ----------------------

	var platformId int
	if err := db.QueryRow(`
		SELECT "id" 
		FROM "platform"
		WHERE
			name = ?
	`, platformName).Scan(&platformId); err == sql.ErrNoRows {
		return false, fmt.Errorf("Unknown platform")
	} else if err != nil {
		return false, err
	}

	// delete all executables resources
	// ----------------------

	var result sql.Result

	if result, err = db.Exec(`
		DELETE FROM "executable_resource"
		WHERE "executable_id" IN (
			SELECT "id" FROM "executable"
			WHERE
				"executable"."filepath" = ?
				AND
				"platform_id" = ?
		)
	`, filepath, platformId); err != nil {
		return false, err
	}

	// delete the executable
	// ----------------------

	if result, err = db.Exec(`
		DELETE FROM "executable"
		WHERE
			"filepath" = ?
			AND
			"platform_id" = ?
	`, filepath, platformId); err != nil {
		return false, err
	}

	// checks the deletion
	// ----------------------

	c, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("while checking for delection: %v", err.Error())
	}

	return c >= 1, nil
}

// CreateExecutable creates the executable in the database.
// Returns the ID of the created executable.
func CreateExecutable(database string, platformName string, executable model.Executable) (int64, bool, error) {

	// database
	// ----------------------

	db, err := sql.Open("sqlite3", database)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// ensure that the platform with the given name exists
	// and retrieve its ID.
	// ----------------------

	var platformId int
	if err := db.QueryRow(`
		SELECT "id" 
		FROM "platform"
		WHERE
			name = ?
	`, platformName).Scan(&platformId); err == sql.ErrNoRows {
		return -1, false, fmt.Errorf("Unknown platform")
	} else if err != nil {
		return -1, false, err
	}

	// checks whether this executable already exists
	// ----------------------

	var c int
	if err := db.QueryRow(`
		SELECT COUNT(*)
		FROM "executable"
		JOIN "platform"
			ON "platform"."id" = "executable"."platform_id"
		WHERE
			"platform"."name" = ?
			AND
			"executable"."filepath" = ?
	`, platformName, executable.Filepath).Scan(&c); err != nil {
		return -1, false, err
	}

	exists := false
	if c != 0 {
		exists = true
	}

	// executes the query
	// ----------------------

	var result sql.Result

	if !exists {
		result, err = db.Exec(`
		INSERT INTO "executable"
		("display_name", "filepath", "platform_id", "description", "genres", "publisher", "developer", "release_date", "players", "rating")
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, executable.Name, executable.Filepath, platformId, executable.Description, executable.Genres, executable.Publisher, executable.Developer, executable.ReleaseDate, executable.Players, executable.Rating)
	} else {
		result, err = db.Exec(`
		UPDATE "executable"
		SET
			"display_name" = ?,
			"description" = ?,
			"genres" = ?,
			"publisher" = ?,
			"developer" = ?,
			"release_date" = ?,
			"players" = ?,
			"rating" = ?
		WHERE
			"executable"."platform_id" = ?
			AND
			"executable"."filepath" = ?
	`, executable.Name, executable.Description, executable.Genres, executable.Publisher, executable.Developer, executable.ReleaseDate, executable.Players, executable.Rating, platformId, executable.Filepath)
	}

	if err != nil {
		return -1, false, fmt.Errorf("[err] Can't create/update a new executable %s in the DB: %s", executable.Name, err.Error())
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, false, fmt.Errorf("[err] During id retrieving of '%s': %s", executable.Name, err.Error())
	}

	return id, exists, nil

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

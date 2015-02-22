package db

import (
	"database/sql"
	"log"

	"model"

	_ "github.com/mattn/go-sqlite3"
)

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
	execStmt, err := tx.Prepare("insert into executable (display_name, filepath, platform_id) values(?, ?, ?)")
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
	result, err := execStmt.Exec(gameInfo.Title, gameInfo.Filepath, platform)
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
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"common"
	"db"
	"model"
	"thegamesdb"
)

// Scraping launches the scraping for
func Scraping(flags Flags) (int, error) {
	// try to read in the env var if
	// everything is available
	// ----------------------

	// mandatory
	platformName := os.Getenv("PLATFORM")
	p := os.Getenv("PLATFORM_ID")
	output := os.Getenv("OUTPUT")
	dir := os.Getenv("DIR")
	exts := os.Getenv("EXTS")
	w := os.Getenv("WIDTH")

	if !StringsHasContent(platformName, p, output, exts) {
		fmt.Println(`Missing parameter.
Mandatory:
	PLATFORM    : name of the platform in TheGamesDB to find data. Use -platforms for a list
	PLATFORM_ID : id of the platform to scrape for
	DIR         : directory containing the executables
	EXTS        : extensions of the executable files, separated with a comma.
	OUTPUT      : output directory for images, etc.
Optional:
	WIDTH       : max width for the downloaded content (default: 768)
		`)
	}

	// parse parameters
	// ----------------------

	if len(output) > 0 && string(output[len(output)-1]) != "/" {
		output = output + "/"
	}
	if len(dir) > 0 && string(dir[len(dir)-1]) != "/" {
		dir = dir + "/"
	}

	var err error
	platformId := -1
	width := 768

	if platformId, err = strconv.Atoi(p); err != nil {
		fmt.Errorf("Bad platform value.")
		os.Exit(-1)
	}
	if len(w) > 0 {
		if width, err = strconv.Atoi(w); err != nil {
			fmt.Errorf("Unparseable width value.")
			os.Exit(-1)
		}
	}

	fmt.Printf("Launch scraping %d %d", platformId, width)

	return scrape(flags, platformName, platformId, dir, exts, output, width)
}

func scrape(flags Flags, platformName string, platformId int, dir, extensions, output string, width int) (int, error) {
	// Extensions array
	split := strings.Split(extensions, ",")
	exts := make([]string, len(split))
	for i, v := range split {
		exts[i] = strings.Trim(v, " ")
	}

	// Platforms array
	split = strings.Split(platformName, ",")
	platforms := make([]string, len(split))
	for i, v := range split {
		platforms[i] = strings.Trim(v, " ")
	}

	// look for files to proceed in the given directory if any
	var filenames []string
	if len(dir) > 0 {
		filenames = lookForFiles(dir, exts)
	} else {
		if len(flag.Args()) == 0 {
			fmt.Println("You should either use the DIR environment variable to provide a directory or provide filepath when calling mehtadata.")
			os.Exit(1)
		}
		filenames = flag.Args()
	}

	gamesinfo := model.NewGamesinfo()
	client := thegamesdb.NewClient()

	for _, filename := range filenames {
		gameinfo, err := client.Find(filename, platforms, dir, output, uint(width))
		if err != nil {
			log.Println("[err] Unable to find info for the game:", filename)
			log.Println(err)
			continue
		}

		// game scraped.
		if len(gameinfo.Title) > 0 {
			fmt.Printf("For '%s', scraped : '%s' on '%s'\n", filename, gameinfo.Title, gameinfo.Platform)
		} else {
			common.FillDefaults(dir, filename, &gameinfo)
			fmt.Printf("Nothing found for '%s'\n", filename)
		}
		gamesinfo.AddGame(gameinfo)
	}

	// writes executables info
	db.WriteDatabase(flags.DestSqlite, platformId, gamesinfo)
	return 0, nil
}

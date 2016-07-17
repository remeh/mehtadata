package main

import (
	"fmt"
	"os"
	"strconv"
)

// Scraping launches the scraping for
func Scraping(flags Flags) (int64, error) {
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

	if !StringsHasContent(platformName, p, output, dir, exts) {
		fmt.Println(`
			Missing parameter.
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

	// TODO(remy): launch the scraping.
	return 0, nil
}

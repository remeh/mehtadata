package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Options given to the CLI
type Flags struct {
	DestSqlite string // If destination is mehstation, the mehstation database write in.

	ShowPlatforms bool // To show the list of available platforms.
	InputGamelist bool // gamelist.xml file
	Scrape        bool
	NewPlatform   bool // name of a new platform to create
}

// ParseFlags parses the CLI options.
func ParseFlags() Flags {
	flags := Flags{}

	flag.StringVar(&(flags.DestSqlite), "meh-db", "database.db", "If destination is mehstation, the mehstation database write in.")

	flag.BoolVar(&(flags.InputGamelist), "import-es", false, "Import an EmulationStation gamelist.xml file.")
	flag.BoolVar(&(flags.ShowPlatforms), "show-platforms", false, "Display all TheGamesDB supported platforms")
	flag.BoolVar(&(flags.NewPlatform), "new-platform", false, "To create a new platform.")
	flag.BoolVar(&(flags.Scrape), "scrape", false, "To scrape content for a platform.")

	flag.Parse()

	return flags
}

func lookForFiles(directory string, extensions []string) []string {
	results := make([]string, 0)

	// list files in the directory
	fileinfos, err := ioutil.ReadDir(directory)
	if err != nil {
		return results
	}

	// for every files existing in the directory
	for _, fileinfo := range fileinfos {
		// don't mind of directories and check that the extension is valid for this scrape session.
		name := fileinfo.Name()
		if !fileinfo.IsDir() {
			// Check extensions
			extension := strings.ToLower(filepath.Ext(name))
			for _, e := range extensions {
				if extension == strings.ToLower(e) {
					results = append(results, name)
					break
				}
			}
		}
	}

	return results
}

func main() {
	flags := ParseFlags()

	// Show platforms mode
	// ----------------------

	if flags.ShowPlatforms {
		ShowPlatforms()
		os.Exit(0)
	}

	// Create platform
	// ----------------------

	if flags.NewPlatform {
		if _, err := NewPlatform(flags); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Println("Platform created.")
		os.Exit(0)
	}

	// Scrape mode
	// ----------------------
	if flags.Scrape {
		var amount int
		var err error

		if amount, err = Scraping(flags); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		fmt.Println("Content for %d scraped.", amount)
		os.Exit(0)
	}

	// Import from EmulationStation mode.
	// ----------------------

	if flags.InputGamelist {
		if err := ImportES(flags); err != nil {
			fmt.Println("Error:", err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	}

	flag.PrintDefaults()
}

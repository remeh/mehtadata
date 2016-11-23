package main

import (
	"flag"
	"fmt"
	"os"
)

// Options given to the CLI
type Flags struct {
	DestSqlite string // If destination is mehstation, the mehstation database write in.

	ShowPlatforms bool // To show the list of available platforms.
	InputGamelist bool // gamelist.xml file
	Scrape        bool
	NewPlatform   bool // name of a new platform to create
	InitSchema    bool
}

// ParseFlags parses the CLI options.
func ParseFlags() Flags {
	flags := Flags{}

	flag.StringVar(&(flags.DestSqlite), "db", "database.db", "If destination is mehstation, the mehstation database write in.")

	flag.BoolVar(&(flags.InputGamelist), "import-es", false, "Import an EmulationStation gamelist.xml file.")
	flag.BoolVar(&(flags.ShowPlatforms), "show-platforms", false, "Display all TheGamesDB supported platforms")
	flag.BoolVar(&(flags.NewPlatform), "new-platform", false, "To create a new platform.")
	flag.BoolVar(&(flags.Scrape), "scrape", false, "To scrape content for a platform.")
	flag.BoolVar(&(flags.InitSchema), "init", false, "Init the schema")

	flag.Parse()

	return flags
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

	if flags.InitSchema {
		if done, err := InitSchema(flags); err != nil {
			fmt.Println("Error:", err.Error())
			os.Exit(1)
		} else {
			fmt.Println("Init schema:", done)
			os.Exit(1)
		}
	}

	flag.PrintDefaults()
}

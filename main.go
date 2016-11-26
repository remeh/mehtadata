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
	Scrape        bool // to launch the scraper
	NewPlatform   bool // new platform
	DelPlatform   bool // delete a platform
	NewExecutable bool // new executbale
	DelExecutable bool // delete an executable
	InitSchema    bool // to run a .sql file
}

// ParseFlags parses the CLI options.
func ParseFlags() Flags {
	flags := Flags{}

	flag.StringVar(&(flags.DestSqlite), "db", "database.db", "If destination is mehstation, the mehstation database write in.")

	flag.BoolVar(&(flags.InputGamelist), "import-es", false, "Import an EmulationStation gamelist.xml file.")
	flag.BoolVar(&(flags.ShowPlatforms), "show-platforms", false, "Display all TheGamesDB supported platforms")
	flag.BoolVar(&(flags.NewPlatform), "new-platform", false, "To create a new platform.")
	flag.BoolVar(&(flags.DelPlatform), "del-platform", false, "Delete a platform and all its executables")
	flag.BoolVar(&(flags.DelExecutable), "del-exec", false, "To delete an executable")
	flag.BoolVar(&(flags.NewExecutable), "new-exec", false, "To create a new executable.")
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
		if _, exists, err := NewPlatform(flags); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		} else {
			if !exists {
				fmt.Println("Platform created")
			} else {
				fmt.Println("Platform updated")
			}
		}
		os.Exit(0)
	}

	// Create executable
	// ----------------------

	if flags.NewExecutable {
		if _, existing, err := NewExecutable(flags); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		} else {
			if !existing {
				fmt.Println("Executable created")
			} else {
				fmt.Println("Executable updated")
			}
		}
		os.Exit(0)
	}

	// Delete a platform
	// ----------------------

	if flags.DelPlatform {
		if deleted, err := DelPlatform(flags); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		} else {
			if !deleted {
				fmt.Println("Platform not deleted")
			} else {
				fmt.Println("Platform deleted")
			}
		}
		os.Exit(0)
	}

	// Delete an executable
	// ----------------------

	if flags.DelExecutable {
		if deleted, err := DelExecutable(flags); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		} else {
			if !deleted {
				fmt.Println("Executable not deleted")
			} else {
				fmt.Println("Executable deleted")
			}
		}
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

	// Run .SQL files
	// ----------------------

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

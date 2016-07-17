package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"common"
	"db"
	"thegamesdb"
)

// Options given to the CLI
type Flags struct {
	DestSqlite    string // If destination is mehstation, the mehstation database write in.
	ShowPlatforms bool   // To show the list of available platforms.
	InputGamelist string // gamelist.xml file

	Scrape      bool
	NewPlatform bool // name of a new platform to create
}

// ParseFlags parses the CLI options.
func ParseFlags() Flags {
	flags := Flags{}

	flag.StringVar(&(flags.DestSqlite), "meh-db", "database.db", "If destination is mehstation, the mehstation database write in.")
	flag.BoolVar(&(flags.ShowPlatforms), "platforms", false, "Display all the available platforms")
	flag.StringVar(&(flags.InputGamelist), "es", "", "gamelist.xml to import (Import from EmulationStation mode)")

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

func printPlatforms() {
	for _, p := range thegamesdb.TGDBPlatforms {
		fmt.Printf("%s\n", p.Name)
	}
}

func importFromEmulationStation(flags Flags) {
	if len(flags.DestSqlite) == 0 {
		fmt.Printf("Parameter error:\nWith the -es flag (import from emulation station), you'll need to\nprovida a meh-db pointing to the 'database.db'\nand a meh-platform value pointing to the dest platform.")
		os.Exit(1)
	}

	file, err := os.Open(flags.InputGamelist)
	if err != nil {
		log.Println("[err]", err.Error())
		os.Exit(1)
	}
	file.Close()

	gamesinfo, err := common.Decode(flags.InputGamelist)
	if err != nil {
		log.Println("[err]", err.Error())
		os.Exit(1)
	}

	// TODO(remy): dest platform must be an env var
	destPlatform := 1

	db.WriteDatabase(flags.DestSqlite, destPlatform, &gamesinfo)

	log.Printf("%s imported.\n", flags.InputGamelist)

	os.Exit(0)
}

func main() {
	flags := ParseFlags()

	// Show platforms mode
	// ----------------------

	if flags.ShowPlatforms {
		printPlatforms()
		os.Exit(0)
	}

	// Create platform
	// ----------------------

	if flags.NewPlatform {
		if _, err := NewPlatform(flags); err != nil {
			fmt.Println(err)
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
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Content for %d scraped.", amount)
		os.Exit(0)
	}

	// Import from EmulationStation mode.

	if len(flags.InputGamelist) > 0 {
		importFromEmulationStation(flags)
	}

}

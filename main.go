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
	"model"
	"thegamesdb"
)

// Options given to the CLI
type Flags struct {
	InputDirectory  string // in which directory to look for games
	OutputDirectory string // in which directory outputing the resulting files (images) and gamelist.xml
	Destination     string // Possible values: "emulationstation", "mehstation"
	DestSqlite      string // If destination is mehstation, the mehstation database write in.
	DestPlatform    int    // If destination is mehstation, the mehstation platform id to write for.
	Extension       string // extension to look for, separated by a space
	Platform        string // Which platform we must use for the scraping
	MaxWidth        uint   // Max width of the cover
	ShowPlatforms   bool   // To show the list of available platforms.
}

// ParseFlags parses the CLI options.
func ParseFlags() Flags {
	flags := Flags{}

	flag.StringVar(&(flags.InputDirectory), "in-dir", "", "Input directories (directory containing games)")
	flag.StringVar(&(flags.OutputDirectory), "out-dir", "", "Output directories for images and cover")
	flag.StringVar(&(flags.Destination), "dest", "emulationstation", "Possible values: emulationstation, mehstation")
	flag.StringVar(&(flags.DestSqlite), "meh-db", "database.db", "If destination is mehstation, the mehstation database write in.")
	flag.IntVar(&(flags.DestPlatform), "meh-platform", -1, "If destination is mehstation, the mehstation database write in.")
	flag.StringVar(&(flags.Extension), "ext", ".zip,.rar", "Accepted extensions")
	flag.StringVar(&(flags.Platform), "p", "", "Platform to use for the scraping")
	flag.UintVar(&(flags.MaxWidth), "w", 768, "Max width for the downloaded cover")
	flag.BoolVar(&(flags.ShowPlatforms), "platforms", false, "Display all the available platforms")

	flag.Parse()

	if flags.Destination != "emulationstation" && flags.Destination != "mehstation" {
		log.Fatalf("Unsupported destination : %s\nPossible values: \"emulationstation\", \"mehstation\"\n", flags.Destination)
	}

	if len(flags.InputDirectory) > 0 && string(flags.InputDirectory[len(flags.InputDirectory)-1]) != "/" {
		flags.InputDirectory = flags.InputDirectory + "/"
	}
	if len(flags.OutputDirectory) > 0 && string(flags.OutputDirectory[len(flags.OutputDirectory)-1]) != "/" {
		flags.OutputDirectory = flags.OutputDirectory + "/"
	}

	if flags.Destination == "mehstation" && flags.DestPlatform == -1 {
		log.Fatalf("To write in a mehstation DB, mehtadata needs the platform ID for which it will scrape metadata.")
	}

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
					// removes the extension
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

func main() {
	var err error
	flags := ParseFlags()

	if flags.ShowPlatforms {
		printPlatforms()
		os.Exit(0)
	}

	// Extensions array
	exts := make([]string, 0)
	split := strings.Split(flags.Extension, ",")
	for _, v := range split {
		exts = append(exts, strings.Trim(v, " "))
	}

	filenames := lookForFiles(flags.InputDirectory, exts)
	gamesinfo := model.NewGamesinfo()

	var file *os.File // For gamelist.xml for ES 2.0

	// Create/lock the gamelist.xml for EmulationStation 2.0
	// Do it now because it's useless to scrap everything if
	// we can't write the result into gamelist.xml
	if flags.Destination == "emulationstation" {
		file, err = os.Create(flags.OutputDirectory + "gamelist.xml")

		if err != nil {
			log.Fatalf("[err] Can't create %sgamelist.xml: %s\n", flags.OutputDirectory, err.Error())
		}
	} else if flags.Destination == "mehstation" {

	}

	client := thegamesdb.NewClient()
	for _, filename := range filenames {
		gameinfo, err := client.Find(filename, flags.Platform, flags.OutputDirectory, flags.MaxWidth)
		if err != nil {
			log.Println("[err] Unable to find info for the game:", filename)
			continue
		}

		// game scraped.
		fmt.Printf("For '%s', scraped : '%s' on '%s'\n", filename, gameinfo.Title, gameinfo.Platform)
		gamesinfo.AddGame(gameinfo)
	}

	if flags.Destination == "emulationstation" {
		// builds the gamelist.xml
		data, err := common.Encode(gamesinfo)
		if err != nil {
			log.Println("[err] Failed to encode the final xml:", err.Error())
			os.Exit(1)
		}

		file.Write(data)
		file.Close()

		log.Println("gamelist.xml written.")
	} else if flags.Destination == "mehstation" {
		db.WriteDatabase(flags.DestSqlite, flags.DestPlatform, gamesinfo)
	}
}

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
	"model"
	"thegamesdb"
)

// Options given to the CLI
type Flags struct {
	Input         string // in which directory to look for games
	Output        string // in which directory outputing the resulting files (images) and gamelist.xml
	Extension     string // extension to look for, separated by a space
	Platform      string // Which platform we must use for the scraping
	MaxWidth      uint   // Max width of the cover
	ShowPlatforms bool   // To show the list of available platforms.
}

// ParseFlags parses the CLI options.
func ParseFlags() Flags {
	flags := Flags{}

	flag.StringVar(&(flags.Input), "in", "", "Input directories (directory containing games)")
	flag.StringVar(&(flags.Output), "out", "", "Output directories for images and cover")
	flag.StringVar(&(flags.Extension), "ext", ".zip,.rar", "Accepted extensions")
	flag.StringVar(&(flags.Platform), "p", "", "Platform to use for the scraping")
	flag.UintVar(&(flags.MaxWidth), "w", 768, "Max width for the downloaded cover")
	flag.BoolVar(&(flags.ShowPlatforms), "platforms", false, "Display all the available platforms")

	flag.Parse()

	if len(flags.Input) > 0 && string(flags.Input[len(flags.Input)-1]) != "/" {
		flags.Input = flags.Input + "/"
	}
	if len(flags.Output) > 0 && string(flags.Output[len(flags.Output)-1]) != "/" {
		flags.Output = flags.Output + "/"
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

	filenames := lookForFiles(flags.Input, exts)
	gamesinfo := model.NewGamesinfo()

	// Create/lock the gamelist.xml for EmulationStation 2.0
	// Do it now because it's useless to scrap everything if
	// we can't write the result into gamelist.xml
	file, err := os.Create(flags.Output + "gamelist.xml")

	if err != nil {
		log.Printf("[err] Can't create %sgamelist.xml: %s\n", err.Error())
		os.Exit(1)
	}

	client := thegamesdb.NewClient()
	for _, filename := range filenames {
		gameinfo, err := client.Find(filename, flags.Platform, flags.Output, flags.MaxWidth)
		if err != nil {
			log.Println("[err] Unable to find info for the game:", filename)
			continue
		}

		// game scraped.
		fmt.Printf("For '%s', scraped : '%s' on '%s'\n", filename, gameinfo.Title, gameinfo.Platform)
		gamesinfo.AddGame(gameinfo)
	}

	// builds the gamelist.xml
	data, err := common.Encode(gamesinfo)
	if err != nil {
		log.Println("[err] Failed to encode the final xml:", err.Error())
		os.Exit(1)
	}

	file.Write(data)
	file.Close()

	log.Println("gamelist.xml written.")
}

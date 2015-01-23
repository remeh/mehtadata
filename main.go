package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"common"
	"model"
	"thegamesdb"
)

// Options given to the CLI
type Flags struct {
	Input     string // in which directory to look for games
	Output    string // in which directory outputing the resulting files (images) and gamelist.xml
	Extension string // extension to look for, separated by a space
	Platform  string // To display the list of platforms
}

// ParseFlags parses the CLI options.
func ParseFlags() Flags {
	flags := Flags{}

	flag.StringVar(&(flags.Input), "in", "", "Input directories (directory containing games)")
	flag.StringVar(&(flags.Output), "out", "", "Output directories for images and cover")
	flag.StringVar(&(flags.Extension), "ext", "zip,rar", "Accepted extensions")
	flag.StringVar(&(flags.Platform), "p", "", "'display' to prints the platform available")

	flag.Parse()

	// TODO be sure that it ends with /

	return flags
}

func lookForFiles(directory string) []string {
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
		if !fileinfo.IsDir() && filepath.Ext(name) == ".zip" { // TODO
			results = append(results, name)
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

	if flags.Platform == "display" {
		printPlatforms()
		os.Exit(0)
	}

	filenames := lookForFiles(flags.Input)
	gamesinfo := model.NewGamesinfo()

	// Create/lock the gamelist.xml for EmulationStation 2.0
	// Do it now because it's useless to scrap everything if
	// we can't write the result here...
	file, err := os.Create(flags.Output + "gamelist.xml")

	if err != nil {
		log.Printf("[err] Can't create %sgamelist.xml: %s\n", err.Error())
		os.Exit(1)
	}

	client := thegamesdb.NewClient()
	for _, filename := range filenames {
		gameinfo, err := client.Find(filename, flags.Platform, flags.Output)
		if err != nil {
			log.Println("[err] Unable to find info for the game:", filename)
			continue
		}
		gamesinfo.AddGame(gameinfo)
	}

	data, err := common.Encode(gamesinfo)
	if err != nil {
		log.Println("[err] Failed to encode the final xml:", err.Error())
		os.Exit(1)
	}

	file.Write(data)
	file.Close()

	log.Println("gamelist.xml written.")
}

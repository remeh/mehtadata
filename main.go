package main

import (
	"flag"
	"io/ioutil"
	"path/filepath"

	"thegamesdb"
)

// Options given to the CLI
type Flags struct {
	Input     string // in which directory to look for games
	Output    string // in which directory outputing the resulting files (images) and gamelist.xml
	Extension string // extension to look for, separated by a space
}

// ParseFlags parses the CLI options.
func ParseFlags() Flags {
	flags := Flags{}

	flag.StringVar(&(flags.Input), "input", "", "Input directories (directory containing games)")
	flag.StringVar(&(flags.Output), "output", "", "Output directories for images and cover")

	flag.Parse()
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

func main() {
	flags := ParseFlags()

	filenames := lookForFiles(flags.Input)

	client := thegamesdb.NewClient()
	for _, filename := range filenames {
		client.Find(filename, "nope")
	}
}

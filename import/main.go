package main

import (
	"flag"
	"log"
	"os"

	"common"
	"db"
)

type Flags struct {
	In       string
	Out      string
	Platform int
}

func main() {
	flags := parseFlags()

	file, err := os.Open(flags.In)
	if err != nil {
		log.Println("[err]", err.Error())
		os.Exit(1)
	}
	file.Close()

	gamesinfo, err := common.Decode(flags.In)
	if err != nil {
		log.Println("[err]", err.Error())
		os.Exit(1)
	}

	db.WriteDatabase(flags.Out, flags.Platform, &gamesinfo)
}

func parseFlags() Flags {
	var filename string
	var out string
	var p int
	flag.StringVar(&filename, "in", "gamelist.xml", "EmulationStation games list.")
	flag.StringVar(&out, "out", "database.db", "Output mehstation database")
	flag.IntVar(&p, "p", 0, "Output platform in mehstation database")
	flag.Parse()
	return Flags{
		In:       filename,
		Out:      out,
		Platform: p,
	}
}

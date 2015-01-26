package model

import "encoding/xml"

func NewGamesinfo() *Gamesinfo {
	return &Gamesinfo{
		Games: make([]Gameinfo, 0),
	}
}

type Gamesinfo struct {
	XMLName xml.Name   `xml:"gameList"`
	Games   []Gameinfo `xml:"game"`
}

func (g *Gamesinfo) AddGame(game Gameinfo) {
	g.Games = append(g.Games, game)
}

// General information on a game.
type Gameinfo struct {
	Filepath        string   `xml:"path"`
	Title           string   `xml:"name"`
	Platform        string   `xml:"-"`
	Description     string   `xml:"desc"`
	Genres          string   `xml:"genre"`
	Players         string   `xml:"players"`
	Publisher       string   `xml:"publisher"`
	Developer       string   `xml:"developer"`
	ReleaseDate     string   `xml:"released"`
	CoverPath       string   `xml:"image"`
	ScreenshotPaths []string `xml:"-"`
	FanartPaths     []string `xml:"-"`
	Rating          float32  `xml:"rating"`
}

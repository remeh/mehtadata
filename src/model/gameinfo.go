package model

import "time"

// General information on a game.
type Gameinfo struct {
	Title       string
	Platform    Platform
	Publisher   string
	Developer   string
	ReleaseDate time.Time
	Cover       string
	Screenshots []string
	Fanarts     []string
	Thumbnails  []string
	Description string
	Rating      float32
	Genres      []string
	Players     string
}

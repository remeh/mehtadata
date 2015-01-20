package model

// General information on a game.
type Gameinfo struct {
	Title           string
	Platform        Platform
	Description     string
	Genres          string
	Players         string
	Publisher       string
	Developer       string
	ReleaseDate     string
	CoverPath       string
	ScreenshotPaths []string
	FanartPaths     []string
	ThumbnailPaths  []string
	Rating          float32
}

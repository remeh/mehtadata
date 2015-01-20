// TheGamesDB - GetGames response
//
// Rémy 'remeh' MATHIEU © 2014

package thegamesdb

import (
	"encoding/xml"
	"strconv"
	"strings"

	"model"
)

type GetGame struct {
	XMLName      xml.Name    `xml:"Data"`
	BaseImageURL string      `xml:"baseImgUrl"`
	Game         GetGameGame `xml:"Game"`
}

type GetGameGame struct {
	Id          int `xml:"id"`
	GameTitle   string
	ReleaseDate string
	Overview    string
	Platform    model.Platform
	Genres      GetGameGenres
	Youtube     string
	ESRB        string
	Publisher   string
	Developer   string
	Rating      string
	Images      GetGameImages
}

type GetGameGenres struct {
	Genres []string `xml:"genre"`
}

type GetGameImages struct {
	Fanarts     []GetGameFanart     `xml:"fanart"`
	Boxarts     []GetGameBoxart     `xml:"boxart"`
	Screenshots []GetGameScreenshot `xml:"screenshot"`
}

type GetGameFanart struct {
	Original string `xml:"original"`
	Thumb    string `xml:"thumb"`
}

type GetGameBoxart struct {
	Side   string `xml:"side,attr"`
	Thumb  string `xml:"thumb,attr"`
	Boxart string `xml:",innerxml"` // NOTE ",innerxml" ?
}

type GetGameScreenshot struct {
	XMLName  xml.Name `xml:"screenshot"`
	Thumb    string   `xml:"thumb"`
	Original string   `xml:"original"`
}

func (gg GetGame) ToGameinfo() model.Gameinfo {
	g := gg.Game

	genre := ""
	genres := g.Genres.Genres
	if len(genres) > 0 {
		genre = strings.Join(genres, ", ")
	}

	rating := 0.0
	rating, _ = strconv.ParseFloat(g.Rating, 32)

	return model.Gameinfo{
		Title:       g.GameTitle,
		Platform:    g.Platform,
		Publisher:   g.Publisher,
		Developer:   g.Developer,
		ReleaseDate: g.ReleaseDate,
		// TODO paths for screenshots etc.
		Description: g.Overview,
		Genres:      genre,
		Rating:      float32(rating),
	}
}

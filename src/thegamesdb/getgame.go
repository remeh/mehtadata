// TheGamesDB - GetGames response
//
// Rémy 'remeh' MATHIEU © 2014

package thegamesdb

import (
	"encoding/xml"

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
	Genre []string `xml:"genre"`
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

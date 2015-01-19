// TheGamesDB - GetGamesList response.
//
// Rémy 'remeh' MATHIEU © 2014

package thegamesdb

import (
	"encoding/xml"
)

type GetGamesList struct {
	XMLName xml.Name           `xml:"Data"`
	Games   []GetGamesListGame `xml:"Game"`
}

type GetGamesListGame struct {
	Id          int `xml:"id"`
	GameTitle   string
	ReleaseDate string
	Platform    string
}

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

func ConcatGetGamesList(first *GetGamesList, second *GetGamesList) {
	if first == nil || second == nil {
		return
	}
	for _, g := range second.Games {
		first.Games = append(first.Games, g)
	}
}

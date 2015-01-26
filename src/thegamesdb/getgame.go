// TheGamesDB - GetGames response
//
// Rémy 'remeh' MATHIEU © 2014

package thegamesdb

import (
	"encoding/xml"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"common"
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
	Platform    string
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
	Boxart string `xml:",innerxml"`
}

type GetGameScreenshot struct {
	XMLName  xml.Name `xml:"screenshot"`
	Thumb    string   `xml:"thumb"`
	Original string   `xml:"original"`
}

// ToGameinfo converts the GetGame to a Gameinfo.
// During the conversion, it downloads the whole set of available images.
// FIXME A bit of refactoring could be done here... 2015-01-21 - remy
func (gg GetGame) ToGameinfo(outputDirectory string, gameFilename string) model.Gameinfo {
	g := gg.Game

	// misc fields

	genre := ""
	genres := g.Genres.Genres
	if len(genres) > 0 {
		genre = strings.Join(genres, ", ")
	}

	rating := 0.0
	rating, _ = strconv.ParseFloat(g.Rating, 32)

	var wg sync.WaitGroup

	// fanarts

	fanarts := make([]string, 0)
	for i, v := range gg.Game.Images.Fanarts {
		wg.Add(1)

		go func(wg *sync.WaitGroup, path string, i int) {
			defer wg.Done()

			ext := filepath.Ext(path)
			filename, err := common.Download(gg.BaseImageURL+path, gameFilename, "-fanart-"+strconv.Itoa(i)+ext, outputDirectory)
			if err != nil {
				log.Println("[err] While downloading ", gg.BaseImageURL+path, ":", err.Error())
			} else {
				fanarts = append(fanarts, outputDirectory+filename)
			}
		}(&wg, v.Original, i)
	}

	// screenshots

	screenshots := make([]string, 0)
	for i, v := range gg.Game.Images.Screenshots {
		wg.Add(1)

		go func(wg *sync.WaitGroup, path string, i int) {
			defer wg.Done()

			ext := filepath.Ext(path)
			filename, err := common.Download(gg.BaseImageURL+path, gameFilename, "-screenshot-"+strconv.Itoa(i)+ext, outputDirectory)
			if err != nil {
				log.Println("[err] While downloading ", gg.BaseImageURL+path, ":", err.Error())
			} else {
				screenshots = append(screenshots, outputDirectory+filename)
			}
		}(&wg, v.Original, i)
	}

	// look for a front cover
	front := gg.havingFront(gg.Game.Images.Boxarts)
	var coverURL string
	var cover string
	if front > -1 {
		coverURL = gg.Game.Images.Boxarts[front].Boxart
	} else if len(gg.Game.Images.Boxarts) > 0 {
		// No front, take something
		coverURL = gg.Game.Images.Boxarts[0].Boxart
	}
	// something to download for the cover
	if coverURL != "" {
		wg.Add(1)

		go func(wg *sync.WaitGroup) {
			defer wg.Done()

			ext := filepath.Ext(coverURL)
			filename, err := common.Download(gg.BaseImageURL+coverURL, gameFilename, "-cover"+ext, outputDirectory)
			if err != nil {
				log.Println("[err] While downloading ", gg.BaseImageURL+coverURL, ":", err.Error())
			} else {
				cover = outputDirectory + filename
			}
		}(&wg)
	}

	wg.Wait()

	return model.Gameinfo{
		Filepath:        outputDirectory + gameFilename,
		Title:           g.GameTitle,
		Platform:        g.Platform,
		Publisher:       g.Publisher,
		Developer:       g.Developer,
		ReleaseDate:     g.ReleaseDate,
		ScreenshotPaths: screenshots,
		FanartPaths:     fanarts,
		CoverPath:       cover,
		Description:     g.Overview,
		Genres:          genre,
		Rating:          float32(rating),
	}
}

// havingFront returns true if there is a front cover.
func (gg GetGame) havingFront(boxarts []GetGameBoxart) int {
	for i, v := range boxarts {
		if v.Side == "front" {
			return i
		}
	}
	return -1
}

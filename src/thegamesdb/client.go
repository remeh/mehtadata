// TheGamesDB - Client
//
// Rémy 'remeh' MATHIEU © 2014

package thegamesdb

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"sync"

	"common"
	. "model"
)

const (
	THEGAMESDB_API_URL = "http://thegamesdb.net/api"

	THEGAMESDB_GETGAMESLIST = "/GetGamesList.php"
	THEGAMESDB_GETGAME      = "/GetGame.php"

	MAX_RETRIEVED_GAMES                = 5
	MINIMUM_RATING_AUTOMATIC_SELECTION = 0.8
)

// Client for TheGamesDB.net API.
type Client struct {
}

func NewClient() *Client {
	return &Client{}
}

// A match of a query
type Match struct {
	Game   GetGamesListGame
	Rating float32
}

// To sort by bests results.

type Matches []Match

func (m Matches) Len() int {
	return len(m)
}

func (m Matches) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func (m Matches) Less(i, j int) bool {
	return m[j].Rating < m[i].Rating
}

// Does the HTTP call to find game information on TheGamesDB
func (c *Client) Find(name string, platform string, outputDirectory string) (Gameinfo, error) {
	url := THEGAMESDB_API_URL + THEGAMESDB_GETGAMESLIST + "?name=" + url.QueryEscape(common.ClearName(name)) + "&platform=" + url.QueryEscape(platform)

	// HTTP call
	resp, err := http.Get(url)
	if err != nil {
		return Gameinfo{}, err
	}

	defer resp.Body.Close()

	// Read the response
	readBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Gameinfo{}, err
	}

	// Unmarshal the XML
	var gamesList GetGamesList
	xml.Unmarshal(readBody, &gamesList)

	list := c.findBestMatches(name, platform, gamesList, MAX_RETRIEVED_GAMES)

	if len(list) == 0 {
		return Gameinfo{}, nil // we can't find anything on TheGamesDB
	}

	// Sort by rating
	sort.Sort(list)

	// If the rating isn't enough, crawl many games
	// to propose a choice to the user.
	if list[0].Rating < MINIMUM_RATING_AUTOMATIC_SELECTION {
		// TODO with the method findSome
	}

	// The first one has a sufficient rating to be automatically used
	gotGame, err := c.FindGame(list[0].Game, platform)
	if err != nil {
		return Gameinfo{}, err
	}

	return gotGame.ToGameinfo(outputDirectory, name), nil
}

func (c *Client) findSome(list Matches, platform string) []GetGame {
	// Asynchronously get information for the games
	// We need many to propose something.
	var waitGroup sync.WaitGroup
	results := make([]GetGame, 0)
	for i, _ := range list {
		// One more to execute
		waitGroup.Add(1)

		go func(waitGroup *sync.WaitGroup, results *[]GetGame, game GetGamesListGame, platform string) {
			defer waitGroup.Done() // Signal the end of the execution of the routine.

			gotGame, err := c.FindGame(game, platform)
			if err == nil {
				*results = append(*results, gotGame)
			} else {
				fmt.Println("[tgdb] [err] While querying for '", game.GameTitle, "':", err.Error())
			}
		}(&waitGroup, &results, list[i].Game, platform)
	}

	waitGroup.Wait()
	return results
}

// FindGame does one call to TheGamesDB to retrieve one game
// information by its ID.
func (c *Client) FindGame(game GetGamesListGame, platform string) (GetGame, error) {
	url := THEGAMESDB_API_URL + THEGAMESDB_GETGAME + "?id=" + url.QueryEscape(fmt.Sprintf("%d", game.Id)) + "&platform=" + url.QueryEscape(platform)

	// HTTP call
	resp, err := http.Get(url)
	if err != nil {
		return GetGame{}, err
	}

	defer resp.Body.Close()

	// Read the response
	readBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return GetGame{}, err
	}

	// Unmarshal the XML
	var gotGame GetGame
	xml.Unmarshal(readBody, &gotGame)

	return gotGame, nil
}

// findBestMatch tries to find with the name and platform the best matching
// game available in the list of responses from the TheGamesDB search query.
// findBestMatches returned an ordered by best list of matches.
func (c *Client) findBestMatches(name string, platform string, gamesList GetGamesList, count int) Matches {
	name = common.ClearName(name)

	results := make(Matches, 0)

	// iter through the results of the search
	// ands assigns them a rating.
	for _, v := range gamesList.Games {
		rating := common.CompareFilename(v.GameTitle, name)
		results = append(results, Match{Game: v, Rating: rating})

		// stop if we have enough
		if len(results) == count {
			return results
		}
	}

	return results
}

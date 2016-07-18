package main

import (
	"fmt"
	"thegamesdb"
)

// TODO(remy): could hit TheGamesDB for a real
// list of supported platforms.
func ShowPlatforms() {
	for _, p := range thegamesdb.TGDBPlatforms {
		fmt.Printf("%s\n", p.Name)
	}
}

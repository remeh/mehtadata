package main

import (
	"fmt"
	"os"

	"github.com/remeh/mehtadata/db"
)

func NewResource(flags Flags) (bool, error) {
	// try to read in the env var if
	// everything is available

	// basic infos
	filepath := os.Getenv("FILEPATH")
	resource := os.Getenv("RESOURCE")
	platformName := os.Getenv("PLATFORM_NAME")
	typ := os.Getenv("TYPE")

	ok := false

	if StringsHasContent(filepath, platformName, typ) {
		ok = true
	}

	if !ok {
		fmt.Println(`Can't create a new resource.
Mandatory infos:
	RESOURCE			: filepath to the resources to add
	FILEPATH			: filepath of the executable to which you want to add a resource
	PLATFORM_NAME		: name of the platform in which there is the executable  
	TYPE				: type of resource (screenshot, fanart, cover, video or logo)
		`)
		return false, fmt.Errorf("Missing fields.")
	}

	switch typ {
	case "screenshot", "fanart", "cover", "video", "logo":
	default:
		return false, fmt.Errorf("Error: unknown resource type: %v", typ)
	}
	return db.CreateResource(flags.DestSqlite, resource, filepath, platformName, typ)
}

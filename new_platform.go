package main

import (
	"fmt"
	"os"
)

func NewPlatform() (int64, error) {
	// try to read in the env var if
	// everything is available

	// basic infos
	name := os.Getenv("NAME")
	command := os.Getenv("COMMAND")
	typ := os.Getenv("TYPE")
	//icon := os.Getenv("ICON")
	//background := os.Getenv("BG")

	// discover mode
	discover_dir := os.Getenv("DIR")
	discover_exts := os.Getenv("EXTS")

	if typ != "cover" && typ != "complete" {
		typ = "complete"
	}

	ok := false
	discover := false

	if StringsHasContent(name, command) {
		ok = true
	}

	if StringsHasContent(discover_dir, discover_exts) {
		discover = true
	}

	if !ok {
		fmt.Println(`Can't create a new platform.
		Mandatory infos:
		NAME      : name of the platform
		COMMAND   : absolute path to the command with the %exec% flag to start the platform on an executable.
		            Ex:  COMMAND="/usr/bin/retroarch -L /usr/lib/libretro/scumm.so %exec%"
		Discover mode:
		DIR       : which directory contains the executables which must be discovered.
		EXTS      : extensions of the executables when scanning the directory.
		Not mandatory:
		TYPE      : display format. Possible values: "complete", "cover"
		ICON      : absolute path to the icon image to use for this platform.
		BG        : absolute path to the background image to use for this platform.
		`)
		return -1, fmt.Errorf("Missing fields.")
	}

	// TODO(remy): create the platform.
	if discover {
	}
	return 0, nil
}

// StringsHasContent tests that every given string
// has a size > 0.
func StringsHasContent(strings ...string) bool {
	for _, s := range strings {
		if len(s) == 0 {
			return false
		}
	}
	return true
}

package thegamesdb

type Platform struct {
	Id    int
	Name  string
	Alias string
}

var Platforms = []Platform{
	Platform{
		25,
		"3DO",
	},
	Platform{
		4911,
		"Amiga",
	},
	Platform{
		4914,
		"Amstrad CPC",
	},
	Platform{
		4916,
		"Android",
	},
	Platform{
		23,
		"Arcade",
	},
	Platform{
		22,
		"Atari 2600",
	},
}

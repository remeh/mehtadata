package main

import (
	"fmt"
	"os"
	"strings"

	"thegamesdb"
)

func main() {
	client := thegamesdb.NewClient()
	fmt.Println(client.Find(strings.Join(os.Args[1:], " "), "Nope"))
}

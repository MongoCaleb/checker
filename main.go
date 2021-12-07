package main

import (
	"checker/parsers"
	"fmt"
)

func main() {
	refMap := parsers.Intersphinx("https://docs.mongodb.com/drivers/go/current/objects.inv")
	fmt.Println(refMap)
}

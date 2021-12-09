package main

import (
	"checker/sources"
	"fmt"
)

func main() {
	roleMap := sources.NewRoleMap(sources.GetLatestSnootyParserTag())
	fmt.Println(roleMap)
}

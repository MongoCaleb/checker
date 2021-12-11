package main

import (
	"checker/sources"
	"fmt"
)

func main() {
	roleMap := sources.NewRoleMap(sources.GetLatestSnootyParserTag())
	for _, role := range roleMap {
		fmt.Printf(role, "foo\n")
	}
}

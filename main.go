package main

import (
	"fmt"
)

func main() {
	roleMap := NewRoleMap(GetLatestSnootyParserTag())
	fmt.Println(roleMap)
}

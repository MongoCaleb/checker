package rstparser

import (
	"fmt"
	"regexp"
	"strings"
)

func Parse() {
	fmt.Println("parser")
}

type Location struct {
	Line   int
	Column int
}

type Ref struct {
	Match    string
	Location Location
}

func refParse(input []byte) []string {
	matches := make([]string, 0)
	re := regexp.MustCompile(`:([\w\s\-\.\d_=+!@#$%^&*(\)]*):\x60(?P<one>[\w\s\-_\.\d\\\/=+!@#$%^&*(\)]*)<?(?P<two>[\w\s\-_\.\d\\\/=+!@#$%^&*(\)]*)>?`)
	allIndexes := re.FindAllSubmatchIndex([]byte(input), -1)
	for _, loc := range allIndexes {
		extract := input[loc[0]:loc[1]]
		innerMatches := re.FindAllStringSubmatch(strings.Join(strings.Fields(string(extract)), ""), -1)
		for _, match := range innerMatches {
			if match[3] == "" {
				matches = append(matches, match[2])
			} else {
				matches = append(matches, match[3])
			}
		}
	}
	return matches
}

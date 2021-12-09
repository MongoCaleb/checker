package main

import (
	"regexp"
	"strings"
)

type RstConstant struct {
	Name   string
	Target string
}

func ParseForConstants(input string) []RstConstant {
	constants := make([]RstConstant, 0)
	re := regexp.MustCompile(`<\{\+([\w\s\-_\.\d\\\/=+!@#$%^&*(\)]*)\+\}(\/[\w\s\-_\.\d\\\/=+!@#$%^&*(\)]*)>\x60`)
	allIndexes := re.FindAllSubmatchIndex([]byte(input), -1)
	for _, loc := range allIndexes {
		extract := input[loc[0]:loc[1]]
		innerMatches := re.FindAllStringSubmatch(strings.Join(strings.Fields(string(extract)), ""), -1)
		for _, match := range innerMatches {
			constants = append(constants, RstConstant{Target: match[2], Name: match[1]})
		}
	}
	return constants
}

package parsers

import (
	"checker/types"
	"regexp"
	"strings"
)

func ParseForConstants(input []byte) []types.RstConstant {
	matches := make([]types.RstConstant, 0)
	re := regexp.MustCompile(`<\{\+(.*)\+\}(\/.*)>\x60`)
	allIndexes := re.FindAllSubmatchIndex([]byte(input), -1)
	for _, loc := range allIndexes {
		extract := input[loc[0]:loc[1]]
		innerMatches := re.FindAllStringSubmatch(strings.Join(strings.Fields(string(extract)), ""), -1)
		for _, match := range innerMatches {
			matches = append(matches, types.RstConstant{Target: match[2], Name: match[1]})
		}
	}
	return matches
}

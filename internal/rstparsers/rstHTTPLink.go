package rstparsers

import (
	"regexp"
	"strings"
)

type RstHTTPLink string

func ParseForHTTPLinks(input []byte) []RstHTTPLink {
	re := regexp.MustCompile(`(https?:\/\/[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9]{1,6}\b[-a-zA-Z0-9@:%_\+.~#?&//=]*)`)

	links := make([]RstHTTPLink, 0)
	allIndexes := re.FindAllSubmatchIndex(input, -1)
	for _, loc := range allIndexes {
		extract := input[loc[0]:loc[1]]
		innerMatches := re.FindAllStringSubmatch(strings.Join(strings.Fields(string(extract)), ""), -1)
		for _, match := range innerMatches {
			links = append(links, RstHTTPLink(match[0]))
		}
	}
	return links
}

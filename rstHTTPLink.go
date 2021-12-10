package main

import (
	"regexp"
	"strings"
)

type RstHTTPLink string

func ParseForLinks(text string) []RstHTTPLink {
	re := regexp.MustCompile(`(https?:\/\/[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9]{1,6}\b[-a-zA-Z0-9@:%_\+.~#?&//=]*)`)

	links := make([]RstHTTPLink, 0)
	allIndexes := re.FindAllSubmatchIndex([]byte(text), -1)
	for _, loc := range allIndexes {
		extract := text[loc[0]:loc[1]]
		innerMatches := re.FindAllStringSubmatch(strings.Join(strings.Fields(string(extract)), ""), -1)
		for _, match := range innerMatches {
			links = append(links, RstHTTPLink(match[0]))
		}
	}
	return links
}

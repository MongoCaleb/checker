package main

import (
	"regexp"
	"strings"
)

type LinkMap struct {
	Files []RstFile
}

func NewLinkMap() *LinkMap {
	return &LinkMap{}
}

func (lm *LinkMap) ParseForLinks(fileName, text string) *LinkMap {
	re := regexp.MustCompile(`(https?:\/\/[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9]{1,6}\b[-a-zA-Z0-9@:%_\+.~#?&//=]*)`)

	var rstFile RstFile
	rstFile.Name = fileName
	rstFile.Links = make([]HTTPLink, 0)
	allIndexes := re.FindAllSubmatchIndex([]byte(text), -1)
	for _, loc := range allIndexes {
		extract := text[loc[0]:loc[1]]
		innerMatches := re.FindAllStringSubmatch(strings.Join(strings.Fields(string(extract)), ""), -1)
		for _, match := range innerMatches {
			rstFile.Links = append(rstFile.Links, HTTPLink(match[0]))
		}
	}
	lm.Files = append(lm.Files, rstFile)
	return lm
}

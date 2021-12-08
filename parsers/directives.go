package parsers

import (
	"checker/types"
	"regexp"
	"strings"
)

func ParseForRoles(input []byte) []types.RstRole {
	matches := make([]types.RstRole, 0)
	re := regexp.MustCompile(`:([\w\s\-\.\d_=+!@#$%^&*(\)]*):\x60([\w\s\-_\.\d\\\/=+!@#$%^&*(\)]*)<?([\w\s\-_\.\d\\\/=+!@#$%^&*(\)]*)>?`)
	allIndexes := re.FindAllSubmatchIndex([]byte(input), -1)
	for _, loc := range allIndexes {
		extract := input[loc[0]:loc[1]]
		innerMatches := re.FindAllStringSubmatch(strings.Join(strings.Fields(string(extract)), ""), -1)
		for _, match := range innerMatches {
			roleType, name := "", ""
			if match[1] == "ref" {
				roleType = "ref"
				name = "ref"
			} else {
				roleType = "role"
				name = match[1]
			}
			if match[3] == "" {
				matches = append(matches, types.RstRole{Target: match[2], RoleType: roleType, Name: name})
			} else {
				matches = append(matches, types.RstRole{Target: match[3], RoleType: roleType, Name: name})
			}
		}
	}
	return matches
}

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

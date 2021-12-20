package rst

import (
	"regexp"
	"strings"
)

var (
	constantRegex = regexp.MustCompile(`<\{\+([\w\s\-_\.\d\\\/=+!@#$%^&*(\)]*)\+\}(\/[\w\s\-_\.\d\\\/=+!@#$%^&*(\)]*)>\x60`)
	httpLinkRegex = regexp.MustCompile(`(https?:\/\/[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9]{1,6}\b[-a-zA-Z0-9@:%_\+.~#?&//=]*)`)
	roleRegex     = regexp.MustCompile(`:([\w\s\-_\.\d\\\/=+!@#$%^&*(\)]*):\x60([\w\s\-_\.\d\\\/=+!@#$%^&*(\)]*)<?([\w\s\-_\.\d\\\/=+!@#$%^&*(\)]*)>?`)
	localRefRegex = regexp.MustCompile(`\.\.\s_([\w\-_=+!@#$%^&(\)]+):`)
)

type RstHTTPLink string

type RstRole struct {
	Target   string
	RoleType string
	Name     string
}

type RstConstant struct {
	Name   string
	Target string
}
type RefTarget struct {
	Target string
	Type   string
}

func parse(input []byte, re regexp.Regexp, fn func(matches []string)) {
	allFound := re.FindAllString(string(input), -1)
	for _, match := range allFound {
		for _, match := range re.FindAllStringSubmatch(strings.Join(strings.Fields(match), ""), -1) {
			fn(match)
		}
	}
}

func ParseForHTTPLinks(input []byte) []RstHTTPLink {
	links := make([]RstHTTPLink, 0)
	parse(input, *httpLinkRegex, func(matches []string) {
		links = append(links, RstHTTPLink(matches[0]))
	})
	return links
}

func ParseForRoles(input []byte) []RstRole {
	roles := make([]RstRole, 0)
	parse(input, *roleRegex, func(matches []string) {
		roleType, name := "", ""
		if matches[1] == "ref" {
			roleType = "ref"
			name = "ref"
		} else {
			roleType = "role"
			name = matches[1]
		}
		if matches[3] == "" {
			roles = append(roles, RstRole{Target: matches[2], RoleType: roleType, Name: name})
		} else {
			roles = append(roles, RstRole{Target: matches[3], RoleType: roleType, Name: name})
		}
	})
	return roles
}

func ParseForConstants(input []byte) []RstConstant {
	constants := make([]RstConstant, 0)
	parse(input, *constantRegex, func(matches []string) {
		constants = append(constants, RstConstant{Target: matches[2], Name: matches[1]})
	})
	return constants
}

func (r *RstConstant) IsHTTPLink() bool {
	return httpLinkRegex.Match([]byte(r.Target))
}

func ParseForLocalRefs(input []byte) []RefTarget {
	localrefs := make([]RefTarget, 0)

	allIndexes := localRefRegex.FindAllString(string(input), -1)
	for _, match := range allIndexes {
		innerMatches := localRefRegex.FindAllStringSubmatch(match, -1)
		for _, match := range innerMatches {
			localrefs = append(localrefs, RefTarget{Target: match[1], Type: "local"})
		}
	}
	return localrefs
}

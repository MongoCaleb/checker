package rst

import (
	"regexp"
	"strings"
)

var (
	constantRegex      = regexp.MustCompile(`<\{\+([\w\s\-_\.\d\\\/=+!@#$%^&*(\)]*)\+\}(\/[\w\s\-_\.\d\\\/=+!@#$%^&*(\)]*)>\x60`)
	httpLinkRegex      = regexp.MustCompile(`(https?:\/\/[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9]{1,6}\b[-a-zA-Z0-9@:%_\+.~#?&//=]*)`)
	roleRegex          = regexp.MustCompile(`:([[:alnum:]\.]+):\x60([^\x60]+)`)
	localRefRegex      = regexp.MustCompile(`\.\. +_([\-_=+!@#$%^&\(\)\w\d\p{P}\p{S} ]+):`)
	sharedIncludeRegex = regexp.MustCompile(`\.\. sharedinclude::\s([\w\-_\.\d\\\/=+!@#$%^&*(\)\[\]\\\<\>'\?]+)`)
	directiveRegex     = regexp.MustCompile(`\.\.\s([[:alnum:]]+)::\s([[:graph:] ]+)`)
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
	Name string
}

type SharedInclude struct {
	Path string
}

type RstDirective struct {
	Name   string
	Target string
}

func parse(input []byte, re regexp.Regexp, fn func(matches []string)) {
	allFound := re.FindAllString(string(input), -1)
	for _, match := range allFound {
		for _, match := range re.FindAllStringSubmatch(match, -1) {
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
	allFound := roleRegex.FindAllString(string(input), -1)
	for _, match := range allFound {
		for _, m := range roleRegex.FindAllStringSubmatch(match, -1) {
			matches := make([]string, 2)
			if strings.TrimSpace(m[1]) != "" {
				matches[0] = m[1]
			}
			if strings.HasSuffix(m[2], ">") {
				lastClosingBracket := strings.LastIndex(m[2], ">")
				lastOpeningBracket := strings.LastIndex(m[2], "<")
				matches[1] = m[2][lastOpeningBracket+1 : lastClosingBracket]
			} else {
				matches[1] = m[2]
			}
			roleType, name := "", ""
			if matches[0] == "ref" {
				roleType = "ref"
				name = "ref"
			} else {
				roleType = "role"
				name = matches[0]
			}
			roles = append(roles, RstRole{Target: matches[1], RoleType: roleType, Name: name})
		}
	}
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
	parse(input, *localRefRegex, func(matches []string) {
		localrefs = append(localrefs, RefTarget{Name: matches[1]})
	})

	return localrefs
}

func ParseForSharedIncludes(input []byte) []SharedInclude {
	shared := make([]SharedInclude, 0)
	parse(input, *sharedIncludeRegex, func(matches []string) {
		shared = append(shared, SharedInclude{Path: matches[1]})
	})
	return shared
}

func ParseForDirectives(input []byte) []RstDirective {
	directives := make([]RstDirective, 0)
	parse(input, *directiveRegex, func(matches []string) {
		directives = append(directives, RstDirective{Name: matches[1], Target: matches[2]})
	})
	return directives
}

package main

import (
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

type RstRole struct {
	Target   string
	RoleType string
	Name     string
}

func ParseForRoles(input string) []RstRole {
	roles := make([]RstRole, 0)
	re := regexp.MustCompile(`:([\w\s\-_\.\d\\\/=+!@#$%^&*(\)]*):\x60([\w\s\-_\.\d\\\/=+!@#$%^&*(\)]*)<?([\w\s\-_\.\d\\\/=+!@#$%^&*(\)]*)>?`)
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
				roles = append(roles, RstRole{Target: match[2], RoleType: roleType, Name: name})
			} else {
				roles = append(roles, RstRole{Target: match[3], RoleType: roleType, Name: name})
			}
		}
	}
	return roles
}

func (r *RstRole) ToHTTPLink(refmap *RefMap, rolemap *RoleMap) HTTPLink {
	switch r.RoleType {
	case "ref":
		if ref, ok := refmap.Get(r.Name); !ok {
			log.Errorf("Could not find ref %s", r.Name)
			return HTTPLink("an.invalid.ref." + r.Name)
		} else {
			return HTTPLink(ref)
		}
	case "role":
		if role, ok := rolemap.Get(r.Name); !ok {
			log.Errorf("Could not find role %s", r.Name)
			return HTTPLink("an.invalid.role." + r.Name)
		} else {
			return HTTPLink(role)
		}
	default:
		return HTTPLink("an.invalid.something." + r.Name)
	}
}

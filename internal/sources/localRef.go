package sources

import (
	"regexp"
)

func FindLocalRefs(input string) map[string]Ref {
	refs := make(map[string]Ref)

	re := regexp.MustCompile(`\.\.\s_([\w\-_=+!@#$%^&(\)]+):`)
	allIndexes := re.FindAllSubmatchIndex([]byte(input), -1)
	for _, loc := range allIndexes {
		extract := input[loc[0]:loc[1]]
		innerMatches := re.FindAllStringSubmatch(extract, -1)
		for _, match := range innerMatches {
			refs[match[1]] = Ref{Target: match[1], Type: "local"}
		}
	}
	return refs
}

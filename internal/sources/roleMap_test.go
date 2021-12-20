package sources

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	roleMapInput = `
[foo]
rfc = "https://tools.ietf.org/html/%s"

[foo.nope]
help = """Shouldn't show up"""
type = {link = "https://nope.come/%s"}

[role.abbr]
help = """Abbreviation with hover text."""
type = "text"

[role.file]
help = """Show a file path."""
type = "text"

[role.icon-fa4]
help = """Show a FontAwesome 4 icon."""
type = "explicit_title"

[role.rfc]
help = """Reference an IETF RFC."""
type = {link = "https://tools.ietf.org/html/%s"}

[role.wikipedia]
help = """Reference a Wikipedia page."""
type = {link = "https://en.wikipedia.org/wiki/%s"}

`
)

func TestRoleMap(t *testing.T) {

	roleMap := NewRoleMap([]byte(roleMapInput))

	if len(roleMap) != 2 {
		t.Errorf("Expected 2 roles, got %d", len(roleMap))
	}

	expected := map[string]string{"rfc": "https://tools.ietf.org/html/%s", "wikipedia": "https://en.wikipedia.org/wiki/%s"}

	assert.EqualValues(t, expected, roleMap, "Expected %v, got %v", expected, roleMap)

}

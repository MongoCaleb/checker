package sources

import (
	"testing"
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
	for k, v := range expected {
		if roleMap[k] != v {
			t.Errorf("Expected %s to be %s, got %s", k, v, roleMap[k])
		}
	}
}

func TestGet(t *testing.T) {

	roleMap := NewRoleMap([]byte(roleMapInput))

	tests := []struct {
		expected bool
		input    string
	}{
		{true, "rfc"},
		{true, "wikipedia"},
		{false, "nope"},
	}

	for _, test := range tests {
		_, got := roleMap.Get(test.input)
		if got != test.expected {
			t.Errorf("Expected %s to be %t, got %t", test.input, test.expected, got)
		}
	}

}

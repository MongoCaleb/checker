package sources

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	rstSpec = `
[directive.default-domain]
argument_type = "string"

[directive.div]
deprecated = true
argument_type = "string"
content_type = "block"

[directive.container]
deprecated = true
argument_type = "string"
content_type = "block"

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

[rstobject."py:class"]

[rstobject."py:meth"]
type = "callable"

[rstobject."js:func"]

[rstobject."mongodb:projection"]
prefix = "proj"

[rstobject."mongodb:method"]
type = "callable"
fields = [["returns", "Returns"]]

[rstobject."mongodb:authrole"]
[rstobject."mongodb:authaction"]

`
)

func TestRoleMap(t *testing.T) {

	roleMap := NewRoleMap([]byte(rstSpec))

	expected := &RstSpec{
		Roles:      map[string]string{"rfc": "https://tools.ietf.org/html/%s", "wikipedia": "https://en.wikipedia.org/wiki/%s"},
		RawRoles:   map[string]bool{"abbr": true, "file": true, "icon-fa4": true, "rfc": true, "wikipedia": true},
		Directives: map[string]bool{"div": true, "container": true, "default-domain": true},
		RstObjects: map[string]bool{"class": true, "meth": true, "func": true, "projection": true, "method": true, "authrole": true, "authaction": true},
	}

	assert.EqualValues(t, expected, roleMap)
}

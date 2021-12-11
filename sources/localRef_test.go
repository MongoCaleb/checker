package sources

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type refTest struct {
	input    string
	expected map[string]Ref
}

func TestFindLocalRefs(t *testing.T) {
	cases := []refTest{{
		input:    "",
		expected: map[string]Ref{},
	}, {
		input:    ".. _:",
		expected: map[string]Ref{},
	}, {
		input: ".. _foo:",
		expected: map[string]Ref{
			"foo": {Target: "foo", Type: "local"},
		},
	}, {
		input: ".. _foo:\n.. _bar:",
		expected: map[string]Ref{
			"foo": {Target: "foo", Type: "local"},
			"bar": {Target: "bar", Type: "local"},
		},
	}, {
		input: ".. _foo:\n.. _bar:\n\n\n\n\n\n.. _baz:",
		expected: map[string]Ref{
			"foo": {Target: "foo", Type: "local"},
			"bar": {Target: "bar", Type: "local"},
			"baz": {Target: "baz", Type: "local"},
		},
	}}

	for _, c := range cases {
		actual := FindLocalRefs(c.input)
		assert.Equal(t, len(c.expected), len(actual), "FindLocalRefs(%q) should return %d refs, got %d", c.input, len(c.expected), len(actual))

		for i, ref := range actual {
			assert.Equal(t, c.expected[i].Target, ref.Target, "FindLocalRefs(%q) should return ref %d as %q, got %q", c.input, i, c.expected[ref.Target], ref)
			assert.Equal(t, c.expected[i].Type, ref.Type, "FindLocalRefs(%q) should return ref %d as %q, got %q", c.input, i, c.expected[ref.Target], ref)
		}
	}

}

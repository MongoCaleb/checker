package sources

import "testing"

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
			"foo": Ref{Target: "foo", Type: "local"},
		},
	}, {
		input: ".. _foo:\n.. _bar:",
		expected: map[string]Ref{
			"foo": Ref{Target: "foo", Type: "local"},
			"bar": Ref{Target: "bar", Type: "local"},
		},
	}, {
		input: ".. _foo:\n.. _bar:\n\n\n\n\n\n.. _baz:",
		expected: map[string]Ref{
			"foo": Ref{Target: "foo", Type: "local"},
			"bar": Ref{Target: "bar", Type: "local"},
			"baz": Ref{Target: "baz", Type: "local"},
		},
	}}

	for _, c := range cases {
		actual := FindLocalRefs(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("FindLocalRefs(%q) == %q, expected %q", c.input, actual, c.expected)
		}

		for i, ref := range actual {
			if ref.Target != c.expected[i].Target || ref.Type != c.expected[i].Type {
				t.Errorf("FindLocalRefs(%q) == %q, expected %q", c.input, actual, c.expected)
			}
		}
	}

}

package main

import (
	"testing"
)

type constantTestCase struct {
	input    string
	expected []RstConstant
}

func TestConstantParser(t *testing.T) {

	testCases := []constantTestCase{{
		input:    "",
		expected: []RstConstant{},
	}, {
		input:    ".. _:",
		expected: []RstConstant{},
	}, {
		input:    ".. _: foo",
		expected: []RstConstant{},
	}, {
		input:    "This is a `constant link that should fail <{+api+}/flibbertypoo>`__",
		expected: []RstConstant{{Target: "/flibbertypoo", Name: "api"}},
	}, {
		input:    "This is a `constant link that should succeed <{+api+}/classes/AggregationCursor.html>`__",
		expected: []RstConstant{{Target: "/classes/AggregationCursor.html", Name: "api"}},
	}, {
		input:    "here is a :ref:`fantastic`",
		expected: []RstConstant{},
	}, {
		input:    "Here is one `constant link <{+api+}/One.html>`__ and a second `constant link <{+api+}/Two.html>`__",
		expected: []RstConstant{{Target: "/One.html", Name: "api"}, {Target: "/Two.html", Name: "api"}},
	},
	}
	for _, test := range testCases {
		got := ParseForConstants(test.input)
		for i, find := range test.expected {
			if len(got) != len(test.expected) {
				t.Errorf("expected length %d, got %d", len(test.expected), len(got))
			}
			if find != test.expected[i] {
				t.Errorf("expected %q, got %q with %q", test.expected[i], find, test.input)
			}
		}
	}
}

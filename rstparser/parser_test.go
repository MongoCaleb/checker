package rstparser

import (
	"testing"
)

type testCase struct {
	input    string
	expected []string
}

//input:    "here is a :ref:`fantastic` here is another :ref:`2 <mediocre-fantastic>` here is a :ref:`\n<not_great-fantastic>",
func TestRefParser(t *testing.T) {

	testCases := []testCase{{
		input:    "",
		expected: []string{""},
	}, {
		input:    ".. _:",
		expected: []string{""},
	}, {
		input:    ".. _: foo",
		expected: []string{""},
	}, {
		input:    "here is a :ref:`fantastic`",
		expected: []string{"fantastic"},
	}, {
		input:    "here is a :ref:`fantastic` here is another :ref:`2 <mediocre-fantastic>` here is a :ref:`\n<not_great-fantastic>`",
		expected: []string{"fantastic", "mediocre-fantastic", "not_great-fantastic"},
	}}

	for _, test := range testCases {
		got := refParse([]byte(test.input))
		for i, find := range got {
			if len(got) != len(test.expected) {
				t.Errorf("expected length %d, got %d with %q", len(test.expected), len(got), find)
			}
			if find != test.expected[i] {
				t.Errorf("refParse(%q) == %q, want %q", test.input, got, test.expected[i])
			}
		}
	}

}

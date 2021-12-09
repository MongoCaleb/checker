package parsers

import (
	"checker/types"
	"testing"
)

type constantTestCase struct {
	input    string
	expected []types.RstConstant
}

type roleTestCase struct {
	input    string
	expected []types.RstRole
}

func TestConstantParser(t *testing.T) {

	testCases := []constantTestCase{{
		input:    "",
		expected: []types.RstConstant{},
	}, {
		input:    ".. _:",
		expected: []types.RstConstant{},
	}, {
		input:    ".. _: foo",
		expected: []types.RstConstant{},
	}, {
		input:    "This is a `constant link that should fail <{+api+}/flibbertypoo>`__",
		expected: []types.RstConstant{{Target: "/flibbertypoo", Name: "api"}},
	}, {
		input:    "This is a `constant link that should succeed <{+api+}/classes/AggregationCursor.html>`__",
		expected: []types.RstConstant{{Target: "/classes/AggregationCursor.html", Name: "api"}},
	}, {
		input:    "here is a :ref:`fantastic`",
		expected: []types.RstConstant{},
	}, {
		input:    "here is a :ref:`fantastic` here is another :ref:`2 <mediocre-fantastic>` here is a :ref:`\n<not_great-fantastic>`",
		expected: []types.RstConstant{},
	}, {
		input:    ":node-api:`foo </AggregationCursor.html>`",
		expected: []types.RstConstant{},
	}, {
		input:    ":node-api:`foo <AggregationCursorz.html>`",
		expected: []types.RstConstant{},
	}, {
		input:    ":node-api:`foo <AggregationCursor.html>`",
		expected: []types.RstConstant{},
	}, {
		input:    "This is a :ref:`valid atlas ref <connect-to-your-cluster>`",
		expected: []types.RstConstant{},
	}, {
		input:    "This is a :ref:`valid server ref <replica-set-read-preference-behavior>`",
		expected: []types.RstConstant{},
	}, {
		input:    "This is an :ref:`nvalid ref <invalid_ref_sucka-fish>`",
		expected: []types.RstConstant{},
	},
	}
	for _, test := range testCases {
		got := ParseForConstants([]byte(test.input))
		for i, find := range got {
			if len(got) != len(test.expected) {
				t.Errorf("expected length %d, got %d with %q", len(test.expected), len(got), find)
			}
			if find != test.expected[i] {
				t.Errorf("expected %q, got %q with %q", test.expected[i], find, test.input)
			}
		}
	}
}

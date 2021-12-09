package parsers

import (
	"checker/types"
	"testing"
)

func TestRefParser(t *testing.T) {

	testCases := []roleTestCase{{
		input:    "",
		expected: []types.RstRole{},
	}, {
		input:    ".. _:",
		expected: []types.RstRole{},
	}, {
		input:    ".. _: foo",
		expected: []types.RstRole{},
	}, {
		input:    "This is a `constant link that should fail <{+api+}/flibbertypoo>`__",
		expected: []types.RstRole{},
	}, {
		input:    "This is a `constant link that should succeed <{+api+}/classes/AggregationCursor.html>`__",
		expected: []types.RstRole{},
	}, {
		input:    "here is a :ref:`fantastic`",
		expected: []types.RstRole{{Target: "fantastic", RoleType: "ref", Name: "ref"}},
	}, {
		input:    "here is a :ref:`fantastic` here is another :ref:`2 <mediocre-fantastic>` here is a :ref:`\n<not_great-fantastic>`",
		expected: []types.RstRole{{Target: "fantastic", RoleType: "ref", Name: "ref"}, {Target: "mediocre-fantastic", RoleType: "ref", Name: "ref"}, {Target: "not_great-fantastic", RoleType: "ref", Name: "ref"}},
	}, {
		input:    ":node-api:`foo </AggregationCursor.html>`",
		expected: []types.RstRole{{Target: "/AggregationCursor.html", RoleType: "role", Name: "node-api"}},
	}, {
		input:    ":node-api:`foo <AggregationCursorz.html>`",
		expected: []types.RstRole{{Target: "AggregationCursorz.html", RoleType: "role", Name: "node-api"}},
	}, {
		input:    ":node-api:`foo <AggregationCursor.html>`",
		expected: []types.RstRole{{Target: "AggregationCursor.html", RoleType: "role", Name: "node-api"}},
	}, {
		input:    "This is a :ref:`valid atlas ref <connect-to-your-cluster>`",
		expected: []types.RstRole{{Target: "connect-to-your-cluster", RoleType: "ref", Name: "ref"}},
	}, {
		input:    "This is a :ref:`valid server ref <replica-set-read-preference-behavior>`",
		expected: []types.RstRole{{Target: "replica-set-read-preference-behavior", RoleType: "ref", Name: "ref"}},
	}, {
		input:    "This is an :ref:`nvalid ref <invalid_ref_sucka-fish>`",
		expected: []types.RstRole{{Target: "invalid_ref_sucka-fish", RoleType: "ref", Name: "ref"}},
	},
	}

	for _, test := range testCases {
		got := ParseForRoles([]byte(test.input))
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

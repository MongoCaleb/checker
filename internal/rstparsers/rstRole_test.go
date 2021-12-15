package rstparsers

import (
	"testing"
)

type roleTestCase struct {
	input    []byte
	expected []RstRole
}

func TestRoleParser(t *testing.T) {
	testCases := []roleTestCase{{
		input:    []byte(""),
		expected: []RstRole{},
	}, {
		input:    []byte(".. _:"),
		expected: []RstRole{},
	}, {
		input:    []byte(".. _: foo"),
		expected: []RstRole{},
	}, {
		input:    []byte("This is a `constant link that should fail <{+api+}/flibbertypoo>`__"),
		expected: []RstRole{},
	}, {
		input:    []byte("This is a `constant link that should succeed <{+api+}/classes/AggregationCursor.html>`__"),
		expected: []RstRole{},
	}, {
		input:    []byte("here is a :ref:`fantastic`"),
		expected: []RstRole{{Target: "fantastic", RoleType: "ref", Name: "ref"}},
	}, {
		input:    []byte("here is a :ref:`fantastic` here is another :ref:`2 <mediocre-fantastic>` here is a :ref:`\n<not_great-fantastic>`"),
		expected: []RstRole{{Target: "fantastic", RoleType: "ref", Name: "ref"}, {Target: "mediocre-fantastic", RoleType: "ref", Name: "ref"}, {Target: "not_great-fantastic", RoleType: "ref", Name: "ref"}},
	}, {
		input:    []byte(":node-api:`foo </AggregationCursor.html>`"),
		expected: []RstRole{{Target: "/AggregationCursor.html", RoleType: "role", Name: "node-api"}},
	}, {
		input:    []byte(":node-api:`foo <AggregationCursorz.html>`"),
		expected: []RstRole{{Target: "AggregationCursorz.html", RoleType: "role", Name: "node-api"}},
	}, {
		input:    []byte(":node-api:`foo <AggregationCursor.html>`"),
		expected: []RstRole{{Target: "AggregationCursor.html", RoleType: "role", Name: "node-api"}},
	}, {
		input:    []byte("This is a :ref:`valid atlas ref <connect-to-your-cluster>`"),
		expected: []RstRole{{Target: "connect-to-your-cluster", RoleType: "ref", Name: "ref"}},
	}, {
		input:    []byte("This is a :ref:`valid server ref <replica-set-read-preference-behavior>`"),
		expected: []RstRole{{Target: "replica-set-read-preference-behavior", RoleType: "ref", Name: "ref"}},
	}, {
		input:    []byte("This is an :ref:`nvalid ref <invalid_ref_sucka-fish>`"),
		expected: []RstRole{{Target: "invalid_ref_sucka-fish", RoleType: "ref", Name: "ref"}},
	}, {
		input:    []byte("This is a `constant link that should fail <{+api+}/flibbertypoo>`__"),
		expected: []RstRole{},
	}, {
		input:    []byte("This is a `constant link that should succeed <{+api+}/classes/AggregationCursor.html>`__"),
		expected: []RstRole{},
	}, {
		input:    []byte("Here is one `constant link <{+api+}/One.html>`__ and a second `constant link <{+api+}/Two.html>`__"),
		expected: []RstRole{},
	},
	}

	for _, test := range testCases {
		got := ParseForRoles(test.input)
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

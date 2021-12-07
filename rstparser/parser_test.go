package rstparser

import (
	"testing"
)

type roleTestCase struct {
	input    string
	expected []RstRole
}

type constantTestCase struct {
	input    string
	expected []RstConstant
}

//input:    "here is a :ref:`fantastic` here is another :ref:`2 <mediocre-fantastic>` here is a :ref:`\n<not_great-fantastic>",
func TestRefParser(t *testing.T) {

	testCases := []roleTestCase{{
		input:    "",
		expected: []RstRole{},
	}, {
		input:    ".. _:",
		expected: []RstRole{},
	}, {
		input:    ".. _: foo",
		expected: []RstRole{},
	}, {
		input:    "This is a `constant link that should fail <{+api+}/flibbertypoo>`__",
		expected: []RstRole{},
	}, {
		input:    "This is a `constant link that should succeed <{+api+}/classes/AggregationCursor.html>`__",
		expected: []RstRole{},
	}, {
		input:    "here is a :ref:`fantastic`",
		expected: []RstRole{{Target: "fantastic", RoleType: "ref", Name: "ref"}},
	}, {
		input:    "here is a :ref:`fantastic` here is another :ref:`2 <mediocre-fantastic>` here is a :ref:`\n<not_great-fantastic>`",
		expected: []RstRole{{Target: "fantastic", RoleType: "ref", Name: "ref"}, {Target: "mediocre-fantastic", RoleType: "ref", Name: "ref"}, {Target: "not_great-fantastic", RoleType: "ref", Name: "ref"}},
	}, {
		input:    ":node-api:`foo </AggregationCursor.html>`",
		expected: []RstRole{{Target: "/AggregationCursor.html", RoleType: "role", Name: "node-api"}},
	}, {
		input:    ":node-api:`foo <AggregationCursorz.html>`",
		expected: []RstRole{{Target: "AggregationCursorz.html", RoleType: "role", Name: "node-api"}},
	}, {
		input:    ":node-api:`foo <AggregationCursor.html>`",
		expected: []RstRole{{Target: "AggregationCursor.html", RoleType: "role", Name: "node-api"}},
	}, {
		input:    "This is a :ref:`valid atlas ref <connect-to-your-cluster>`",
		expected: []RstRole{{Target: "connect-to-your-cluster", RoleType: "ref", Name: "ref"}},
	}, {
		input:    "This is a :ref:`valid server ref <replica-set-read-preference-behavior>`",
		expected: []RstRole{{Target: "replica-set-read-preference-behavior", RoleType: "ref", Name: "ref"}},
	}, {
		input:    "This is an :ref:`nvalid ref <invalid_ref_sucka-fish>`",
		expected: []RstRole{{Target: "invalid_ref_sucka-fish", RoleType: "ref", Name: "ref"}},
	},
	}

	for _, test := range testCases {
		got := roleParse([]byte(test.input))
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
		input:    "here is a :ref:`fantastic` here is another :ref:`2 <mediocre-fantastic>` here is a :ref:`\n<not_great-fantastic>`",
		expected: []RstConstant{},
	}, {
		input:    ":node-api:`foo </AggregationCursor.html>`",
		expected: []RstConstant{},
	}, {
		input:    ":node-api:`foo <AggregationCursorz.html>`",
		expected: []RstConstant{},
	}, {
		input:    ":node-api:`foo <AggregationCursor.html>`",
		expected: []RstConstant{},
	}, {
		input:    "This is a :ref:`valid atlas ref <connect-to-your-cluster>`",
		expected: []RstConstant{},
	}, {
		input:    "This is a :ref:`valid server ref <replica-set-read-preference-behavior>`",
		expected: []RstConstant{},
	}, {
		input:    "This is an :ref:`nvalid ref <invalid_ref_sucka-fish>`",
		expected: []RstConstant{},
	},
	}
	for _, test := range testCases {
		got := constantParse([]byte(test.input))
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

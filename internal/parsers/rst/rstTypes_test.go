package rst

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindLocalRefs(t *testing.T) {
	cases := []struct {
		input    string
		expected []RefTarget
	}{{
		input:    "",
		expected: []RefTarget{},
	}, {
		input:    ".. _:",
		expected: []RefTarget{},
	}, {
		input: ".. _foo:",
		expected: []RefTarget{
			{Target: "foo", Type: "local"},
		},
	}, {
		input: ".. _foo:\n.. _bar:",
		expected: []RefTarget{
			{Target: "foo", Type: "local"},
			{Target: "bar", Type: "local"},
		},
	}, {
		input: ".. _foo:\n.. _bar:\n\n\n\n\n\n.. _baz:",
		expected: []RefTarget{
			{Target: "foo", Type: "local"},
			{Target: "bar", Type: "local"},
			{Target: "baz", Type: "local"},
		},
	}, {
		input:    ".. _version-4.1:",
		expected: []RefTarget{{Raw: "version-4.1", Target: "version-4.1", Type: "local"}},
	},
	}

	for _, c := range cases {
		actual := ParseForLocalRefs([]byte(c.input))
		assert.Equal(t, len(c.expected), len(actual), "FindLocalRefs(%q) should return %d refs, got %d", c.input, len(c.expected), len(actual))

		for i, ref := range actual {
			assert.Equal(t, c.expected[i].Target, ref.Target, "FindLocalRefs(%q) should return ref %d as %q, got %q", c.input, i, c.expected, ref)
			assert.Equal(t, c.expected[i].Type, ref.Type, "FindLocalRefs(%q) should return ref %d as %q, got %q", c.input, i, c.expected, ref)
		}
	}

}

func TestConstantParser(t *testing.T) {

	cases := []struct {
		input    string
		expected []RstConstant
	}{{
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
	for _, test := range cases {
		got := ParseForConstants([]byte(test.input))
		assert.ElementsMatch(t, test.expected, got, "ParseForConstants(%q) should return %v, got %v", test.input, test.expected, got)
	}
}

func TestFindLinkInConstant(t *testing.T) {
	cases := []struct {
		input    RstConstant
		expected bool
	}{{
		input:    RstConstant{Target: "https://www.google.com", Name: "api"},
		expected: true,
	}, {
		input:    RstConstant{Target: "v1.8.0", Name: "api"},
		expected: false,
	}}

	for _, c := range cases {
		actual := c.input.IsHTTPLink()
		assert.Equal(t, c.expected, actual, "IsLink(%q) should return %v, got %v", c.input, c.expected, actual)
	}
}

func TestLinkParser(t *testing.T) {
	cases := []struct {
		input    string
		expected []RstHTTPLink
	}{{
		input:    "",
		expected: []RstHTTPLink{},
	}, {
		input:    "\n\n\n",
		expected: []RstHTTPLink{},
	}, {
		input:    "// code comments \n /* and more comments */ \n // and yet more!",
		expected: []RstHTTPLink{},
	}, {
		input:    "we can say http and www without any links being found",
		expected: []RstHTTPLink{},
	}, {
		input:    "markdown links are found\n\t\t [some markdown link](https://www.google.com)\\n\" +\n\t\t\"   [some other link](https://a.bad.url)\\n\" +",
		expected: []RstHTTPLink{RstHTTPLink("https://www.google.com"), RstHTTPLink("https://a.bad.url")},
	}, {
		input:    "http links in rst are found\n\t\t\"   this is a bad `url <https://www.flibbertypip.com>`__\\n\" +\n\t\t\"   this is a good `url <https://www.github.com>`__",
		expected: []RstHTTPLink{RstHTTPLink("https://www.flibbertypip.com"), RstHTTPLink("https://www.github.com")},
	},
	}
	for _, test := range cases {
		got := ParseForHTTPLinks([]byte(test.input))
		assert.ElementsMatch(t, test.expected, got, "ParseForConstants(%q) should return %v, got %v", test.input, test.expected, got)
	}
}

//go:embed testdata/makesGoUnhappy.txt
var edge []byte

func TestRoleParser(t *testing.T) {
	cases := []struct {
		input    []byte
		expected []RstRole
	}{{
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
	}, {
		input:    edge,
		expected: []RstRole{{Target: "/reference/operator/update/positional-filtered/", RoleType: "role", Name: "manual"}},
	}, {
		input:    []byte("here is a :ref:`fantastic`"),
		expected: []RstRole{{Target: "fantastic", RoleType: "ref", Name: "ref"}},
	}, {
		input:    []byte(":ref:`What information does the MongoDB Compatibility table show? <mongodb-compatibility-table-about-node>`"),
		expected: []RstRole{{Target: "mongodb-compatibility-table-about-node", RoleType: "ref", Name: "ref"}},
	}}

	for _, test := range cases {
		got := ParseForRoles(test.input)
		assert.ElementsMatch(t, test.expected, got, "ParseForConstants(%q) should return %v, got %v", test.input, test.expected, got)
	}
}

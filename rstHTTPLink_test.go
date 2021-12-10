package main

import (
	"testing"
)

const (
	linkInput = "\n" +
		"\n" +
		"\n" +
		"// a code comment\n" +
		"// another comment\n" +
		"\n" +
		"here is a :ref:`fantastic` here is another :ref:`2 <mediocre-fantastic>` here is a :ref:`\n" +
		"<not_great-fantastic>`\n" +
		"\n" +
		"   http://a.bad.url\n" +
		"\n" +
		"   [some markdown link](https://www.google.com)\n" +
		"   [some other link](https://a.bad.url)\n" +
		"\n" +
		"   this is a bad `url <https://www.flibbertypoo.com>`__\n" +
		"   this is a good `url <https://www.github.com>`__\n" +
		"\n" +
		"   * - D\n" +
		"     - | `Vibe.D native MongoDB driver <http://vibed.org/docs#mongo>`__\n" +
		"\n" +
		"   * - Dart\n" +
		"     - `mongo_dart 0.4.0 <http://pub.dartlang.org/packages/mongo_dart>`__"
)

func contains(s []RstHTTPLink, str RstHTTPLink) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func TestLinkParser(t *testing.T) {

	expected := []RstHTTPLink{
		"http://a.bad.url",
		"https://www.google.com",
		"https://www.flibbertypoo.com",
		"https://a.bad.url",
		"https://www.github.com",
		"http://vibed.org/docs#mongo",
		"http://pub.dartlang.org/packages/mongo_dart",
	}
	actual := ParseForLinks(linkInput)
	if len(actual) != len(expected) {
		t.Errorf("Expected %d files, got %d", len(expected), len(actual))
	}
	for _, link := range actual {
		if !contains(expected, link) {
			t.Errorf("Expected to find link %s in %s", link, expected)
		}
	}
}

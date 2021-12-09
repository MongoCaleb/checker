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

func contains(s []HTTPLink, str HTTPLink) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func makeLinkMap(name string) *LinkMap {
	return &LinkMap{
		Files: []RstFile{{
			Name: name,
			Links: []HTTPLink{
				"http://a.bad.url",
				"https://www.google.com",
				"https://www.flibbertypoo.com",
				"https://a.bad.url",
				"https://www.github.com",
				"http://vibed.org/docs#mongo",
				"http://pub.dartlang.org/packages/mongo_dart",
			},
		}},
	}
}

func TestLinkParser(t *testing.T) {

	expected := makeLinkMap("testFile")
	actual := NewLinkMap()
	actual.ParseForLinks("testFile", linkInput)
	if len(actual.Files) != len(expected.Files) {
		t.Errorf("Expected %d files, got %d", len(expected.Files), len(actual.Files))
	}
	if actual.Files[0].Name != expected.Files[0].Name {
		t.Errorf("Expected %s, got %s", expected.Files[0].Name, actual.Files[0].Name)
	}
	if len(actual.Files[0].Links) != len(expected.Files[0].Links) {
		t.Errorf("Expected %d links, got %d", len(expected.Files[0].Links), len(actual.Files[0].Links))
	}
	for _, link := range actual.Files[0].Links {
		if !contains(expected.Files[0].Links, link) {
			t.Errorf("Expected to find link %s in %s", link, expected.Files[0].Links)
		}
	}
}

func TestMultipleLinkFiles(t *testing.T) {

	expected := makeLinkMap("testFile1")
	expected.Files = append(expected.Files, makeLinkMap("testFile2").Files...)
	expected.Files = append(expected.Files, makeLinkMap("testFile3").Files...)

	actual := NewLinkMap()
	actual.ParseForLinks("testFile1", linkInput)
	actual.ParseForLinks("testFile2", linkInput)
	actual.ParseForLinks("testFile3", linkInput)

	if len(actual.Files) != len(expected.Files) {
		t.Errorf("Expected %d files, got %d", len(expected.Files), len(actual.Files))
	}

	for i := range actual.Files {
		if actual.Files[i].Name != expected.Files[i].Name {
			t.Errorf("Expected %s, got %s", expected.Files[i].Name, actual.Files[i].Name)
		}

		if len(actual.Files[i].Links) != len(expected.Files[i].Links) {
			t.Errorf("Expected %d links, got %d", len(expected.Files[i].Links), len(actual.Files[i].Links))
		}

		for _, link := range actual.Files[i].Links {
			if !contains(expected.Files[i].Links, link) {
				t.Errorf("Expected to find link %s in %s", link, expected.Files[0].Links)
			}
		}
	}

}

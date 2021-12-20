package intersphinx

import (
	"bytes"
	"checker/internal/parsers/rst"
	"compress/zlib"
	"io/ioutil"
	"log"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNoContents(t *testing.T) {

	logrus.SetOutput(ioutil.Discard)

	header := ""

	resp := Intersphinx([]byte(header), "test")
	assert.Nil(t, resp, "Expected nil, got %v", resp)
}

func TestInvalidHeader(t *testing.T) {

	logrus.SetOutput(ioutil.Discard)
	header := []byte(`# Sphinx inventory version 2
# Project: golang
# Version:
`)
	resp := Intersphinx(header, "test")
	assert.Nil(t, resp, "Expected nil, got %v", resp)
}
func TestHeaderNoContent(t *testing.T) {

	header := []byte(`# Sphinx inventory version 2
# Project: golang
# Version:
# The remainder of this file is compressed using zlib.
`)

	resp := Intersphinx(header, "test")
	assert.Nil(t, resp, "Expected nil, got %v", resp)
}

func TestInvalidContent(t *testing.T) {
	logrus.SetOutput(ioutil.Discard)

	header := []byte(`# Sphinx inventory version 2
# Project: golang
# Version:
# The remainder of this file is compressed using zlib.
`)
	zText := []byte(`whats-new std:doc -1 whats-new/ What's New
compatibility std:doc -1 compatibility/ Compatibility
fundamentals std:doc -1 fundamentals/ Fundamentals
usage-examples std:doc -1 usage-examples/ Usage Examples`)

	resp := Intersphinx(append(header, zText...), "test")

	assert.Nil(t, resp, "Expected nil, got %v", resp)
}

func TestSomeContent(t *testing.T) {
	logrus.SetOutput(ioutil.Discard)

	header := []byte(`# Sphinx inventory version 2
# Project: golang
# Version:
# The remainder of this file is compressed using zlib.
`)
	zText := []byte(`whats-new std:doc -1 whats-new/ What's New
compatibility std:doc -1 compatibility/ Compatibility
fundamentals std:doc -1 fundamentals/ Fundamentals
usage-examples std:doc -1 usage-examples/ Usage Examples`)

	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	if _, err := w.Write(zText); err != nil {
		log.Fatal(err)
	}
	w.Close()

	resp := Intersphinx(append(header, b.Bytes()...), "https://test.com/")

	expected := SphinxMap{
		"https://test.com/": {
			"whats-new":      rst.RefTarget{Target: "https://test.com/whats-new/%s", Type: "intersphinx"},
			"compatibility":  rst.RefTarget{Target: "https://test.com/compatibility/%s", Type: "intersphinx"},
			"fundamentals":   rst.RefTarget{Target: "https://test.com/fundamentals/%s", Type: "intersphinx"},
			"usage-examples": rst.RefTarget{Target: "https://test.com/usage-examples/%s", Type: "intersphinx"},
		},
	}

	assert.EqualValues(t, expected, resp, "Expected %v, got %v", expected, resp)
}

func TestJoinSphinxes(t *testing.T) {
	input := []SphinxMap{{
		"https://test1.com/": {
			"whats-new":      rst.RefTarget{Target: "https://test1.com/whats-new/%s", Type: "intersphinx"},
			"compatibility":  rst.RefTarget{Target: "https://test1.com/compatibility/%s", Type: "intersphinx"},
			"fundamentals":   rst.RefTarget{Target: "https://test1.com/fundamentals/%s", Type: "intersphinx"},
			"usage-examples": rst.RefTarget{Target: "https://test1.com/usage-examples/%s", Type: "intersphinx"},
		},
		"https://test2.com/": {
			"whats-new":      rst.RefTarget{Target: "https://test2.com/whats-new/%s", Type: "intersphinx"},
			"compatibility":  rst.RefTarget{Target: "https://test2.com/compatibility/%s", Type: "intersphinx"},
			"fundamentals":   rst.RefTarget{Target: "https://test2.com/fundamentals/%s", Type: "intersphinx"},
			"usage-examples": rst.RefTarget{Target: "https://test2.com/usage-examples/%s", Type: "intersphinx"},
		},
	}}

	expected := SphinxMap{
		"https://test1.com/": {
			"whats-new":      rst.RefTarget{Target: "https://test1.com/whats-new/%s", Type: "intersphinx"},
			"compatibility":  rst.RefTarget{Target: "https://test1.com/compatibility/%s", Type: "intersphinx"},
			"fundamentals":   rst.RefTarget{Target: "https://test1.com/fundamentals/%s", Type: "intersphinx"},
			"usage-examples": rst.RefTarget{Target: "https://test1.com/usage-examples/%s", Type: "intersphinx"},
		},
		"https://test2.com/": {
			"whats-new":      rst.RefTarget{Target: "https://test2.com/whats-new/%s", Type: "intersphinx"},
			"compatibility":  rst.RefTarget{Target: "https://test2.com/compatibility/%s", Type: "intersphinx"},
			"fundamentals":   rst.RefTarget{Target: "https://test2.com/fundamentals/%s", Type: "intersphinx"},
			"usage-examples": rst.RefTarget{Target: "https://test2.com/usage-examples/%s", Type: "intersphinx"},
		},
	}

	actual := JoinSphinxes(input)

	assert.EqualValues(t, expected, actual, "expected %v, got %v", expected, actual)
}

package intersphinx

import (
	"bytes"
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
		"whats-new":      true,
		"compatibility":  true,
		"fundamentals":   true,
		"usage-examples": true,
	}

	assert.EqualValues(t, expected, resp, "Expected %v, got %v", expected, resp)
}

func TestJoinSphinxes(t *testing.T) {
	input := []SphinxMap{
		{
			"whats-new":      true,
			"compatibility":  true,
			"fundamentals":   true,
			"usage-examples": true,
		}, {
			"foowhats-new":      true,
			"foocompatibility":  true,
			"foofundamentals":   true,
			"foousage-examples": true,
		}}

	expected := SphinxMap{
		"whats-new":         true,
		"compatibility":     true,
		"fundamentals":      true,
		"usage-examples":    true,
		"foowhats-new":      true,
		"foocompatibility":  true,
		"foofundamentals":   true,
		"foousage-examples": true,
	}

	actual := JoinSphinxes(input)

	assert.EqualValues(t, expected, actual, "expected %v, got %v", expected, actual)
}

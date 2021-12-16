package intersphinx

import (
	"bytes"
	"checker/internal/parsers/rst"
	"compress/zlib"
	"io/ioutil"
	"log"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestNoContents(t *testing.T) {

	logrus.SetOutput(ioutil.Discard)

	header := ""

	resp := Intersphinx([]byte(header), "test")
	if resp != nil {
		t.Errorf("Expected nil, got %v", resp)
	}

}

func TestInvalidHeader(t *testing.T) {

	logrus.SetOutput(ioutil.Discard)
	header := []byte(`# Sphinx inventory version 2
# Project: golang
# Version:
`)
	resp := Intersphinx(header, "test")
	if resp != nil {
		t.Errorf("Expected nil, got %v", resp)
	}

}
func TestHeaderNoContent(t *testing.T) {

	header := []byte(`# Sphinx inventory version 2
# Project: golang
# Version:
# The remainder of this file is compressed using zlib.
`)

	resp := Intersphinx(header, "test")
	if resp != nil {
		t.Errorf("Expected nil, got %v", resp)
	}

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

	if resp != nil {
		t.Errorf("Expected nil, got %v", resp)
	}

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

	resp := Intersphinx(append(header, b.Bytes()...), "https://docs.mongodb.com/drivers/go/current/")

	if len(resp) != 4 {
		t.Errorf("Expected 4 entries, got %v", len(resp))
	}

	expected := RefMap{
		"whats-new":      rst.LocalRef{Target: "https://docs.mongodb.com/drivers/go/current/whats-new/", Type: "intersphinx"},
		"compatibility":  rst.LocalRef{Target: "https://docs.mongodb.com/drivers/go/current/compatibility/", Type: "intersphinx"},
		"fundamentals":   rst.LocalRef{Target: "https://docs.mongodb.com/drivers/go/current/fundamentals/", Type: "intersphinx"},
		"usage-examples": rst.LocalRef{Target: "https://docs.mongodb.com/drivers/go/current/usage-examples/", Type: "intersphinx"},
	}

	for k, v := range resp {
		if v != expected[k] {
			t.Errorf("Expected %v, got %v", expected[k], v)
		}
	}

}

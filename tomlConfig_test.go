package main

import (
	"testing"
)

const tomlConfigInput = `
name = "this is a test"
title = "TEST"
toc_landing_pages = [
    "/fundamentals/connection",
    "/fundamentals/crud",
    "/usage-examples"
]

intersphinx = [
  "https://docs.mongodb.com/manual/objects.inv",
  "https://docs.atlas.mongodb.com/objects.inv",
]

[constants]
docs-branch = "master" # always set this to the docs branch (i.e. master, 1.7, 1.8, etc.)
version = "v1.8.0" # always set this to the driver branch (i.e. v1.7.0, v1.8.0, etc.)
example = "https://raw.githubusercontent.com/mongodb/docs-golang/{+docs-branch+}/source/includes/usage-examples/code-snippets"
api = "https://pkg.go.dev/go.mongodb.org/mongo-driver@{+version+}"
foo = "{+bar+}@{+version+}"
bar = "{+docs-branch+}"
baz = "baz+{+foo+}"
`

func TestSnootyToml(t *testing.T) {
	cfg, err := NewTomlConfig(tomlConfigInput)
	if err != nil {
		t.Errorf("error parsing toml: %v", err)
	}

	if cfg.Name != "this is a test" {
		t.Errorf("expected name to be 'go', got %s", cfg.Name)
	}

	intersphinxes := []string{"https://docs.mongodb.com/manual/objects.inv", "https://docs.atlas.mongodb.com/objects.inv"}
	for i, intersphinx := range cfg.Intersphinx {
		if intersphinx != intersphinxes[i] {
			t.Errorf("expected intersphinx to be %s, got %s", intersphinxes[i], intersphinx)
		}
	}
}

func TestConstantResolution(t *testing.T) {
	cfg, err := NewTomlConfig(tomlConfigInput)
	if err != nil {
		t.Errorf("error parsing toml: %v", err)
	}
	expected := map[string]string{"docs-branch": "master", "version": "v1.8.0", "example": "https://raw.githubusercontent.com/mongodb/docs-golang/master/source/includes/usage-examples/code-snippets", "api": "https://pkg.go.dev/go.mongodb.org/mongo-driver@v1.8.0", "foo": "master@v1.8.0", "bar": "master", "baz": "baz+master@v1.8.0"}
	for k, v := range expected {
		if cfg.Constants[k] != v {
			t.Errorf("expected %s to be %s, got %s", k, v, cfg.Constants[k])
		}
	}
}

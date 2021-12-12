package collectors

import (
	"io"
	"path/filepath"
	"runtime"
	"testing"

	log "github.com/sirupsen/logrus"
	iowrap "github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

var basepath string

func init() {
	FS = iowrap.NewMemMapFs()
	FSUtil = &iowrap.Afero{Fs: FS}
	_, b, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("Could not get caller")
	}
	basepath = filepath.Dir(b)

}

func afterTest(t *testing.T) {
	t.Cleanup(func() {
		FS.RemoveAll(basepath)
	})

}

func TestSnootyTomlNonExist(t *testing.T) {
	defer afterTest(t)
	log.SetOutput(io.Discard)

	assert.False(t, snootyTomlExists(), "Snooty.toml should not exist")
}

func TestChecksIfSnootyTomlExists(t *testing.T) {
	defer afterTest(t)

	iowrap.WriteFile(FS, filepath.Join(basepath, "snooty.toml"), []byte(""), 0644)

	assert.True(t, snootyTomlExists(), "Snooty.toml should exist")
}

func TestFailsIfNoSourceDirectory(t *testing.T) {
	defer afterTest(t)
	log.SetOutput(io.Discard)
	assert.False(t, sourceDirectoryExists(), "Source directory should not exist")
}

func TestFindsSourceDirectory(t *testing.T) {
	defer afterTest(t)
	log.SetOutput(io.Discard)

	FS.MkdirAll(filepath.Join(basepath, "source/"), 0755)
	iowrap.WriteFile(FS, filepath.Join(basepath, "snooty.toml"), []byte("test"), 0644)

	assert.True(t, sourceDirectoryExists(), "Source directory found")

}

func TestGatherXPanicsIfNoSourceOrSnootyToml(t *testing.T) {
	defer afterTest(t)
	log.SetOutput(io.Discard)
	assert.Panics(t, func() { gatherFiles() }, "gatherRole should panic if no source or Snooty.toml")
}

func TestGatherFiles(t *testing.T) {
	defer afterTest(t)

	FS.MkdirAll(filepath.Join(basepath, "source"), 0755)
	FS.MkdirAll(filepath.Join(basepath, "source", "fundamentals"), 0755)
	iowrap.WriteFile(FS, filepath.Join(basepath, "snooty.toml"), []byte("test"), 0644)

	iowrap.WriteFile(FS, filepath.Join(basepath, "source", "foo.txt"), []byte("test"), 0644)
	iowrap.WriteFile(FS, filepath.Join(basepath, "source", "bar.txt"), []byte("test"), 0644)
	iowrap.WriteFile(FS, filepath.Join(basepath, "source", "fundamentals", "baz.txt"), []byte("test"), 0644)
	iowrap.WriteFile(FS, filepath.Join(basepath, "source", "fundamentals", "biz.txt"), []byte("test"), 0644)

	expected := []string{filepath.Join(basepath, "source", "foo.txt"), filepath.Join(basepath, "source", "bar.txt"), filepath.Join(basepath, "source", "fundamentals", "baz.txt"), filepath.Join(basepath, "source", "fundamentals", "biz.txt")}
	actual := gatherFiles()

	assert.ElementsMatch(t, expected, actual, "gatherFiles should return all files in source directory")

}

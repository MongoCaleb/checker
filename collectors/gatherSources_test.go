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

func init() {
	FS = iowrap.NewMemMapFs()
	FSUtil = &iowrap.Afero{Fs: FS}
}

func TestSnootyTomlNonExist(t *testing.T) {
	log.SetOutput(io.Discard)

	assert.False(t, snootyTomlExists(), "Snooty.toml should not exist")
}

func TestChecksIfSnootyTomlExists(t *testing.T) {
	_, b, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("Could not get caller")
	}
	basepath := filepath.Dir(b)

	FS.MkdirAll(filepath.Join(basepath, "source"), 0755)
	iowrap.WriteFile(FS, filepath.Join(basepath, "snooty.toml"), []byte(""), 0644)

	assert.True(t, snootyTomlExists(), "Snooty.toml should exist")
}

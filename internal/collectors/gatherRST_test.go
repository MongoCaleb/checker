package collectors

import (
	"checker/internal/rstparsers"
	_ "embed"
	"io"
	"path/filepath"
	"testing"

	log "github.com/sirupsen/logrus"
	iowrap "github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

var (
	//go:embed testdata/index.txt
	indexFile string

	//go:embed testdata/aggregation.txt
	aggregationsFile string
)

func init() {
	FS = iowrap.NewMemMapFs()
	FSUtil = &iowrap.Afero{Fs: FS}

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

func TestGatherRoles(t *testing.T) {
	defer afterTest(t)

	FS.MkdirAll(filepath.Join(basepath, "source"), 0755)
	FS.MkdirAll(filepath.Join(basepath, "source", "fundamentals"), 0755)
	iowrap.WriteFile(FS, filepath.Join(basepath, "snooty.toml"), []byte("test"), 0644)

	iowrap.WriteFile(FS, filepath.Join(basepath, "source", "index.txt"), []byte(indexFile), 0644)
	iowrap.WriteFile(FS, filepath.Join(basepath, "source", "fundamentals", "aggregation.txt"), []byte(aggregationsFile), 0644)

	expected := map[string][]rstparsers.RstRole{
		"/source/fundamentals/aggregation.txt": {
			{Target: "/core/aggregation-pipeline-limits/", RoleType: "role", Name: "manual"},
			{Target: "/reference/limits/#mongodb-limit-BSON-Document-Size", RoleType: "role", Name: "manual"},
			{Target: "/reference/operator/aggregation/graphLookup/", RoleType: "role", Name: "manual"}, {Target: "/reference/operator/aggregation/", RoleType: "role", Name: "manual"},
			{Target: "/core/aggregation-pipeline/", RoleType: "role", Name: "manual"},
			{Target: "/meta/aggregation-quick-reference/#stages", RoleType: "role", Name: "manual"},
			{Target: "/meta/aggregation-quick-reference/#operator-expressions", RoleType: "role", Name: "manual"},
			{Target: "/fundamentals/connection", RoleType: "role", Name: "doc"},
			{Target: "/reference/operator/aggregation/match/", RoleType: "role", Name: "manual"},
			{Target: "/reference/operator/aggregation/group/", RoleType: "role", Name: "manual"},
		},
		"/source/index.txt": {
			{Target: "/quick-start", RoleType: "role", Name: "doc"},
			{Target: "/usage-examples", RoleType: "role", Name: "doc"},
			{Target: "/faq", RoleType: "role", Name: "doc"},
			{Target: "/issues-and-help", RoleType: "role", Name: "doc"},
			{Target: "/compatibility", RoleType: "role", Name: "doc"},
			{Target: "What", RoleType: "role", Name: "doc"}},
	}

	actual := gatherRoles(gatherFiles())

	assert.EqualValues(t, expected, actual, "gatherRoles should return all roles in source directory")

}

func TestGatherConstants(t *testing.T) {
	defer afterTest(t)

	FS.MkdirAll(filepath.Join(basepath, "source"), 0755)
	FS.MkdirAll(filepath.Join(basepath, "source", "fundamentals"), 0755)
	iowrap.WriteFile(FS, filepath.Join(basepath, "snooty.toml"), []byte("test"), 0644)

	iowrap.WriteFile(FS, filepath.Join(basepath, "source", "index.txt"), []byte(indexFile), 0644)
	iowrap.WriteFile(FS, filepath.Join(basepath, "source", "fundamentals", "aggregation.txt"), []byte(aggregationsFile), 0644)

	expected := map[string][]rstparsers.RstConstant{
		"/source/fundamentals/aggregation.txt": {
			{Name: "api", Target: "/interfaces/AggregateOptions.html"},
			{Name: "api", Target: "/classes/Collection.html#aggregate"},
		},
		"/source/index.txt": {},
	}

	actual := gatherConstants(gatherFiles())

	assert.EqualValues(t, expected, actual, "gatherConstants should return all constants in source directory")

}

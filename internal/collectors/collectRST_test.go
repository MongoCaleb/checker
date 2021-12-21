package collectors

import (
	"checker/internal/parsers/rst"
	_ "embed"
	"io"
	"os"
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

	//go:embed testdata/gridfs.txt
	grifsFile string
)

func init() {
	FS = iowrap.NewMemMapFs()
	FSUtil = &iowrap.Afero{Fs: FS}
	basepath, _ = os.Getwd()

}

func check(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func afterTest(t *testing.T) {
	t.Cleanup(func() {
		if err := FS.RemoveAll(basepath); err != nil {
			log.Fatal(err)
		}
	})

}

func TestSnootyTomlNonExist(t *testing.T) {
	defer afterTest(t)
	log.SetOutput(io.Discard)

	assert.False(t, snootyTomlExists(basepath), "Snooty.toml should not exist")
}

func TestChecksIfSnootyTomlExists(t *testing.T) {
	defer afterTest(t)

	check(iowrap.WriteFile(FS, filepath.Join(basepath, "snooty.toml"), []byte(""), 0644))

	assert.True(t, snootyTomlExists(basepath), "Snooty.toml should exist")
}

func TestFailsIfNoSourceDirectory(t *testing.T) {
	defer afterTest(t)
	log.SetOutput(io.Discard)
	assert.False(t, sourceDirectoryExists(basepath), "Source directory should not exist")
}

func TestFindsSourceDirectory(t *testing.T) {
	defer afterTest(t)
	log.SetOutput(io.Discard)

	check(FS.MkdirAll(filepath.Join(basepath, "source/"), 0755))
	check(iowrap.WriteFile(FS, filepath.Join(basepath, "snooty.toml"), []byte("test"), 0644))

	assert.True(t, sourceDirectoryExists(basepath), "Source directory found")

}

func TestGatherXPanicsIfNoSourceOrSnootyToml(t *testing.T) {
	defer afterTest(t)
	log.SetOutput(io.Discard)
	assert.Panics(t, func() { GatherFiles(basepath) }, "gatherRole should panic if no source or Snooty.toml")
}

func TestGatherFiles(t *testing.T) {
	defer afterTest(t)

	check(FS.MkdirAll(filepath.Join(basepath, "source"), 0755))
	check(FS.MkdirAll(filepath.Join(basepath, "source", "fundamentals"), 0755))
	check(iowrap.WriteFile(FS, filepath.Join(basepath, "snooty.toml"), []byte("test"), 0644))
	check(iowrap.WriteFile(FS, filepath.Join(basepath, "source", "foo.txt"), []byte("test"), 0644))
	check(iowrap.WriteFile(FS, filepath.Join(basepath, "source", "bar.txt"), []byte("test"), 0644))
	check(iowrap.WriteFile(FS, filepath.Join(basepath, "source", "fundamentals", "baz.txt"), []byte("test"), 0644))
	check(iowrap.WriteFile(FS, filepath.Join(basepath, "source", "fundamentals", "biz.txt"), []byte("test"), 0644))
	expected := []string{filepath.Join(basepath, "source", "foo.txt"), filepath.Join(basepath, "source", "bar.txt"), filepath.Join(basepath, "source", "fundamentals", "baz.txt"), filepath.Join(basepath, "source", "fundamentals", "biz.txt")}
	actual := GatherFiles(basepath)

	assert.ElementsMatch(t, expected, actual, "gatherFiles should return all files in source directory")

}

func TestGatherRoles(t *testing.T) {
	defer afterTest(t)

	check(FS.MkdirAll(filepath.Join(basepath, "source"), 0755))
	check(FS.MkdirAll(filepath.Join(basepath, "source", "fundamentals"), 0755))
	check(iowrap.WriteFile(FS, filepath.Join(basepath, "snooty.toml"), []byte("test"), 0644))
	check(iowrap.WriteFile(FS, filepath.Join(basepath, "source", "index.txt"), []byte(indexFile), 0644))
	check(iowrap.WriteFile(FS, filepath.Join(basepath, "source", "fundamentals", "aggregation.txt"), []byte(aggregationsFile), 0644))
	check(iowrap.WriteFile(FS, filepath.Join(basepath, "source", "fundamentals", "gridfs.txt"), []byte(grifsFile), 0644))

	expected := map[rst.RstRole]string{
		{Target: "/classes/gridfsbucket.html#drop", RoleType: "role", Name: "node-api-4.0"}:                   "/source/fundamentals/gridfs.txt",
		{Target: "/compatibility", RoleType: "role", Name: "doc"}:                                             "/source/index.txt",
		{Target: "/core/aggregation-pipeline-limits/", RoleType: "role", Name: "manual"}:                      "/source/fundamentals/aggregation.txt",
		{Target: "/core/aggregation-pipeline/", RoleType: "role", Name: "manual"}:                             "/source/fundamentals/aggregation.txt",
		{Target: "/core/gridfs", RoleType: "role", Name: "manual"}:                                            "/source/fundamentals/gridfs.txt",
		{Target: "/core/gridfs/#gridfs-indexes", RoleType: "role", Name: "manual"}:                            "/source/fundamentals/gridfs.txt",
		{Target: "/faq", RoleType: "role", Name: "doc"}:                                                       "/source/index.txt",
		{Target: "/fundamentals/connection", RoleType: "role", Name: "doc"}:                                   "/source/fundamentals/aggregation.txt",
		{Target: "/fundamentals/crud/read-operations/", RoleType: "role", Name: "doc"}:                        "/source/fundamentals/gridfs.txt",
		{Target: "/fundamentals/crud/read-operations/cursor", RoleType: "role", Name: "doc"}:                  "/source/fundamentals/gridfs.txt",
		{Target: "/issues-and-help", RoleType: "role", Name: "doc"}:                                           "/source/index.txt",
		{Target: "/meta/aggregation-quick-reference/#operator-expressions", RoleType: "role", Name: "manual"}: "/source/fundamentals/aggregation.txt",
		{Target: "/meta/aggregation-quick-reference/#stages", RoleType: "role", Name: "manual"}:               "/source/fundamentals/aggregation.txt",
		{Target: "/quick-start", RoleType: "role", Name: "doc"}:                                               "/source/index.txt",
		{Target: "/reference/limits/#mongodb-limit-BSON-Document-Size", RoleType: "role", Name: "manual"}:     "/source/fundamentals/aggregation.txt",
		{Target: "/reference/operator/aggregation/", RoleType: "role", Name: "manual"}:                        "/source/fundamentals/aggregation.txt",
		{Target: "/reference/operator/aggregation/graphLookup/", RoleType: "role", Name: "manual"}:            "/source/fundamentals/aggregation.txt",
		{Target: "/reference/operator/aggregation/group/", RoleType: "role", Name: "manual"}:                  "/source/fundamentals/aggregation.txt",
		{Target: "/reference/operator/aggregation/match/", RoleType: "role", Name: "manual"}:                  "/source/fundamentals/aggregation.txt",
		{Target: "/usage-examples", RoleType: "role", Name: "doc"}:                                            "/source/index.txt",
		{Target: "What", RoleType: "role", Name: "doc"}:                                                       "/source/index.txt",
		{Target: "classes/gridfsbucket.html#delete", RoleType: "role", Name: "node-api-4.0"}:                  "/source/fundamentals/gridfs.txt",
		{Target: "classes/gridfsbucket.html#rename", RoleType: "role", Name: "node-api-4.0"}:                  "/source/fundamentals/gridfs.txt",
		{Target: "gridfs-create-bucket", RoleType: "ref", Name: "ref"}:                                        "/source/fundamentals/gridfs.txt",
		{Target: "gridfs-delete-bucket", RoleType: "ref", Name: "ref"}:                                        "/source/fundamentals/gridfs.txt",
		{Target: "gridfs-delete-files", RoleType: "ref", Name: "ref"}:                                         "/source/fundamentals/gridfs.txt",
		{Target: "gridfs-download-files", RoleType: "ref", Name: "ref"}:                                       "/source/fundamentals/gridfs.txt",
		{Target: "gridfs-rename-files", RoleType: "ref", Name: "ref"}:                                         "/source/fundamentals/gridfs.txt",
		{Target: "gridfs-retrieve-file-info", RoleType: "ref", Name: "ref"}:                                   "/source/fundamentals/gridfs.txt",
		{Target: "gridfs-upload-files", RoleType: "ref", Name: "ref"}:                                         "/source/fundamentals/gridfs.txt",
	}

	actual := GatherRoles(GatherFiles(basepath))

	assert.EqualValues(t, expected, actual, "gatherRoles should return all roles in source directory")

}

func TestRstRoleMapGet(t *testing.T) {
	targets := []rst.RstRole{
		{Target: "gridfs-delete-files", RoleType: "ref", Name: "ref"},
		{Target: "gridfs-create-bucket", RoleType: "ref", Name: "ref"},
		{Target: "gridfs-delete-bucket", RoleType: "ref", Name: "ref"},
		{Target: "gridfs-delete-files", RoleType: "ref", Name: "ref"},
		{Target: "gridfs-download-files", RoleType: "ref", Name: "ref"},
		{Target: "gridfs-rename-files", RoleType: "ref", Name: "ref"},
		{Target: "gridfs-retrieve-file-info", RoleType: "ref", Name: "ref"},
		{Target: "gridfs-upload-files", RoleType: "ref", Name: "ref"},
	}

	localRefs := RefTargetMap{
		{Target: "gridfs-create-bucket", Type: "local"}:        "/source/fundamentals/gridfs.txt",
		{Target: "gridfs-delete-bucket", Type: "local"}:        "/source/fundamentals/gridfs.txt",
		{Target: "gridfs-delete-files", Type: "local"}:         "/source/fundamentals/gridfs.txt",
		{Target: "gridfs-download-files", Type: "local"}:       "/source/fundamentals/gridfs.txt",
		{Target: "gridfs-rename-files", Type: "local"}:         "/source/fundamentals/gridfs.txt",
		{Target: "gridfs-retrieve-file-info", Type: "local"}:   "/source/fundamentals/gridfs.txt",
		{Target: "gridfs-upload-files", Type: "local"}:         "/source/fundamentals/gridfs.txt",
		{Target: "nodejs-aggregation-overview", Type: "local"}: "/source/fundamentals/aggregation.txt",
	}
	for _, target := range targets {
		_, ok := localRefs.Get(&target)
		assert.True(t, ok, "localRefs should contain %s", target.Target)
	}
}

func TestGatherConstants(t *testing.T) {
	defer afterTest(t)

	check(FS.MkdirAll(filepath.Join(basepath, "source"), 0755))
	check(FS.MkdirAll(filepath.Join(basepath, "source", "fundamentals"), 0755))
	check(iowrap.WriteFile(FS, filepath.Join(basepath, "snooty.toml"), []byte("test"), 0644))
	check(iowrap.WriteFile(FS, filepath.Join(basepath, "source", "index.txt"), []byte(indexFile), 0644))
	check(iowrap.WriteFile(FS, filepath.Join(basepath, "source", "fundamentals", "aggregation.txt"), []byte(aggregationsFile), 0644))

	expected := map[rst.RstConstant]string{
		{Name: "api", Target: "/classes/Collection.html#aggregate"}: "/source/fundamentals/aggregation.txt",
		{Name: "api", Target: "/interfaces/AggregateOptions.html"}:  "/source/fundamentals/aggregation.txt",
	}

	actual := GatherConstants(GatherFiles(basepath))

	assert.EqualValues(t, expected, actual, "gatherConstants should return all constants in source directory")

}
func TestGatherHTTPLinks(t *testing.T) {
	defer afterTest(t)

	check(FS.MkdirAll(filepath.Join(basepath, "source"), 0755))
	check(FS.MkdirAll(filepath.Join(basepath, "source", "fundamentals"), 0755))
	check(iowrap.WriteFile(FS, filepath.Join(basepath, "snooty.toml"), []byte("test"), 0644))
	check(iowrap.WriteFile(FS, filepath.Join(basepath, "source", "index.txt"), []byte(indexFile), 0644))
	check(iowrap.WriteFile(FS, filepath.Join(basepath, "source", "fundamentals", "aggregation.txt"), []byte(aggregationsFile), 0644))

	expected := map[rst.RstHTTPLink]string{
		"https://developer.mongodb.com/community/forums/tag/node-js":                                                         "/source/index.txt",
		"https://developer.mongodb.com/learn/?content=Articles&text=Node.js":                                                 "/source/index.txt",
		"https://github.com/mongodb/node-mongodb-native/":                                                                    "/source/index.txt",
		"https://github.com/mongodb/node-mongodb-native/releases/":                                                           "/source/index.txt",
		"https://university.mongodb.com/courses/M220JS/about":                                                                "/source/index.txt",
		"https://www.mongodb.com/blog/post/quick-start-nodejs--mongodb--how-to-analyze-data-using-the-aggregation-framework": "/source/fundamentals/aggregation.txt",
	}

	actual := GatherHTTPLinks(GatherFiles(basepath))

	assert.EqualValues(t, expected, actual, "gatherConstants should return all constants in source directory")

}

func TestGatherLocalRefs(t *testing.T) {
	defer afterTest(t)

	check(FS.MkdirAll(filepath.Join(basepath, "source"), 0755))
	check(FS.MkdirAll(filepath.Join(basepath, "source", "fundamentals"), 0755))
	check(iowrap.WriteFile(FS, filepath.Join(basepath, "snooty.toml"), []byte("test"), 0644))
	check(iowrap.WriteFile(FS, filepath.Join(basepath, "source", "fundamentals", "aggregation.txt"), []byte(aggregationsFile), 0644))
	check(iowrap.WriteFile(FS, filepath.Join(basepath, "source", "fundamentals", "gridfs.txt"), []byte(grifsFile), 0644))

	expected := map[rst.RefTarget]string{
		{Target: "gridfs-create-bucket", Type: "local"}:        "/source/fundamentals/gridfs.txt",
		{Target: "gridfs-delete-bucket", Type: "local"}:        "/source/fundamentals/gridfs.txt",
		{Target: "gridfs-delete-files", Type: "local"}:         "/source/fundamentals/gridfs.txt",
		{Target: "gridfs-download-files", Type: "local"}:       "/source/fundamentals/gridfs.txt",
		{Target: "gridfs-rename-files", Type: "local"}:         "/source/fundamentals/gridfs.txt",
		{Target: "gridfs-retrieve-file-info", Type: "local"}:   "/source/fundamentals/gridfs.txt",
		{Target: "gridfs-upload-files", Type: "local"}:         "/source/fundamentals/gridfs.txt",
		{Target: "nodejs-aggregation-overview", Type: "local"}: "/source/fundamentals/aggregation.txt",
	}

	actual := GatherLocalRefs(GatherFiles(basepath))

	assert.EqualValues(t, expected, actual, "gatherConstants should return all constants in source directory")

}

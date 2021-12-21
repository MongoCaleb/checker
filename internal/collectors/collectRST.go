package collectors

import (
	"checker/internal/parsers/rst"
	"os"
	"path/filepath"
	"strings"

	iowrap "github.com/spf13/afero"

	log "github.com/sirupsen/logrus"
)

var (
	FS       iowrap.Fs
	FSUtil   *iowrap.Afero
	basepath string
)

func init() {
	FS = iowrap.NewOsFs()
	FSUtil = &iowrap.Afero{Fs: FS}
}

func exists(path string) bool {

	if _, err := FS.Stat(path); os.IsNotExist(err) {
		log.Errorf("%s does not exist", path)
		return false
	}
	return true
}

func snootyTomlExists(path string) bool {
	return exists(filepath.Join(path, "snooty.toml"))
}

func sourceDirectoryExists(path string) bool {

	return exists(filepath.Join(path, "source"))
}

func GatherFiles(path string) []string {
	basepath = path
	if !snootyTomlExists(path) || !sourceDirectoryExists(path) {
		log.Panic("snooty.toml or source directory does not exist")
	}

	files := make([]string, 0)

	err := FSUtil.Walk(basepath, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".rst" || filepath.Ext(path) == ".txt" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return files
}

func gather(files []string, fn func(filename string, data []byte)) {
	for _, file := range files {
		dat, err := FSUtil.ReadFile(file)
		if err != nil {
			log.Panic(err)
		}

		fileName := strings.Replace(file, basepath, "", 1)
		fn(fileName, dat)
	}
}

type RstRoleMap map[rst.RstRole]string

func GatherRoles(files []string) RstRoleMap {
	roles := make(map[rst.RstRole]string, len(files))
	gather(files, func(filename string, data []byte) {
		for _, role := range rst.ParseForRoles(data) {
			roles[role] = filename
		}
	})
	return roles
}

func (r *RstRoleMap) Get(key string) (*rst.RstRole, bool) {
	for k := range *r {
		if k.Name == key {
			return &k, true
		}
	}
	return nil, false
}

func GatherConstants(files []string) map[rst.RstConstant]string {
	consts := make(map[rst.RstConstant]string, len(files))
	gather(files, func(filename string, data []byte) {
		for _, con := range rst.ParseForConstants(data) {
			consts[con] = filename
		}
	})
	return consts
}

func GatherHTTPLinks(files []string) map[rst.RstHTTPLink]string {
	links := make(map[rst.RstHTTPLink]string, len(files))
	gather(files, func(filename string, data []byte) {
		for _, link := range rst.ParseForHTTPLinks(data) {
			links[link] = filename
		}
	})
	return links
}

type RefTargetMap map[rst.RefTarget]string

func GatherLocalRefs(files []string) RefTargetMap {
	refs := make(map[rst.RefTarget]string, len(files))
	gather(files, func(filename string, data []byte) {
		for _, ref := range rst.ParseForLocalRefs(data) {
			refs[ref] = filename
		}
	})
	return refs
}

func (r *RefTargetMap) Get(ref *rst.RstRole) (*rst.RefTarget, bool) {
	for k := range *r {
		if k.Target == ref.Target {
			return &k, true
		}
	}
	return nil, false
}

package collectors

import (
	"checker/internal/rstparsers"
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
	basepath, _ = os.Getwd()
}

func exists(path string) bool {

	if _, err := FS.Stat(filepath.Join(basepath, path)); os.IsNotExist(err) {
		log.Errorf("%s does not exist", path)
		return false
	}
	return true
}

func snootyTomlExists() bool {
	return exists("snooty.toml")
}

func sourceDirectoryExists() bool {

	return exists("source")
}

func gatherFiles() []string {
	if !snootyTomlExists() || !sourceDirectoryExists() {
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

func gatherRoles(files []string) map[string][]rstparsers.RstRole {
	roles := make(map[string][]rstparsers.RstRole, len(files))
	gather(files, func(filename string, data []byte) {
		roles[filename] = rstparsers.ParseForRoles(data)
	})
	return roles
}

func gatherConstants(files []string) map[string][]rstparsers.RstConstant {
	consts := make(map[string][]rstparsers.RstConstant, len(files))
	gather(files, func(filename string, data []byte) {
		consts[filename] = rstparsers.ParseForConstants(data)
	})
	return consts
}

func gatherHTTPLinks(files []string) map[string][]rstparsers.RstHTTPLink {
	links := make(map[string][]rstparsers.RstHTTPLink, len(files))
	gather(files, func(filename string, data []byte) {
		links[filename] = rstparsers.ParseForHTTPLinks(data)
	})
	return links
}

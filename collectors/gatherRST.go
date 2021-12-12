package collectors

import (
	"checker/rstparsers"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	iowrap "github.com/spf13/afero"

	log "github.com/sirupsen/logrus"
)

var (
	FS       iowrap.Fs
	FSUtil   *iowrap.Afero
	basepath string
	trimlen  int
)

func init() {
	FS = iowrap.NewOsFs()
	FSUtil = &iowrap.Afero{Fs: FS}
	_, b, _, ok := runtime.Caller(0)
	if !ok {
		log.Panic("Could not get caller")
	}
	basepath = filepath.Dir(b)
	trimlen = len(basepath)
}

func exists(path string) bool {
	_, b, _, ok := runtime.Caller(0)
	if !ok {
		log.Panic("Could not get caller")
	}
	basepath := filepath.Dir(b)
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
func gatherRoles() map[string][]rstparsers.RstRole {

	files := gatherFiles()

	type results struct {
		file  string
		roles []rstparsers.RstRole
	}

	var wg sync.WaitGroup

	allroles := make(map[string][]rstparsers.RstRole, 0)
	queue := make(chan results, len(files))

	for _, file := range files {
		wg.Add(1)
		go func(file string) {
			dat, err := FSUtil.ReadFile(file)
			if err != nil {
				log.Panic(err)
			}
			roles := rstparsers.ParseForRoles(string(dat))
			queue <- results{file: file[trimlen:], roles: roles}
		}(file)
	}

	go func() {
		for t := range queue {
			var roles []rstparsers.RstRole
			if allroles[t.file] == nil {
				roles = make([]rstparsers.RstRole, 0)
			} else {
				roles = allroles[t.file]
			}

			roles = append(roles, t.roles...)
			allroles[t.file] = roles

			wg.Done()
		}
	}()

	wg.Wait()

	return allroles
}

func gatherConstants() map[string][]rstparsers.RstConstant {

	files := gatherFiles()

	type results struct {
		file      string
		constants []rstparsers.RstConstant
	}

	var wg sync.WaitGroup

	allconstants := make(map[string][]rstparsers.RstConstant, 0)
	queue := make(chan results, len(files))

	for _, file := range files {
		wg.Add(1)
		go func(file string) {
			dat, err := FSUtil.ReadFile(file)
			if err != nil {
				log.Panic(err)
			}
			constants := rstparsers.ParseForConstants(string(dat))
			queue <- results{file: file[trimlen:], constants: constants}
		}(file)
	}

	go func() {
		for t := range queue {
			var constants []rstparsers.RstConstant
			if allconstants[t.file] == nil {
				constants = make([]rstparsers.RstConstant, 0)
			} else {
				constants = allconstants[t.file]
			}
			constants = append(constants, t.constants...)
			allconstants[t.file] = constants
			wg.Done()

		}
	}()

	wg.Wait()

	return allconstants
}

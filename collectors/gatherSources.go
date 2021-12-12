package collectors

import (
	"checker/rstparsers"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	iowrap "github.com/spf13/afero"

	log "github.com/sirupsen/logrus"
)

var (
	FS     iowrap.Fs
	FSUtil *iowrap.Afero
)

func init() {
	FS = iowrap.NewOsFs()
	FSUtil = &iowrap.Afero{Fs: FS}
}

func exists(path string) bool {
	_, b, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("Could not get caller")
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
	_, b, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("Could not get caller")
	}
	basepath := filepath.Dir(b)
	files := make([]string, 0)

	err := FSUtil.Walk(basepath, func(path string, info os.FileInfo, err error) error {
		log.Info(path)
		if filepath.Ext(path) == ".rst" || filepath.Ext(path) == ".txt" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return files
}
func gatherRoles() map[string]rstparsers.RstRole {

	files := gatherFiles()

	var wg sync.WaitGroup

	allroles := make([]rstparsers.RstRole, 0)
	queue := make(chan []rstparsers.RstRole, len(files))

	for _, file := range files {
		wg.Add(1)
		go func(file string) {
			dat, err := os.ReadFile(file)
			if err != nil {
				log.Fatal(err)
			}
			queue <- rstparsers.ParseForRoles(string(dat))
		}(file)
	}

	go func() {
		for t := range queue {
			allroles = append(allroles, t...)
			wg.Done()
		}
	}()

	wg.Wait()

	for _, role := range allroles {
		fmt.Println(role)
	}

	return make(map[string]rstparsers.RstRole)
}

package collectors

import (
	"checker/rstparsers"
	"errors"
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

func snootyTomlExists() bool {
	_, b, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("Could not get caller")
	}
	basepath := filepath.Dir(b)

	_, err := FS.Open(filepath.Join(basepath, "snooty.toml"))
	if errors.Is(err, os.ErrNotExist) {
		log.Error(errors.New("snooty.toml does not exist"))
		return false
	}
	return true
}

func main() {
	var files []string

	root := "source"
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get working directory: %v", err)
	}
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".rst" || filepath.Ext(path) == ".txt" {
			files = append(files, filepath.Join(cwd, path))
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

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

}

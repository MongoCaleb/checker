package sources

import (
	"context"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/google/go-github/v41/github"
	log "github.com/sirupsen/logrus"
)

const (
	rstSpecBase = "https://raw.githubusercontent.com/mongodb/snooty-parser/"
)

func GetLatestSnootyParserTag() string {
	ghClient := github.NewClient(nil)

	gctx, gcancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer gcancel()

	// get the latest release
	tags, _, err := ghClient.Repositories.ListTags(gctx, "mongodb", "snooty-parser", nil)
	if err != nil {
		log.Fatal(err)
	}

	latest := tags[0].Name
	return rstSpecBase + *latest + "/snooty/rstspec.toml"
}

func GetIntersphinxFile(input string) []byte {
	resp, err := http.Get(input)
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err)
	}
	return body
}

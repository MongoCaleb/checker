package utils

import (
	"context"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	"github.com/google/go-github/v41/github"
	log "github.com/sirupsen/logrus"
)

const (
	rstSpecBase = "https://raw.githubusercontent.com/mongodb/snooty-parser/"
)

var httpLinkRegex = regexp.MustCompile(`(https?:\/\/[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9]{1,6}\b[-a-zA-Z0-9@:%_\+.~#?&//=]*)`)

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

func GetNetworkFile(input string) []byte {
	resp, err := http.Get(input)
	if err != nil {
		log.Panicf("Could not get file %s: ", input, err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err)
	}
	return body
}

func GetLocalFile(input string) []byte {
	body, err := ioutil.ReadFile(input)
	if err != nil {
		log.Panic(err)
	}
	return body
}

func IsHTTPLink(input string) bool {
	return httpLinkRegex.MatchString(input)
}

func IsReachable(url string) bool {

	response, errors := http.Get(url)

	if errors != nil {
		log.Debug(url)
		log.Error(errors)
		return false
	}

	if response.StatusCode == 200 {
		return true
	}
	return false
}

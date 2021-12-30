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

var (
	httpLinkRegex = regexp.MustCompile(`(https?:\/\/[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9]{1,6}\b[-a-zA-Z0-9@:%_\+.~#?&//=]*)`)
	client        *http.Client
)

func init() {
	client = &http.Client{
		Timeout: time.Second * 4,
	}
}

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
	req, err := http.NewRequest("GET", input, nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Panicf("Could not get file %s: %v", input, err)
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

func IsReachable(url string) (string, bool) {
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Connection", "Keep-Alive")
	req.Header.Set("Accept-Language", "en-US")
	req.Header.Set("User-Agent", "Mozilla/5.0")
	if err != nil {
		log.Fatal(err)
	}
	response, errors := client.Do(req)

	if errors != nil {
		if response != nil {
			sc := response.StatusCode
			if sc < 400 {
				return "", true
			} else {
				return response.Status, false
			}
		} else {
			return "", false
		}
	}

	if response.StatusCode < 400 {
		return "", true
	}

	return "", false
}

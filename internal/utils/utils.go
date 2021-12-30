package utils

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/google/go-github/v41/github"
	log "github.com/sirupsen/logrus"
)

const (
	rstSpecBase = "https://raw.githubusercontent.com/mongodb/snooty-parser/"
)

type validRedirects [7]int

func (v validRedirects) contains(i int) bool {
	for _, a := range v {
		if a == i {
			return true
		}
	}
	return false
}

var (
	httpLinkRegex = regexp.MustCompile(`(https?:\/\/[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9]{1,6}\b[-a-zA-Z0-9@:%_\+.~#?&//=]*)`)
	client        *http.Client
	redirects     = validRedirects{301, 302, 303, 304, 305, 307, 308}
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

func IsReachable(url string) (error, bool) {

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Connection", "Keep-Alive")
	req.Header.Set("Accept-Language", "en-US")
	req.Header.Set("User-Agent", "Mozilla/5.0")
	if err != nil {
		log.Fatal(err)
	}
	response, err := client.Do(req)

	if err != nil {
		if strings.Contains(err.Error(), "stopped after 10 redirects") {
			if redirects.contains(response.StatusCode) {
				return nil, true
			}
		} else {
			return err, false
		}
	}
	if response.StatusCode == 200 {
		return nil, true
	} else {
		return fmt.Errorf("%d", response.StatusCode), false
	}
}

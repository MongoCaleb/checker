package utils

import (
	"context"
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

type HttpResponse struct {
	Code     int
	Filename string
	Message  string
}

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
	log.SetFormatter(&log.TextFormatter{
		PadLevelText:           true,
		DisableLevelTruncation: false,
	})
	client = &http.Client{
		Timeout: time.Second * 30,
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

func IsReachable(uri string) (HttpResponse, bool) {
	// check to see if there's a way to avoid triggering page viewws
	// block add blockers
	// test net.DialTCP
	// look at muffet to see what they do to make sure a url is valid

	var r HttpResponse

	req, err := http.NewRequest("GET", uri, nil)
	req.Header.Set("Connection", "Keep-Alive")
	req.Header.Set("Accept-Language", "en-US")
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	if err != nil {
		log.Fatal(err)
	}
	response, err := client.Do(req)

	if err != nil {
		var code int
		if response != nil {
			code = response.StatusCode
		}
		if strings.Contains(err.Error(), "stopped after 10 redirects") {
			if redirects.contains(code) {
				return r, true
			}
		} else {
			r.Code = code
			return r, false
		}
	}
	defer response.Body.Close()

	if response.StatusCode == 200 {
		return r, true
	} else {
		r.Code = response.StatusCode
		r.Message = req.URL.Path
		return r, false
	}
}

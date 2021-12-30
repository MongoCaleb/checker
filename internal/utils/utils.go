package utils

import (
	"context"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	"golang.org/x/time/rate"

	"github.com/google/go-github/v41/github"
	log "github.com/sirupsen/logrus"
)

const (
	rstSpecBase = "https://raw.githubusercontent.com/mongodb/snooty-parser/"
)

var (
	client        *RLHTTPClient
	httpLinkRegex = regexp.MustCompile(`(https?:\/\/[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9]{1,6}\b[-a-zA-Z0-9@:%_\+.~#?&//=]*)`)
)

func init() {
	rl := rate.NewLimiter(rate.Every(10*time.Second), 100)
	client = NewClient(rl)
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

func IsReachable(url string) (*http.Response, bool) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	response, errors := client.Do(req)

	if errors != nil {
		return response, false
	}

	if response.StatusCode == 200 {
		return nil, true
	}

	return nil, false
}

type RLHTTPClient struct {
	client      *http.Client
	Ratelimiter *rate.Limiter
}

func (c *RLHTTPClient) Do(req *http.Request) (*http.Response, error) {
	ctx := context.Background()
	err := c.Ratelimiter.Wait(ctx)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func NewClient(rl *rate.Limiter) *RLHTTPClient {
	c := &RLHTTPClient{
		client:      http.DefaultClient,
		Ratelimiter: rl,
	}
	return c
}

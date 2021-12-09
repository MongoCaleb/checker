package types

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
)

const (
	rstSpec = "https://raw.githubusercontent.com/mongodb/snooty-parser/master/snooty/rstspec.toml"
)

type RawMap struct {
	Roles map[string]interface{} `toml:"role"`
}

// RoleMap contains roles from rstSpec.toml
type RoleMap struct {
	Roles map[string]string
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	Client HTTPClient
)

func init() {
	Client = &http.Client{}
}

func NewRoleMap() RoleMap {

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", rstSpec, nil)
	if err != nil {
		log.Errorf("Error creating request to url %s: %v", rstSpec, err)
	}
	resp, err := Client.Do(req)
	if err != nil {
		log.Errorf("Error getting response from url %s: %v", rstSpec, err)
	}
	defer resp.Body.Close()

	buff, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	var rawmap RawMap
	_, err = toml.Decode(string(buff), &rawmap)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	for k, v := range rawmap.Roles {
		switch (v.(map[string]interface{})["type"]).(type) {
		case map[string]interface{}:
			continue
		default:
			delete(rawmap.Roles, k)
		}
	}

	for k, v := range rawmap.Roles {
		urlObj := v.(map[string]interface{})["type"]
		target := urlObj.(map[string]interface{})["link"]
		rawmap.Roles[k] = target
	}
	var roleMap RoleMap
	roleMap.Roles = make(map[string]string, len(rawmap.Roles))
	for k, v := range rawmap.Roles {
		if v != nil {
			roleMap.Roles[k] = v.(string)
		}
	}
	return roleMap
}

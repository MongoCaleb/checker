package sources

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
)

type RawMap struct {
	Roles map[string]interface{} `toml:"role"`
}

// RoleMap contains roles from rstSpec.toml
type RoleMap map[string]string

func init() {
	Client = &http.Client{}
}

func NewRoleMap(rstSpecURL string) RoleMap {

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	// get the release version of rstspec.toml
	req, err := http.NewRequestWithContext(ctx, "GET", rstSpecURL, nil)
	if err != nil {
		log.Errorf("Error creating request to url %s: %v", rstSpecURL, err)
	}
	resp, err := Client.Do(req)
	if err != nil {
		log.Errorf("Error getting response from url %s: %v", rstSpecURL, err)
	}
	defer resp.Body.Close()

	buff, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// populate a raw role map that is map[string]interface{}
	var rawmap RawMap
	_, err = toml.Decode(string(buff), &rawmap)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// filter out roles that aren't links, and convert to a RoleMap
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
	roleMap := make(map[string]string, len(rawmap.Roles))
	for k, v := range rawmap.Roles {
		if v != nil {
			roleMap[k] = v.(string)
		}
	}
	return roleMap
}

func (r *RoleMap) Get(key string) (string, bool) {
	val, ok := (*r)[key]
	return val, ok
}

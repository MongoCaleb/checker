package sources

import (
	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
)

type rawMap struct {
	Roles map[string]interface{} `toml:"role"`
}

// RoleMap contains roles from rstSpec.toml
type RoleMap map[string]string

func NewRoleMap(input []byte) RoleMap {

	// populate a raw role map that is map[string]interface{}
	var rawmap rawMap
	_, err := toml.Decode(string(input), &rawmap)
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

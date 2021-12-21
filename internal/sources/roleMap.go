package sources

import (
	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
)

type RawMap struct {
	Roles map[string]interface{} `toml:"role"`
}

type RoleMap struct {
	Links LinkRoleMap
	Raw   RawMap
}

// LinkRoleMap contains roles from rstspec.toml
type LinkRoleMap map[string]string

// OtherRoleMap contains other roles from rstspec.toml, like guilabel
type OtherRoleMap map[string]string

func NewRoleMap(input []byte) *RoleMap {

	// populate a raw role map that is map[string]interface{}
	var rawmap RawMap
	_, err := toml.Decode(string(input), &rawmap)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	var copyRaw RawMap
	copyRaw.Roles = make(map[string]interface{}, len(rawmap.Roles))
	for k, v := range rawmap.Roles {
		copyRaw.Roles[k] = v
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
	roleMap := make(LinkRoleMap, len(rawmap.Roles))
	for k, v := range rawmap.Roles {
		if v != nil {
			roleMap[k] = v.(string)
		}
	}
	return &RoleMap{
		Links: roleMap,
		Raw:   copyRaw,
	}
}

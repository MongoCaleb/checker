package sources

import (
	"strings"

	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
)

type RawRstSpec struct {
	Roles      map[string]interface{} `toml:"role"`
	RstObjects map[string]interface{} `toml:"rstobject"`
	Directives map[string]interface{} `toml:"directive"`
}

type RstSpec struct {
	Roles      RolesMap
	RawRoles   map[string]bool
	Directives map[string]bool
	RstObjects map[string]bool
}

// RolesMap contains roles from rstspec.toml
type RolesMap map[string]string

// OtherRoleMap contains other roles from rstspec.toml, like guilabel
type OtherRoleMap map[string]string

func NewRoleMap(input []byte) *RstSpec {

	var rstSpec RstSpec

	// populate a raw role map that is map[string]interface{}
	var rawmap RawRstSpec
	_, err := toml.Decode(string(input), &rawmap)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// log.SetLevel(log.DebugLevel)
	// log.Debugf("rstspec: %v", rawmap.Directives)
	rstSpec.populateRoles(&rawmap)
	rstSpec.populateDirectives(&rawmap)
	rstSpec.populateRstObjects(&rawmap)
	return &rstSpec
}

func (r *RstSpec) populateRoles(raw *RawRstSpec) {
	r.RawRoles = make(map[string]bool, len(raw.Roles))
	for k := range raw.Roles {
		r.RawRoles[k] = true
	}

	// filter out roles that aren't links, and convert to a RoleMap
	for k, v := range raw.Roles {
		switch (v.(map[string]interface{})["type"]).(type) {
		case map[string]interface{}:
			continue
		default:
			delete(raw.Roles, k)
		}
	}
	for k, v := range raw.Roles {
		urlObj := v.(map[string]interface{})["type"]
		target := urlObj.(map[string]interface{})["link"]
		raw.Roles[k] = target
	}
	roleMap := make(RolesMap, len(raw.Roles))
	for k, v := range raw.Roles {
		if v != nil {
			roleMap[k] = v.(string)
		}
	}
	r.Roles = roleMap

}

func (r *RstSpec) populateDirectives(raw *RawRstSpec) {
	r.Directives = make(map[string]bool, len(raw.Directives))

	for k := range raw.Directives {
		r.Directives[k] = true
	}
}

func (r *RstSpec) populateRstObjects(raw *RawRstSpec) {
	r.RstObjects = make(map[string]bool, len(raw.RstObjects))

	for k := range raw.RstObjects {
		target := strings.Split(k, ":")
		if len(target) > 1 {
			r.RstObjects[target[1]] = true
		} else {
			r.RstObjects[k] = true
		}
	}
}

package types

import (
	"regexp"

	"github.com/BurntSushi/toml"

	log "github.com/sirupsen/logrus"
)

type TomlConfig struct {
	Name        string            `toml:"name"`
	Title       string            `toml:"title"`
	Constants   map[string]string `toml:"constants"`
	Intersphinx []string          `toml:"intersphinx"`
}

func NewTomlConfig(input string) (*TomlConfig, error) {
	var cfg TomlConfig
	_, err := toml.Decode(input, &cfg)
	if err != nil {
		return nil, err
	}

	cfg.Constants = cfg.resolveConstants()

	return &cfg, nil
}

func (cfg *TomlConfig) resolveConstants() map[string]string {
	newMap := make(map[string]string, len(cfg.Constants))
	re := regexp.MustCompile(`\{\+([\w\s\-\.\d_=+!@#$%^&*(\)]*)\+\}`)
	for k, v := range cfg.Constants {
		loc := re.FindIndex([]byte(v))
		if len(loc) == 0 {
			newMap[k] = v
		} else {
			newMap[k] = descendConstants(cfg.Constants, v, 0)
		}
	}
	return newMap
}

func descendConstants(constantMap map[string]string, value string, depth int8) string {
	if depth > 4 {
		log.Warnf("Constant interpolation is reaching ridiculous levels. Resolving %s and have reached a depth of %d", value, depth)
	}
	re := regexp.MustCompile(`\{\+([\w\s\-\.\d_=+!@#$%^&*(\)]*)\+\}`)
	loc := re.FindIndex([]byte(value))
	if len(loc) == 0 {
		return value
	}
	toFind := value[loc[0]+len("{+") : loc[1]-len("+}")]
	lookup, ok := constantMap[toFind]
	if !ok {
		log.Errorf("Could not find constant %s", toFind)
	}
	newVal := value[:loc[0]] + lookup + value[loc[1]:]
	return descendConstants(constantMap, newVal, depth+1)
}
